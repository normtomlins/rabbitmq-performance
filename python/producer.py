#!/usr/bin/env python3
import pika
import json
import datetime
import time

# Establish connection to RabbitMQ
connection_params = pika.ConnectionParameters(
    host='localhost',
    port=5673,  # Changed from default 5672
    credentials=pika.PlainCredentials('user', 'password')
)

try:
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()

    # Declare a queue (this creates the queue if it doesn't exist)
    channel.queue_declare(queue='hello_queue', durable=True)

    # Send messages
    for i in range(10):
        message = {
            'id': i,
            'content': f'Hello World #{i}',
            'timestamp': datetime.datetime.now().isoformat()
        }
        
        # Publish message
        channel.basic_publish(
            exchange='',  # Using default exchange
            routing_key='hello_queue',  # Queue name
            body=json.dumps(message),
            properties=pika.BasicProperties(
                delivery_mode=2,  # Make message persistent
            )
        )
        
        print(f'Sent: {message}')
        time.sleep(1)  # Wait 1 second between messages

    # Close connection
    connection.close()
    print('All messages sent!')

except pika.exceptions.AMQPConnectionError:
    print('Failed to connect to RabbitMQ. Is it running?')
except Exception as e:
    print(f'Error: {e}')
