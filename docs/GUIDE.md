# RabbitMQ and Python Integration Guide

This guide demonstrates how RabbitMQ and Python work together in a message queue architecture.

## What is RabbitMQ?

RabbitMQ is a message broker that implements the Advanced Message Queuing Protocol (AMQP). It acts as an intermediary for messaging, allowing applications to communicate with each other by sending and receiving messages through queues.

## Architecture Overview

```
Producer → RabbitMQ Queue → Consumer
```

- **Producer**: Application that sends messages
- **RabbitMQ**: Message broker that stores and routes messages
- **Consumer**: Application that receives and processes messages

## Running the Demo

### Option 1: Local Development (Recommended for beginners)

1. Start RabbitMQ:
   ```bash
   docker compose up -d
   ```

2. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
   ```

3. In terminal 1, start the consumer:
   ```bash
   python consumer.py
   ```

4. In terminal 2, run the producer:
   ```bash
   python producer.py
   ```

### Option 2: Full Docker Environment

Use this if you want everything containerized:

```bash
docker compose -f docker-compose-full.yml up --build
```

## Code Explanation

### Producer (Sender)

The producer script does the following:

1. **Connects to RabbitMQ**:
   ```python
   connection = pika.BlockingConnection(pika.ConnectionParameters(
       host='localhost',
       credentials=pika.PlainCredentials('user', 'password')
   ))
   ```

2. **Creates a channel and declares a queue**:
   ```python
   channel = connection.channel()
   channel.queue_declare(queue='hello_queue', durable=True)
   ```
   - `durable=True` makes the queue survive RabbitMQ restarts

3. **Sends messages**:
   ```python
   channel.basic_publish(
       exchange='',
       routing_key='hello_queue',
       body=json.dumps(message),
       properties=pika.BasicProperties(delivery_mode=2)
   )
   ```
   - `exchange=''` uses the default exchange
   - `delivery_mode=2` makes messages persistent

### Consumer (Receiver)

The consumer script:

1. **Connects to RabbitMQ** (same as producer)

2. **Defines a callback function**:
   ```python
   def callback(ch, method, properties, body):
       message = json.loads(body)
       print(f'Received: {message}')
       ch.basic_ack(delivery_tag=method.delivery_tag)
   ```
   - Processes the message
   - Acknowledges receipt (important for reliability)

3. **Sets up consumption**:
   ```python
   channel.basic_qos(prefetch_count=1)
   channel.basic_consume(
       queue='hello_queue',
       on_message_callback=callback,
       auto_ack=False
   )
   ```
   - `prefetch_count=1` processes one message at a time
   - `auto_ack=False` requires manual acknowledgment

## Key Concepts

### Message Persistence
- Queue durability: `durable=True`
- Message persistence: `delivery_mode=2`
- Ensures messages survive RabbitMQ restarts

### Message Acknowledgment
- Consumer acknowledges message processing
- Prevents message loss if consumer crashes
- RabbitMQ redelivers unacknowledged messages

### Quality of Service (QoS)
- `prefetch_count` controls how many messages a consumer gets
- Helps distribute work evenly among multiple consumers

## Monitoring

Access the RabbitMQ Management UI:
- URL: http://localhost:15673
- Username: user
- Password: password

Here you can:
- View queues and their message counts
- Monitor message rates
- Manage exchanges and bindings
- View consumer connections

## Advanced Topics to Explore

1. **Exchanges**: Different routing patterns (direct, topic, fanout, headers)
2. **Dead Letter Queues**: Handling failed messages
3. **Priority Queues**: Message prioritization
4. **RPC Pattern**: Request-response over RabbitMQ
5. **Clustering**: High availability setup
6. **Federation**: Connecting multiple RabbitMQ instances

## Troubleshooting

### Connection Refused
- Ensure RabbitMQ is running: `docker compose ps`
- Check if ports 5672 and 15672 are available
- Wait a few seconds after starting RabbitMQ

### Messages Not Being Consumed
- Check if consumer is connected in Management UI
- Verify queue names match exactly
- Ensure consumer is acknowledging messages

### Messages Lost on Restart
- Use `durable=True` for queues
- Use `delivery_mode=2` for messages
- Check if volumes are configured in docker compose file

## Best Practices

1. Always declare queues as durable in production
2. Use acknowledgments for reliable message processing
3. Set appropriate prefetch counts for work distribution
4. Handle connection errors and implement retry logic
5. Monitor queue depths and consumer performance
6. Use meaningful queue and exchange names
7. Implement proper error handling and logging

## Next Steps

1. Try modifying the message format
2. Add multiple consumers to see work distribution
3. Experiment with different exchange types
4. Implement error handling and dead letter queues
5. Build a real-world application using these patterns
