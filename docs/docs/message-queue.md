# Message Queue Infrastructure

## Overview

The message queue infrastructure plays a pivotal role in ensuring efficient communication between producers and consumers in the system. It facilitates the storage and transfer of messages from producers to consumers, providing a reliable means of asynchronous communication.

## Message Queue Generation and Handling

### Message Queue Creation

Once a connection is established between the producer broker and a remote machine, a message queue is dynamically generated on the remote machine.

### Message Queue Structure

- **Single Queue per Topic:** Each topic corresponds to a single message queue.
- **Multiple Topics per Queue:** However, a single message queue can accommodate messages from multiple topics.

### Message Storage and Relay

- **Message Storage:** Upon receipt, messages from producers are stored in the appropriate message queue based on their associated topic.
- **Relay to Consumers:** Messages stored in the message queue are then relayed to consumers via the consumer channel.

## Process Flow

1. **Message Storage:**
   - Messages from producers are stored in the message queue associated with their respective topics.

2. **Relay to Consumers:**
   - The stored messages are relayed to Load distributor through the consumer channel.