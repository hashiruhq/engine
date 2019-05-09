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

- Done: Fill a new buy/sell order or add it to the list if it can't be filled yet
- Cancel an existing order and remove it from the lists
- Listen for new incomming commands/orders
- Send every trade made to an external system
- Handle over 1.000.000 trades/second
- Done: 100% test converage
- Fully tested with historical data

### Questions/Improvements

- How do we handle server crashes and maintain the trading lists?
- How/When do we acknowledge orders?
- How do we handle persistance while maintaining the high load?
- How do we handle upgrades and redeployments?
- How do we prevent double spending?
- How do we prevent memory leaks?
- Perform benchmarks on the entire system and check what happens when the load is increased.

## Benchmarks

### Process orders rates

This bechmark checks the performance of executing orders after they are available in memory.

```
// Old benchmarks

Total Orders: 1000000
Total Trades: 996903
Orders/second: 667924.145191
Trades/second: 665855.584113
Pending Buy: 2838
Lowest Ask: 9940629
Pending Sell: 2838
Highest Bid: 7007361
Duration (seconds): 1.497176

1000000	      1497 ns/op	     642 B/op	       6 allocs/op

// New Benchmarks without decoding/encoding orders/trades
Total Orders: 100
Orders/second: 980392.156863
Pending Buy: 18
Lowest Ask: 999705600000000
Pending Sell: 10
Highest Bid: 999622700000000
Duration (seconds): 0.000102

Total Orders: 10000
Orders/second: 1511944.360448
Pending Buy: 591
Lowest Ask: 996920900000000
Pending Sell: 662
Highest Bid: 996859500000000
Duration (seconds): 0.006614

Total Orders: 1000000
Orders/second: 1160515.547427
Pending Buy: 2244
Lowest Ask: 996670000000000
Pending Sell: 1038
Highest Bid: 700948400000000
Duration (seconds): 0.861686

Total Orders: 2000000
Orders/second: 1058599.293173
Pending Buy: 4207
Lowest Ask: 994008100000000
Pending Sell: 1350
Highest Bid: 402202500000000
Duration (seconds): 1.889289

 2000000	       944 ns/op	     505 B/op	       3 allocs/op

// New Benchmarks with decoding/encoding orders/trades

Total Orders: 100
Orders/second: 735294.117647
Pending Buy: 18
Lowest Ask: 999705600000000
Pending Sell: 10
Highest Bid: 999622700000000
Duration (seconds): 0.000136

Total Orders: 10000
Orders/second: 831393.415364
Pending Buy: 591
Lowest Ask: 996920900000000
Pending Sell: 662
Highest Bid: 996859500000000
Duration (seconds): 0.012028

Total Orders: 1000000
Orders/second: 600300.990917
Pending Buy: 2244
Lowest Ask: 996670000000000
Pending Sell: 1038
Highest Bid: 700948400000000
Duration (seconds): 1.665831

 1000000	      1665 ns/op	     698 B/op	       7 allocs/op
```

### Consumer and Producer

The following benchmark refers to consuming orders from a Kafka server, processing them
by the trading engine and saving generated trades back to the Kafka server.

```
Total Orders: 300000
Total Trades Generated: 295126
Orders/second: 248929.396155
Trades Generated/second: 244885.123232
Pending Buy: 3970
Lowest Ask: 3967957
Pending Sell: 3970
Highest Bid: 3706215
Duration (seconds): 1.205161

 300000          4043 ns/op       2437 B/op       22 allocs/op
```

## Further optimizations

- Add support for FIX formats
- Maybe Switch Kafka Client to: https://github.com/confluentinc/confluent-kafka-go
- Optimize Kafka Producer rate
- Use the golang profiler to see where the bottlenecks are in the code
- Write to multiple partitions at once when generating trades
- Optimize allocations with buffers bytes (see minute 23 of the video): https://www.youtube.com/watch?v=ZuQcbqYK0BY

## Resources

**Architecture Design**

