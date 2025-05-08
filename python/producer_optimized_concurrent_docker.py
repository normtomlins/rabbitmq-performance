#!/usr/bin/env python3
import pika
import json
import datetime
import time
import os
from concurrent.futures import ProcessPoolExecutor, ThreadPoolExecutor
import multiprocessing

class OptimizedProducer:
    def __init__(self):
        self.connection = None
        self.channel = None
    
    def setup_connection(self):
        """Create a single shared connection"""
        connection_params = pika.ConnectionParameters(
            host='rabbitmq',
            port=5672,
            credentials=pika.PlainCredentials('user', 'password')
        )
        
        self.connection = pika.BlockingConnection(connection_params)
        self.channel = self.connection.channel()
        
        # Use a non-durable queue for maximum speed
        self.channel.queue_declare(queue='optimized_queue', durable=False)
        
        # Optimize channel settings
        self.channel.confirm_delivery()
    
    def send_batch_single_connection(self, messages):
        """Send a batch of messages using the shared connection"""
        sent = 0
        for msg in messages:
            try:
                self.channel.basic_publish(
                    exchange='',
                    routing_key='optimized_queue',
                    body=json.dumps(msg),
                    properties=pika.BasicProperties(
                        delivery_mode=1,  # Non-persistent
                    )
                )
                sent += 1
            except Exception as e:
                print(f"Error: {e}")
        return sent
    
    def run_single_thread_batched(self, message_count=50000):
        """Single thread, but with batched sending"""
        self.setup_connection()
        
        start_time = time.time()
        batch_size = 1000
        total_sent = 0
        
        print(f"Single thread batched: Sending {message_count} messages...")
        
        for batch_start in range(0, message_count, batch_size):
            # Create batch
            batch = []
            for i in range(batch_start, min(batch_start + batch_size, message_count)):
                batch.append({
                    'id': i,
                    'content': f'Hello World #{i}',
                    'timestamp': datetime.datetime.now().isoformat()
                })
            
            # Send batch
            sent = self.send_batch_single_connection(batch)
            total_sent += sent
            
            # Progress
            if total_sent % 10000 == 0:
                elapsed = time.time() - start_time
                rate = total_sent / elapsed
                print(f'Sent {total_sent} messages - Rate: {rate:.0f} msgs/sec')
        
        # Final stats
        end_time = time.time()
        total_time = end_time - start_time
        rate = total_sent / total_time
        
        self.connection.close()
        
        print(f'\nCompleted!')
        print(f'Total messages: {total_sent}')
        print(f'Total time: {total_time:.2f} seconds')
        print(f'Average rate: {rate:.0f} messages/second')
        
        return rate


def process_worker(worker_id, start_id, message_count):
    """Independent process worker with its own connection"""
    connection_params = pika.ConnectionParameters(
        host='rabbitmq',
        port=5672,
        credentials=pika.PlainCredentials('user', 'password')
    )
    
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()
    channel.queue_declare(queue='process_queue', durable=False)
    
    start_time = time.time()
    
    for i in range(start_id, start_id + message_count):
        message = {
            'id': i,
            'content': f'Hello World #{i}',
            'timestamp': datetime.datetime.now().isoformat()
        }
        
        channel.basic_publish(
            exchange='',
            routing_key='process_queue',
            body=json.dumps(message),
            properties=pika.BasicProperties(delivery_mode=1)
        )
    
    connection.close()
    
    end_time = time.time()
    duration = end_time - start_time
    rate = message_count / duration if duration > 0 else 0
    
    return worker_id, message_count, duration, rate


def run_multiprocess_test(total_messages=50000, num_processes=4):
    """Use multiprocessing to bypass GIL"""
    print(f"\nMultiprocess test: {num_processes} processes, {total_messages} messages")
    
    messages_per_process = total_messages // num_processes
    
    start_time = time.time()
    
    with ProcessPoolExecutor(max_workers=num_processes) as executor:
        futures = []
        for i in range(num_processes):
            start_id = i * messages_per_process
            future = executor.submit(process_worker, i, start_id, messages_per_process)
            futures.append(future)
        
        # Collect results
        total_sent = 0
        for future in futures:
            worker_id, sent, duration, rate = future.result()
            total_sent += sent
            print(f'Process {worker_id}: {sent} messages in {duration:.2f}s ({rate:.0f} msgs/sec)')
    
    end_time = time.time()
    total_time = end_time - start_time
    overall_rate = total_sent / total_time
    
    print(f'\nTotal messages: {total_sent}')
    print(f'Total time: {total_time:.2f} seconds')
    print(f'Overall rate: {overall_rate:.0f} messages/second')
    
    return overall_rate


def run_optimized_single_thread(message_count=50000):
    """Optimized single thread without any overhead"""
    connection_params = pika.ConnectionParameters(
        host='rabbitmq',
        port=5672,
        credentials=pika.PlainCredentials('user', 'password')
    )
    
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()
    channel.queue_declare(queue='fast_queue', durable=False)
    
    start_time = time.time()
    
    print(f"Optimized single thread: Sending {message_count} messages...")
    
    # Pre-create message template
    message_template = {
        'id': 0,
        'content': 'Hello World #0',
        'timestamp': ''
    }
    
    for i in range(message_count):
        # Minimal message updates
        message_template['id'] = i
        message_template['content'] = f'Hello World #{i}'
        
        channel.basic_publish(
            exchange='',
            routing_key='fast_queue',
            body=json.dumps(message_template),
            properties=pika.BasicProperties(delivery_mode=1)
        )
        
        # Progress
        if i > 0 and i % 10000 == 0:
            elapsed = time.time() - start_time
            rate = i / elapsed
            print(f'Sent {i} messages - Rate: {rate:.0f} msgs/sec')
    
    connection.close()
    
    end_time = time.time()
    total_time = end_time - start_time
    rate = message_count / total_time
    
    print(f'\nCompleted!')
    print(f'Total messages: {message_count}')
    print(f'Total time: {total_time:.2f} seconds')
    print(f'Average rate: {rate:.0f} messages/second')
    
    return rate


if __name__ == "__main__":
    print("RabbitMQ Concurrency Performance Analysis")
    print("========================================\n")
    
    # Test 1: Optimized single thread (baseline)
    print("Test 1: Optimized single thread (baseline)")
    rate1 = run_optimized_single_thread(50000)
    
    time.sleep(2)
    
    # Test 2: Single thread with batching
    print("\nTest 2: Single thread with batching")
    producer = OptimizedProducer()
    rate2 = producer.run_single_thread_batched(50000)
    
    time.sleep(2)
    
    # Test 3: Multiprocessing (bypasses GIL)
    print("\nTest 3: Multiprocessing")
    rate3 = run_multiprocess_test(50000, 4)
    
    # Summary
    print("\n" + "="*50)
    print("PERFORMANCE SUMMARY")
    print("="*50)
    print(f"Optimized single thread: {rate1:.0f} msgs/sec")
    print(f"Single thread batched:   {rate2:.0f} msgs/sec")
    print(f"Multiprocessing (4):     {rate3:.0f} msgs/sec")
    print("\nConclusion:")
    print("- Python's GIL prevents true thread parallelism")
    print("- Multiple connections add overhead")
    print("- Single optimized thread often outperforms multiple threads")
    print("- Multiprocessing can help but has its own overhead")
