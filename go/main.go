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

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(url string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Producer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *Producer) Close() {
	p.channel.Close()
	p.conn.Close()
}

func (p *Producer) declareQueue(queueName string) error {
	_, err := p.channel.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// Single goroutine producer
func runSingleGoroutine(messageCount int) {
	fmt.Println("\nTest 1: Single Goroutine Producer")
	fmt.Println("================================")

	producer, err := NewProducer("amqp://user:password@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer producer.Close()

	if err := producer.declareQueue("go_queue"); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	startTime := time.Now()
	ctx := context.Background()
	
	for i := 0; i < messageCount; i++ {
		msg := Message{
			ID:        i,
			Content:   fmt.Sprintf("Hello World #%d", i),
			Timestamp: time.Now(),
		}

		body, _ := json.Marshal(msg)

		err := producer.channel.PublishWithContext(
			ctx,
			"",         // exchange
			"go_queue", // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Transient, // non-persistent
			})

		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		}

		if i > 0 && i%10000 == 0 {
			elapsed := time.Since(startTime).Seconds()
			rate := float64(i) / elapsed
			fmt.Printf("Sent %d messages - Rate: %.0f msgs/sec\n", i, rate)
		}
	}

	elapsed := time.Since(startTime).Seconds()
	rate := float64(messageCount) / elapsed

	fmt.Printf("\nCompleted!\n")
	fmt.Printf("Total messages: %d\n", messageCount)
	fmt.Printf("Total time: %.2f seconds\n", elapsed)
	fmt.Printf("Average rate: %.0f messages/second\n", rate)
}

// Multiple goroutines with worker pool
func runWorkerPool(messageCount int, numWorkers int) {
	fmt.Printf("\nTest 2: Worker Pool (%d workers)\n", numWorkers)
	fmt.Println("================================")

	startTime := time.Now()
	var wg sync.WaitGroup
	messagesCh := make(chan int, 1000)
	var messagesSent int64

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			producer, err := NewProducer("amqp://user:password@rabbitmq:5672/")
			if err != nil {
				log.Printf("Worker %d: Failed to connect: %v", workerID, err)
				return
			}
			defer producer.Close()

			if err := producer.declareQueue("go_worker_queue"); err != nil {
				log.Printf("Worker %d: Failed to declare queue: %v", workerID, err)
				return
			}

			ctx := context.Background()

			for msgID := range messagesCh {
				msg := Message{
					ID:        msgID,
					Content:   fmt.Sprintf("Hello World #%d", msgID),
					Timestamp: time.Now(),
				}

				body, _ := json.Marshal(msg)

				err := producer.channel.PublishWithContext(
					ctx,
					"",                // exchange
					"go_worker_queue", // routing key
					false,             // mandatory
					false,             // immediate
					amqp.Publishing{
						ContentType:  "application/json",
						Body:         body,
						DeliveryMode: amqp.Transient,
					})

				if err != nil {
					log.Printf("Worker %d: Failed to publish message: %v", workerID, err)
				}

				sent := atomic.AddInt64(&messagesSent, 1)
				if sent%10000 == 0 {
					elapsed := time.Since(startTime).Seconds()
					rate := float64(sent) / elapsed
					fmt.Printf("Sent %d messages - Rate: %.0f msgs/sec\n", sent, rate)
				}
			}
		}(w)
	}

	// Feed messages to workers
	go func() {
		for i := 0; i < messageCount; i++ {
			messagesCh <- i
		}
		close(messagesCh)
	}()

	// Wait for all workers to complete
	wg.Wait()

	elapsed := time.Since(startTime).Seconds()
	rate := float64(messagesSent) / elapsed

	fmt.Printf("\nCompleted!\n")
	fmt.Printf("Total messages: %d\n", messagesSent)
	fmt.Printf("Total time: %.2f seconds\n", elapsed)
	fmt.Printf("Average rate: %.0f messages/second\n", rate)
}

// Batch publishing with goroutines
func runBatchPublishing(messageCount int, numGoroutines int, batchSize int) {
	fmt.Printf("\nTest 3: Batch Publishing (%d goroutines, batch size %d)\n", numGoroutines, batchSize)
	fmt.Println("========================================================")

	startTime := time.Now()
	var wg sync.WaitGroup
	var messagesSent int64

	messagesPerGoroutine := messageCount / numGoroutines

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			producer, err := NewProducer("amqp://user:password@rabbitmq:5672/")
			if err != nil {
				log.Printf("Goroutine %d: Failed to connect: %v", goroutineID, err)
				return
			}
			defer producer.Close()

			if err := producer.declareQueue("go_batch_queue"); err != nil {
				log.Printf("Goroutine %d: Failed to declare queue: %v", goroutineID, err)
				return
			}

			ctx := context.Background()
			startID := goroutineID * messagesPerGoroutine
			endID := startID + messagesPerGoroutine

			// Process in batches
			for i := startID; i < endID; i += batchSize {
				batchEnd := i + batchSize
				if batchEnd > endID {
					batchEnd = endID
				}

				// Send batch
				for j := i; j < batchEnd; j++ {
					msg := Message{
						ID:        j,
						Content:   fmt.Sprintf("Hello World #%d", j),
						Timestamp: time.Now(),
					}

					body, _ := json.Marshal(msg)

					err := producer.channel.PublishWithContext(
						ctx,
						"",               // exchange
						"go_batch_queue", // routing key
						false,            // mandatory
						false,            // immediate
						amqp.Publishing{
							ContentType:  "application/json",
							Body:         body,
							DeliveryMode: amqp.Transient,
						})

					if err != nil {
						log.Printf("Goroutine %d: Failed to publish message: %v", goroutineID, err)
					}
				}

				sent := atomic.AddInt64(&messagesSent, int64(batchEnd-i))
				if sent%10000 == 0 {
					elapsed := time.Since(startTime).Seconds()
					rate := float64(sent) / elapsed
					fmt.Printf("Sent %d messages - Rate: %.0f msgs/sec\n", sent, rate)
				}
			}
		}(g)
	}

	wg.Wait()

	elapsed := time.Since(startTime).Seconds()
	rate := float64(messagesSent) / elapsed

	fmt.Printf("\nCompleted!\n")
	fmt.Printf("Total messages: %d\n", messagesSent)
	fmt.Printf("Total time: %.2f seconds\n", elapsed)
	fmt.Printf("Average rate: %.0f messages/second\n", rate)
}

func main() {
	fmt.Println("RabbitMQ Go Producer Performance Test")
	fmt.Println("====================================")
	fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	messageCount := 50000

	// Test 1: Single goroutine
	runSingleGoroutine(messageCount)
	time.Sleep(2 * time.Second)

	// Test 2: Worker pool
	runWorkerPool(messageCount, runtime.NumCPU())
	time.Sleep(2 * time.Second)

	// Test 3: Batch publishing with goroutines
	runBatchPublishing(messageCount, runtime.NumCPU(), 100)

	fmt.Println("\nAll tests completed!")
}
