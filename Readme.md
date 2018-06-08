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
  type: 0, // 0=market, 1=limit, 2=stop-loss, 3=... // see https://www.investopedia.com/university/intro-to-order-types/ for other possible types
  price: 100,
  amount: 124,
  order_id: 51987123,
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

## Resources

Trading Engine
- Matcher: https://github.com/fmstephe/matching_engine/blob/master/matcher/matcher.go
- Order book: https://github.com/hookercookerman/trading_engine/blob/master/trading_engine/order_book.go
- https://github.com/bloq/cpptrade#summary
- https://www.investopedia.com/university/intro-to-order-types/ 

Golang Docs
- https://gobyexample.com/goroutines
- https://www.apress.com/gp/book/9781484226919