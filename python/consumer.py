#!/usr/bin/env python3
import pika
import json
import sys

# Establish connection to RabbitMQ
connection_params = pika.ConnectionParameters(
    host='localhost',
    port=5673,  # Changed from default 5672
    credentials=pika.PlainCredentials('user', 'password')
)

try:
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()

    # Declare the same queue (idempotent operation)
    channel.queue_declare(queue='hello_queue', durable=True)

    # Define callback function to process messages
    def callback(ch, method, properties, body):
        message = json.loads(body)
        print(f'Received: {message}')
        
        # Simulate some work
        print(f'Processing message ID: {message["id"]}')
        
        # Acknowledge the message was received and processed
        ch.basic_ack(delivery_tag=method.delivery_tag)

    # Set up consumer
    channel.basic_qos(prefetch_count=1)  # Process one message at a time
    channel.basic_consume(
        queue='hello_queue',
        on_message_callback=callback,
        auto_ack=False  # Manual acknowledgment
    )

    print('Waiting for messages. Press CTRL+C to exit.')
    
    # Start consuming messages
    channel.start_consuming()

except KeyboardInterrupt:
    print('\nStopping consumer...')
    connection.close()
    sys.exit(0)
except pika.exceptions.AMQPConnectionError:
    print('Failed to connect to RabbitMQ. Is it running?')
except Exception as e:
    print(f'Error: {e}')
