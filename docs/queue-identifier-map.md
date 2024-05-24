# QIM Structure

## Overview

The Queue Identifier Map (QIM) is a critical component in the message transfer process between the producer broker and remote machines. It uniquely identifies a message queue and contains essential information for establishing communication and transferring messages.

## Structure of QIM

The QIM consists of two main parts: the header and the body.

### 1. Header

The header contains the following information about the remote machine:

- **IP Address:** The IP address of the remote machine.
- **Additional Metadata:** Any other relevant information required to identify and establish a connection with the remote machine.

### 2. Body

The body includes details necessary for the message transfer:

- **Protocol:** The protocol to be used for transferring the message (e.g., HTTP, RPC, etc.).
- **Topic Name:** The name of the topic associated with the message queue.

## Process Flow

1. **QIM Request:**
   - When a new producer needs to be registered, the producer broker requests a QIM from a remote machine using an available IP address.
   
2. **QIM Generation:**
   - The remote machine generates a new QIM, including a header with its IP address and metadata, and a body containing the transfer protocol and topic name.
   
3. **QIM Reception:**
   - The generated QIM is sent back to the producer broker, which then updates its Producer Broker - QIM map with the new entry.

4. **Message Transfer:**
   - The producer broker uses the QIM to determine the protocol and topic name.
   - A logical connection is established with the remote machine based on the QIM header.
   - The message is pushed to the message queue as specified in the QIM body.
