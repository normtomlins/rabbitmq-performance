#!/usr/bin/env python3
import pika
import json
import datetime
import time
import os

def connect_to_rabbitmq():
    """Establish connection to RabbitMQ"""
    connection_params = pika.ConnectionParameters(
        host='rabbitmq',
        port=5672,
        credentials=pika.PlainCredentials('user', 'password')
    )
    
    max_retries = 5
    for retry_count in range(max_retries):
        try:
            connection = pika.BlockingConnection(connection_params)
            print(f'Connected to RabbitMQ')
            return connection
        except pika.exceptions.AMQPConnectionError:
            print(f'Connection attempt {retry_count + 1} failed. Retrying in 5 seconds...')
            time.sleep(5)
    
    raise Exception('Failed to connect to RabbitMQ after multiple attempts')

def test_persistent_mode(message_count=1000):
    """Test with persistent messages and durable queue"""
    connection = connect_to_rabbitmq()
    channel = connection.channel()
    channel.queue_declare(queue='persistent_queue', durable=True)
    
    start_time = time.time()
    
    for i in range(message_count):
        message = {
            'id': i,
            'content': f'Hello World #{i}',
            'timestamp': datetime.datetime.now().isoformat()
        }
        
        channel.basic_publish(
            exchange='',
            routing_key='persistent_queue',
            body=json.dumps(message),
            properties=pika.BasicProperties(
                delivery_mode=2,  # Persistent
            )
        )
    
    end_time = time.time()
    total_time = end_time - start_time
    rate = message_count / total_time
    
    connection.close()
    return rate, total_time

def test_non_persistent_mode(message_count=1000):
    """Test with non-persistent messages and non-durable queue"""
    connection = connect_to_rabbitmq()
    channel = connection.channel()
    channel.queue_declare(queue='fast_queue', durable=False)
    
    start_time = time.time()
    
    for i in range(message_count):
        message = {
            'id': i,
            'content': f'Hello World #{i}',
            'timestamp': datetime.datetime.now().isoformat()
        }
        
        channel.basic_publish(
            exchange='',
            routing_key='fast_queue',
            body=json.dumps(message),
            properties=pika.BasicProperties(
                delivery_mode=1,  # Non-persistent
            )
        )
    
    end_time = time.time()
    total_time = end_time - start_time
    rate = message_count / total_time
    
    connection.close()
    return rate, total_time

def test_with_sleep(message_count=1000, sleep_time=0.001):
    """Test with sleep delay"""
    connection = connect_to_rabbitmq()
    channel = connection.channel()
    channel.queue_declare(queue='sleep_queue', durable=True)
    
    start_time = time.time()
    
    for i in range(message_count):
        message = {
            'id': i,
            'content': f'Hello World #{i}',
            'timestamp': datetime.datetime.now().isoformat()
        }
        
        channel.basic_publish(
            exchange='',
            routing_key='sleep_queue',
            body=json.dumps(message),
            properties=pika.BasicProperties(
                delivery_mode=2,  # Persistent
            )
        )
        
        time.sleep(sleep_time)
    
    end_time = time.time()
    total_time = end_time - start_time
    rate = message_count / total_time
    
    connection.close()
    return rate, total_time

if __name__ == "__main__":
    print("RabbitMQ Performance Comparison")
    print("==============================\n")
    
    message_count = 5000
    
    # Test 1: Persistent mode
    print("Test 1: Persistent messages + Durable queue")
    rate, total_time = test_persistent_mode(message_count)
    print(f"Rate: {rate:.0f} msgs/sec")
    print(f"Time: {total_time:.2f} seconds\n")
    
    # Test 2: Non-persistent mode
    print("Test 2: Non-persistent messages + Non-durable queue")
    rate, total_time = test_non_persistent_mode(message_count)
    print(f"Rate: {rate:.0f} msgs/sec")
    print(f"Time: {total_time:.2f} seconds\n")
    
    # Test 3: With sleep (0.001s)
    print("Test 3: Persistent + Sleep(0.001)")
    rate, total_time = test_with_sleep(message_count, 0.001)
    print(f"Rate: {rate:.0f} msgs/sec")
    print(f"Time: {total_time:.2f} seconds\n")
    
    # Test 4: With sleep (0.0001s)
    print("Test 4: Persistent + Sleep(0.0001)")
    rate, total_time = test_with_sleep(message_count, 0.0001)
    print(f"Rate: {rate:.0f} msgs/sec")
    print(f"Time: {total_time:.2f} seconds\n")
