#!/usr/bin/env python3
import pika
import json
import sys
import os
import time

# Establish connection to RabbitMQ
connection_params = pika.ConnectionParameters(
    host='rabbitmq',  # Use container name instead of localhost
    port=5672,  # Use internal port, not mapped port
    credentials=pika.PlainCredentials('user', 'password')
)

# Retry connection logic for Docker environment
max_retries = 5
retry_count = 0
connection = None

while retry_count < max_retries:
    try:
        connection = pika.BlockingConnection(connection_params)
        print(f'Connected to RabbitMQ')
        break
    except pika.exceptions.AMQPConnectionError:
        retry_count += 1
        print(f'Connection attempt {retry_count} failed. Retrying in 5 seconds...')
        time.sleep(5)

if connection is None:
    print('Failed to connect to RabbitMQ after multiple attempts')
    exit(1)

# Statistics tracking (moved outside try block)
message_count = 0
start_time = None

try:
    channel = connection.channel()
    channel.queue_declare(queue='hello_queue', durable=True)

    # Define callback function to process messages
    def callback(ch, method, properties, body):
        global message_count, start_time
        
        if start_time is None:
            start_time = time.time()
        
        message = json.loads(body)
        message_count += 1
        
        # Print progress every 1000 messages
        if message_count % 1000 == 0:
            elapsed = time.time() - start_time
            rate = message_count / elapsed
            print(f'Processed {message_count} messages - Rate: {rate:.0f} msgs/sec')
        
        # Acknowledge the message
        ch.basic_ack(delivery_tag=method.delivery_tag)

    # Set up consumer
    channel.basic_qos(prefetch_count=100)  # Increased for better performance
    channel.basic_consume(
        queue='hello_queue',
        on_message_callback=callback,
        auto_ack=False
    )

    print('Waiting for messages. Press CTRL+C to exit.')
    
    # Start consuming messages
    channel.start_consuming()

except KeyboardInterrupt:
    if start_time:
        elapsed = time.time() - start_time
        rate = message_count / elapsed if elapsed > 0 else 0
        print(f'\n\nStopping consumer...')
        print(f'Total messages processed: {message_count}')
        print(f'Total time: {elapsed:.2f} seconds')
        print(f'Average rate: {rate:.0f} messages/second')
    connection.close()
    sys.exit(0)
except Exception as e:
    print(f'Error: {e}')
