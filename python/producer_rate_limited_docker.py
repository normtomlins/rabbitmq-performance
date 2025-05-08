#!/usr/bin/env python3
import pika
import json
import datetime
import time
import os

# Establish connection to RabbitMQ
connection_params = pika.ConnectionParameters(
    host='rabbitmq',
    port=5672,
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

try:
    channel = connection.channel()
    channel.queue_declare(queue='hello_queue', durable=True)

    # Send messages at a controlled rate
    message_count = 10000
    target_rate = 1000  # messages per second
    start_time = time.time()
    
    print(f'Starting to send {message_count} messages at {target_rate} msgs/sec...')
    
    for i in range(message_count):
        loop_start = time.time()
        
        message = {
            'id': i,
            'content': f'Hello World #{i}',
            'timestamp': datetime.datetime.now().isoformat()
        }
        
        # Publish message
        channel.basic_publish(
            exchange='',
            routing_key='hello_queue',
            body=json.dumps(message),
            properties=pika.BasicProperties(
                delivery_mode=2,  # Persistent
            )
        )
        
        # Print progress every 1000 messages
        if i > 0 and i % 1000 == 0:
            elapsed = time.time() - start_time
            rate = i / elapsed
            print(f'Sent {i} messages - Rate: {rate:.0f} msgs/sec')
        
        # Rate limiting logic
        elapsed_time = time.time() - start_time
        expected_time = (i + 1) / target_rate
        sleep_time = expected_time - elapsed_time
        
        if sleep_time > 0:
            time.sleep(sleep_time)

    # Print final statistics
    end_time = time.time()
    total_time = end_time - start_time
    actual_rate = message_count / total_time
    
    connection.close()
    print(f'\nCompleted!')
    print(f'Total messages: {message_count}')
    print(f'Total time: {total_time:.2f} seconds')
    print(f'Average rate: {actual_rate:.0f} messages/second')

except Exception as e:
    print(f'Error: {e}')
