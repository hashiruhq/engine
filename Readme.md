# Trading Engine

## How to build a trading engine
In this section I will try to document step by step how to build a trading engine in Golang from scratch.

### Orders 
Every trading engine needs to listen for orders from some sort of external system. This can be a web socket, a queue system, a database, etc.
The point is that we need a way to encapsulate each type of order received.

The trading engine will therefore keep a list of trade orders and match incomming orders against that list based on the type of order that is received.

For example if there is an unprocessed buy order of 200 for 0.41 and an incomming order of sell for 100 for 0.40 then the incomming order is satisfied on the spot, the buy order is decreased by 100 and a buy signal is triggerred.

Such an order may look like the following data structure:
```js
{
  action: 0, // 0=sell, 1=buy, 2=cancel-sell, 3=cancel-buy
  category: 0, // 0=market, 1=limit, 2=stop-loss, 3=... // see https://www.investopedia.com/university/intro-to-order-types/ for other possible types
  price: 100,
  amount: 124,
  id: 51987123,
}
```

### Engine

The engine will maintain two lists of orders. Buy orders and Sell orders.
It should then listen for new orders and try to satisfy them against these lists.

The engine should support the following actions:
- Fill a new buy/sell order or add it to the list if it can't be filled yet
- Cancel an existing order and remove it from the lists
- Listen for new incomming commands/orders
- Send every trade made to an external system
- Handle over 1.000.000 trades/second
- 100% test converage
- Fully tested with historical data

### Questions/Improvements
- How do we handle server crashes and maintain the trading lists?
- How do we order acknowledgement?
- How do we handle persistance while maintaining the high load?
- How do we handle upgrades and redeployments?
- How do we prevent double spending?
- How do we prevent memory leaks?
- Perform benchmarks on the entire system and check what happens when the load is increased.

## Benchmarks

### Process orders rates

This bechmark checks the performance of executing orders after they are available in memory.

> Total Orders: 1000000\
> Total Trades: 997065\
> Orders/second: 614837.504596\
> Trades/second: 613032.956520\
> Pending Buy: 3022\
> Lowest Ask: 7950529.000000\
> Pending Sell: 3022\
> Highest Bid: 7007034.000000\
> Duration (seconds): 1.626446\
> \
> 1000000 &nbsp; &nbsp; &nbsp; 1626 ns/op &nbsp; &nbsp; &nbsp; 597 B/op &nbsp; &nbsp; &nbsp; 7 allocs/op

### Consumer and Producer

The following benchmark refers to consuming orders from a Kafka server, processing them 
by the trading engine and saving generated trades back to the Kafka server.

> Total Orders: 300000 \
> Total Trades Generated: 293410 \
> Orders/second: 207613.031455\
> Trades Generated/second: 203052.465197\
> Pending Buy: 5117\
> Lowest Ask: 3867492.000000\
> Pending Sell: 5117\
> Highest Bid: 3735132.000000\
> Duration (seconds): 1.444996\
> \
> 300000 &nbsp; &nbsp; &nbsp; 4840 ns/op &nbsp; &nbsp; &nbsp; 2583 B/op &nbsp; &nbsp; &nbsp; 25 allocs/op

## Further optimizations
- Decrease the size of the Kafka message value by switching the JSON encoded order to either Protobufs or FIX formats
- Switch Kafka Client to: https://github.com/confluentinc/confluent-kafka-go
- Optimize Kafka Producer rate
- Use the golang profiler to see where the bottlenecks are in the code
- Write to multiple partitions at once when generating trades
- Optimize allocations with buffers bytes (see minute 23 of the video): https://www.youtube.com/watch?v=ZuQcbqYK0BY

## Resources

Trading Engine
- Matcher: https://github.com/fmstephe/matching_engine/blob/master/matcher/matcher.go
- Order book: https://github.com/hookercookerman/trading_engine/blob/master/trading_engine/order_book.go
- https://github.com/bloq/cpptrade#summary
- https://www.investopedia.com/university/intro-to-order-types/ 

Golang Docs
- https://gobyexample.com/goroutines
- https://www.apress.com/gp/book/9781484226919

Communication
- Apache Kafka: https://godoc.org/github.com/Shopify/sarama
- Alternative Kafka Client: https://github.com/confluentinc/confluent-kafka-go
- JSON Decoder IO Reader: https://golang.org/pkg/encoding/json/#Decoder.Decode
- Unix Domain Socket: https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2

Performance
- https://blog.golang.org/profiling-go-programs
- https://github.com/golang/go/wiki/Performance
- https://www.quora.com/How-does-one-become-a-low-latency-programmer
- https://www.akkadia.org/drepper/cpumemory.pdf
- https://engineering.linkedin.com/kafka/benchmarking-apache-kafka-2-million-writes-second-three-cheap-machines
- FS buffering: https://zensrc.wordpress.com/2016/06/15/golang-buffer-vs-non-buffer-in-io/
- Fwd: https://blog.kintoandar.com/2016/08/fwd-the-little-forwarder-that-could.html


Algorithms
- http://igoro.com/archive/skip-lists-are-fascinating/