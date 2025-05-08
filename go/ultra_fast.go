package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UltraFastProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewUltraFastProducer(url string) (*UltraFastProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Optimize channel settings
	ch.Qos(1000, 0, false) // prefetch count

	return &UltraFastProducer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *UltraFastProducer) Close() {
	p.channel.Close()
	p.conn.Close()
}

// Ultra-fast version with all optimizations
func runUltraFast(messageCount int) {
	fmt.Println("\nUltra Fast Go Producer")
	fmt.Println("=====================")

	numWorkers := runtime.NumCPU() * 2 // Oversubscribe
	fmt.Printf("Workers: %d\n", numWorkers)

	startTime := time.Now()
	var wg sync.WaitGroup
	var messagesSent int64

	// Pre-allocate message template
	msgTemplate := map[string]interface{}{
		"id":        0,
		"content":   "",
		"timestamp": "",
	}

	// Create a pool of connections
	producers := make([]*UltraFastProducer, numWorkers)
	for i := 0; i < numWorkers; i++ {
		producer, err := NewUltraFastProducer("amqp://user:password@rabbitmq:5672/")
		if err != nil {
			log.Fatalf("Failed to create producer %d: %v", i, err)
		}
		
		_, err = producer.channel.QueueDeclare(
			"ultra_fast_queue",
			false, // not durable
			false, // don't delete when unused
			false, // not exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare queue: %v", err)
		}
		
		producers[i] = producer
	}
	defer func() {
		for _, p := range producers {
			p.Close()
		}
	}()

	messagesPerWorker := messageCount / numWorkers
	ctx := context.Background()
	
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			producer := producers[workerID]
			startID := workerID * messagesPerWorker
			endID := startID + messagesPerWorker
			
			// Pre-serialize template
			msgTemplate["id"] = 0
			msgTemplate["content"] = "Hello World #0"
			msgTemplate["timestamp"] = time.Now().Format(time.RFC3339)
			templateBytes, _ := json.Marshal(msgTemplate)
			
			// Create a batch of messages
			batchSize := 100
			batch := make([]amqp.Publishing, batchSize)
			
			for i := 0; i < batchSize; i++ {
				batch[i] = amqp.Publishing{
					ContentType:  "application/json",
					Body:         make([]byte, len(templateBytes)),
					DeliveryMode: amqp.Transient,
				}
			}
			
			for i := startID; i < endID; i += batchSize {
				batchEnd := i + batchSize
				if batchEnd > endID {
					batchEnd = endID
				}
				
				// Update batch messages
				for j := 0; j < batchEnd-i; j++ {
					msgID := i + j
					msgTemplate["id"] = msgID
					msgTemplate["content"] = fmt.Sprintf("Hello World #%d", msgID)
					body, _ := json.Marshal(msgTemplate)
					batch[j].Body = body
				}
				
				// Send entire batch
				for j := 0; j < batchEnd-i; j++ {
					err := producer.channel.PublishWithContext(
						ctx,
						"",                // exchange
						"ultra_fast_queue", // routing key
						false,             // mandatory
						false,             // immediate
						batch[j])
					
					if err != nil {
						log.Printf("Worker %d: Failed to publish: %v", workerID, err)
					}
				}
				
				sent := atomic.AddInt64(&messagesSent, int64(batchEnd-i))
				if sent%10000 == 0 {
					elapsed := time.Since(startTime).Seconds()
					rate := float64(sent) / elapsed
					fmt.Printf("Sent %d messages - Rate: %.0f msgs/sec\n", sent, rate)
				}
			}
		}(w)
	}

	wg.Wait()

	elapsed := time.Since(startTime).Seconds()
	rate := float64(messagesSent) / elapsed

	fmt.Printf("\nCompleted!\n")
	fmt.Printf("Total messages: %d\n", messagesSent)
	fmt.Printf("Total time: %.2f seconds\n", elapsed)
	fmt.Printf("Average rate: %.0f messages/second\n\n", rate)
}

func main() {
	// Maximize CPU usage
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	fmt.Println("Ultra Fast RabbitMQ Go Producer")
	fmt.Println("==============================")
	fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	messageCount := 500000 // Test with more messages
	runUltraFast(messageCount)
}