- https://martinfowler.com/articles/lmax.html
- https://web.archive.org/web/20110314042933/http://howtohft.wordpress.com/
- http://blog.bitfinex.com/announcements/introducing-hive/
- https://web.archive.org/web/20110310171841/http://www.quantcup.org/home/spec
- https://web.archive.org/web/20110315000556/http://drdobbs.com:80/high-performance-computing/210604448

**Trading Engine**

- - Matcher: https://github.com/fmstephe/matching_engine/blob/master/matcher/matcher.go
- Order book: https://github.com/hookercookerman/trading_engine/blob/master/trading_engine/order_book.go
- https://github.com/bloq/cpptrade#summary
- - https://www.investopedia.com/university/intro-to-order-types/
- - https://github.com/IanLKaplan/matchingEngine/wiki/Market-Order-Matching-Engines
- - https://medium.com/@Oskarr3/implementing-cqrs-using-kafka-and-sarama-library-in-golang-da7efa3b77fe
- https://hackernoon.com/a-blockchain-experiment-with-apache-kafka-97ee0ab6aefc

**Golang Docs**

- https://www.apress.com/gp/book/9781484226919

**Communication**

- Apache Kafka: https://godoc.org/github.com/Shopify/sarama
- Alternative Kafka Client: https://github.com/confluentinc/confluent-kafka-go
- JSON Decoder IO Reader: https://golang.org/pkg/encoding/json/#Decoder.Decode
- Unix Domain Socket: https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2

**Performance**

- Great presentation on how to use go profilling tools: https://www.youtube.com/watch?v=N3PWzBeLX2M
- LMAX disruptor: https://lmax-exchange.github.io/disruptor/
  - https://dzone.com/articles/using-apache-kafka-for-real-time-event-processing
  - https://www.slideshare.net/InfoQ/low-latency-trading-architecture-at-lmax-exchange
- https://www.youtube.com/watch?v=DJ4d_PZ6Gns
  - Use a ring buffer instead of channels
  - Understand the internals of the code you call
  -
- https://blog.golang.org/profiling-go-programs
- https://github.com/golang/go/wiki/Performance
- https://www.quora.com/How-does-one-become-a-low-latency-programmer
- https://www.akkadia.org/drepper/cpumemory.pdf
- https://engineering.linkedin.com/kafka/benchmarking-apache-kafka-2-million-writes-second-three-cheap-machines
- FS buffering: https://zensrc.wordpress.com/2016/06/15/golang-buffer-vs-non-buffer-in-io/
- Fwd: https://blog.kintoandar.com/2016/08/fwd-the-little-forwarder-that-could.html
- Effective kafka topic partitioning: https://blog.newrelic.com/2018/03/13/effective-strategies-kafka-topic-partitioning/

**Performance: To Read**

- http://bhomnick.net/building-a-simple-limit-order-in-go/
- https://github.com/ProofSuite/amp-matching-engine
- https://github.com/robaho/go-trader
- https://github.com/i25959341/orderbook
- https://github.com/crypto-bank

**Algorithms**

- http://igoro.com/archive/skip-lists-are-fascinating/

**Liquidity**

- https://github.com/HydroProtocol/amm-bots

## Testing

`go test -bench=^BenchmarkWithRandomData -run=^$ -timeout 10s -benchmem -cpu 8 -cpuprofile cpu.prof -memprofile mem.prof -gcflags="-m" ./benchmarks`
`go-torch benchmarks.test cpu.prof`

**Create flame chart out of active server**

```
go-torch -u http://127.0.0.1:6060
```

**Start GoConvey**

```
govendor install
goconvey .
```

### Testing References

**Benchmarking**

- https://godoc.org/github.com/codahale/hdrhistogram
- https://github.com/tylertreat/bench

**Testing**

- Mock sarama producers: https://medium.com/@Oskarr3/implementing-cqrs-using-kafka-and-sarama-library-in-golang-da7efa3b77fe
