# Research 

## Table of contents



## Chapter 1: Introduction into Trading Systems

### Algorithmic trading system architecture

URL: http://www.turingfinance.com/algorithmic-trading-system-architecture-post/

According to the software engineering institute an architectural tactic is a means of satisfying a quality requirement by manipulating some aspect of a quality attribute model through architectural design decisions. A simple example used in the algorithmic trading system architecture is 'manipulating' an operational data store (ODS) with a continuous querying component. This component would continuously analyse the ODS to identify and extract complex events. The following tactics are used in the architecture:

- The disruptor pattern in the event and order queues
- Shared memory for the event and order queues
- Continuous querying language (CQL) on the ODS
- Data filtering with the filter design pattern on incoming data
- Congestion avoidance algorithms on all incoming and outbound connections
- Active queue management (AQM) and explicit congestion notification
- Commodity computing resources with capacity for upgrade (scalable)
- Active redundancy for all single points of failure
- Indexation and optimized persistence structures in the ODS
- Schedule regular data backup and clean-up scripts for ODS
- Transaction histories on all databases
- Checksums for all orders to detect faults
- Annotate events with timestamps to skip 'stale' events
- Order validation rules e.g. maximum trade quantities
- Automated trader components use an in-memory database for analysis
- Two stage authentication for user interfaces connecting to the ATs
- Encryption on user interfaces and connections to the ATs
- Observer design pattern for the MVC to manage views

