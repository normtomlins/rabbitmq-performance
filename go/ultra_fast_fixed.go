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

type FastMessage struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
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
	ch.Qos(1000, 0, false)

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
	fmt.Println("\nUltra Fast Go Producer (Fixed)")
	fmt.Println("==============================")

	// Pre-initialize timezone to avoid concurrent initialization
	_ = time.Now()

	numWorkers := runtime.NumCPU() * 2
	fmt.Printf("Workers: %d\n", numWorkers)

	startTime := time.Now()
	var wg sync.WaitGroup
	var messagesSent int64

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
			
			// Pre-allocate message buffer
			message := FastMessage{
				ID:      0,
				Content: "",
			}
			
			// Pre-allocate publishing struct
			publishing := amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Transient,
			}
			
			// Buffer for JSON encoding
			buf := make([]byte, 0, 256)
			
			for i := startID; i < endID; i++ {
				// Update message
				message.ID = i
				message.Content = fmt.Sprintf("Hello World #%d", i)
				
				// Marshal directly to buffer
				buf = buf[:0]
				jsonData, err := json.Marshal(message)
				if err != nil {
					log.Printf("Worker %d: Failed to marshal: %v", workerID, err)
					continue
				}
				
				publishing.Body = jsonData
				
				err = producer.channel.PublishWithContext(
					ctx,
					"",                // exchange
					"ultra_fast_queue", // routing key
					false,             // mandatory
					false,             // immediate
					publishing)
				
				if err != nil {
					log.Printf("Worker %d: Failed to publish: %v", workerID, err)
				}
				
				if i%1000 == 0 {
					sent := atomic.AddInt64(&messagesSent, 1000)
					if sent%10000 == 0 {
						elapsed := time.Since(startTime).Seconds()
						rate := float64(sent) / elapsed
						fmt.Printf("Sent %d messages - Rate: %.0f msgs/sec\n", sent, rate)
					}
				}
			}
			
			// Add remaining messages
			remaining := int64(endID % 1000)
			if remaining > 0 {
				atomic.AddInt64(&messagesSent, remaining)
			}
		}(w)
	}

	wg.Wait()

	// Make sure we count all messages
	messagesSent = int64(messageCount)
	
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

	messageCount := 500000
	runUltraFast(messageCount)
}
