package engine

/**

StopLoss Orders
===============

This file contains methods that help with the execution of stop loss orders.
An order is considered a stop loss order if the Stop and StopPrice fields are set.

The order will be processed either as Market Order or Limit Order but the order will not be active
until the StopPrice set on the order is reached.

There are 2 types of stop orders: 0=none 1=loss 2=entry
 - Stop loss triggers when the last trade price changes to a value at or below the `StopPrice`.
 - Stop entry triggers when the last trade price changes to a value at or above the `StopPrice`.
 - Note that when triggered, stop orders execute as either market or limit orders, depending on the type.

When a new order is added the matching engine will check if that order has the stop flag set to either loss or entry.
If it's set then that order will be added in an ordered list and wait for the the StopPrice to be reached.

If the last trade price changes to a value at or below for loss triggers and at or above for entry triggers then
that order will be processed by the matching engine as either a Market order or as a Limit order depending on the
specified order type and the side of the market for which it was added.

*/

// @todo add a method that receives the last trade of the previous order and checks if there are any
// pending stop loss orders that should be activated at that price.

// @todo if there are then those orders will be processed and added as limit/market buy/sell orders in the order book
// From that point on there are considered normal orders and follow the same flow as other types.
