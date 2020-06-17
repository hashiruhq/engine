# Research

## Table of contents

## Chapter 1: Introduction into Trading Systems

### Common terms in Finantial markets

**Slippage**
Slippage refers to the difference between the price you expect and the price at which the trade is actually filled. In fast-moving markets, slippage can be substantial and make the difference between a winning and losing trade.

**Long Trade: Profit from a rising market**
Long trades are the classic method of buying with the intention of profiting from a rising market. All brokers support long trades and you won’t need a margin account – assuming you have the funds to cover the trade. Even though losses could be substantial, they are considered limited because price can only go as low as $0 if the trade moves in the wrong direction.

**Short Trade: Profit from a falling market**
Short trades, on the other hand, are entered with the intention of profiting from a falling market. Once price reaches your target level, you buy back the shares (or buy to cover) to replace what you originally borrowed from your broker. Since you borrow shares/contracts from a broker to sell short, you have to have a margin account to complete the transaction.

**Market Order**
A market order is the most basic type of trade order. It’s an order to buy or sell at the best available current price. Most order entry interfaces – including trading apps – have handy "buy" and "sell" buttons to make these orders quick and easy. As long as there is adequate liquidity, this type of order is executed immediately.

**Limit Order**
In a limit order the user can set how much to buy/set and set a maximum/minimum price to fulfill the order on.
Properties:
- amount: the amount you want to buy/sell
- price: the maximum/minimum price to fulfill the order for

**Stop Loss Order**

### References:
- Information about order params: https://docs.gdax.com/#orders



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
