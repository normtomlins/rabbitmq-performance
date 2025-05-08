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

type BatchProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewBatchProducer(url string) (*BatchProducer, error) {
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
	ch.Qos(5000, 0, false)

	return &BatchProducer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *BatchProducer) Close() {
	p.channel.Close()
	p.conn.Close()
}

func runUltraBatch(messageCount int) {
	fmt.Println("\nUltra Batch Go Producer")
	fmt.Println("======================")

	numWorkers := runtime.NumCPU()
	batchSize := 1000
	fmt.Printf("Workers: %d, Batch size: %d\n", numWorkers, batchSize)

	startTime := time.Now()
	var wg sync.WaitGroup
	var messagesSent int64

	// Create a pool of connections
	producers := make([]*BatchProducer, numWorkers)
	for i := 0; i < numWorkers; i++ {
		producer, err := NewBatchProducer("amqp://user:password@rabbitmq:5672/")
		if err != nil {
			log.Fatalf("Failed to create producer %d: %v", i, err)
		}
		
		_, err = producer.channel.QueueDeclare(
			"batch_queue",
			false, // not durable
			false, // delete when unused
			false, // exclusive
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
			
			// Pre-allocate batch
			batch := make([][]byte, batchSize)
			for i := range batch {
				batch[i] = make([]byte, 0, 256)
			}
			
			// Pre-allocate messages
			messages := make([]map[string]interface{}, batchSize)
			for i := range messages {
				messages[i] = map[string]interface{}{
					"id":      0,
					"content": "",
				}
			}
			
			for i := startID; i < endID; i += batchSize {
				batchEnd := i + batchSize
				if batchEnd > endID {
					batchEnd = endID
				}
				currentBatchSize := batchEnd - i
				
				// Prepare batch
				for j := 0; j < currentBatchSize; j++ {
					msgID := i + j
					messages[j]["id"] = msgID
					messages[j]["content"] = fmt.Sprintf("Hello World #%d", msgID)
					
					// Marshal directly to pre-allocated buffer
					jsonData, err := json.Marshal(messages[j])
					if err != nil {
						log.Printf("Worker %d: Failed to marshal: %v", workerID, err)
						continue
					}
					batch[j] = jsonData
				}
				
				// Send batch using transactions for better throughput
				err := producer.channel.Tx()
				if err != nil {
					log.Printf("Worker %d: Failed to start transaction: %v", workerID, err)
					continue
				}
				
				for j := 0; j < currentBatchSize; j++ {
					err = producer.channel.PublishWithContext(
						ctx,
						"",           // exchange
						"batch_queue", // routing key
						false,        // mandatory
						false,        // immediate
						amqp.Publishing{
							ContentType:  "application/json",
							Body:         batch[j],
							DeliveryMode: amqp.Transient,
						})
					
					if err != nil {
						log.Printf("Worker %d: Failed to publish: %v", workerID, err)
						producer.channel.TxRollback()
						break
					}
				}
				
				err = producer.channel.TxCommit()
				if err != nil {
					log.Printf("Worker %d: Failed to commit transaction: %v", workerID, err)
					continue
				}
				
				sent := atomic.AddInt64(&messagesSent, int64(currentBatchSize))
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	fmt.Println("Ultra Batch RabbitMQ Go Producer")
	fmt.Println("===============================")
	fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	messageCount := 500000
	runUltraBatch(messageCount)
}
