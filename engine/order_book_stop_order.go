package engine

/*
 * Copyright © 2006-2019 Around25 SRL <office@around25.com>
 *
 * Licensed under the Around25 Exchange License Agreement (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.around25.com/licenses/EXCHANGE_LICENSE
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Cosmin Harangus <cosmin@around25.com>
 * @copyright 2006-2019 Around25 SRL <office@around25.com>
 * @license 	EXCHANGE_LICENSE
 */

import (
	"gitlab.com/around25/products/matching-engine/model"
)

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


### Stop implementation

Any order that is received by the engine can have a __stop flag__ set along with a __stop price__.
These flags control when the order will be added in the orderbook with the type set at creation: Market or Limit.

The `StopPrice` is __required__ for any Stop order.

There are three types of Stop orders:
- __None__: which means that the order is not a stop order and the it should be processed imediatelly
- __Loss__: which means that the order triggers when the last trade price changes to a value at or below the `StopPrice`
- __Entry__: which means that the order triggers when the last trade price changes to a value at or above the `StopPrice`

Stop orders will only be activated by the next trade.
When it's added no checks will be done based on the previous last trade price.

The processing flow for stop orders is as follows:
1. @done When an order that has the stop flag and price set the order is added to an ordered list based on the flag.
2. @done When a new batch of trades is generated the system will check if the last price can activate any stop order.
3. @done If it can then the found orders are removed from the pending list.
4. @done Every stop order activated is added in the system as limit/market order in the correct side of the market.
5. @done For every activated order an event with order activated is generated by the engine.
6. @done The order is then processed normally by the engine and can generate trades.
7. @todo A user can cancel a stop order before or after it was activated

*/

/**
 * 1. Add stop order as pending
 */

// Process a new order and return a list of events resulted from the exchange
func (book *orderBook) addStopOrder(order model.Order) {
	if order.Stop == model.StopLoss_None || order.StopPrice == 0 {
		return
	}
	switch order.Stop {
	case model.StopLoss_Loss:
		book.StopLossOrders.addOrder(order.StopPrice, order)
		if book.HighestLossPrice == 0 || order.StopPrice > book.HighestLossPrice {
			book.HighestLossPrice = order.StopPrice
		}
	case model.StopLoss_Entry:
		book.StopEntryOrders.addOrder(order.StopPrice, order)
		if book.LowestEntryPrice == 0 || order.StopPrice < book.LowestEntryPrice {
			book.LowestEntryPrice = order.StopPrice
		}
	}
}

/**
 * 2. Activate pending stop orders
 */

func (book *orderBook) GetFirstTradePriceFromEvents(events *[]model.Event) uint64 {
	price := uint64(0)
	if events == nil {
		return price
	}
	eventCount := len(*events)
	if eventCount == 0 {
		return price
	}
	for i := 0; i <= eventCount-1; i++ {
		if (*events)[i].Type == model.EventType_NewTrade {
			trade := (*events)[i].GetTrade()
			price = trade.Price
			break
		}
	}
	return price
}

func (book *orderBook) GetLastTradePriceFromEvents(events *[]model.Event) uint64 {
	lastPrice := uint64(0)
	if events == nil {
		return lastPrice
	}
	eventCount := len(*events)
	if eventCount == 0 {
		return lastPrice
	}
	for i := eventCount - 1; i >= 0; i-- {
		if (*events)[i].Type == model.EventType_NewTrade {
			trade := (*events)[i].GetTrade()
			lastPrice = trade.Price
			break
		}
	}
	return lastPrice
}

// Activate stop orders that have the price
func (book *orderBook) ActivateStopOrders(events *[]model.Event) *[]model.Order {
	orders := &[]model.Order{}
	firstPrice := book.GetFirstTradePriceFromEvents(events)
	lastPrice := book.GetLastTradePriceFromEvents(events)
	min, max := firstPrice, lastPrice
	if firstPrice > lastPrice {
		min = lastPrice
		max = firstPrice
	}
	// activate stop entry orders
	if max != 0 {
		book.activateStopEntryOrders(max, events, orders)
	}

	// activate stop loss price
	if min != 0 {
		book.activateStopLossOrders(min, events, orders)
	}

	return orders
}

//traverse pending stop entry orders and activate the ones that have a stop price lower or equal to the last price
func (book *orderBook) activateStopEntryOrders(lastPrice uint64, events *[]model.Event, orders *[]model.Order) {
	if book.LowestEntryPrice == 0 || book.LowestEntryPrice > lastPrice {
		return
	}
	iterator := book.StopEntryOrders.Seek(book.LowestEntryPrice)
	if iterator == nil {
		return
	}

	for lastPrice >= book.LowestEntryPrice {
		price := iterator.Key()
		pricePoint := iterator.Value()
		*orders = append(*orders, pricePoint.Entries...)
		for _, order := range pricePoint.Entries {
			book.LastEventSeqID++
			*events = append(*events, model.NewOrderActivatedEvent(book.LastEventSeqID, order.Market, order.Type, order.Side, order.ID, order.OwnerID, order.Price, order.Amount, order.Funds, order.Status))
		}
		book.StopEntryOrders.Delete(price)

		if ok := iterator.Next(); ok {
			book.LowestEntryPrice = iterator.Key()
		} else {
			book.LowestEntryPrice = 0
			break
		}
	}
	iterator.Close()
}

//traverse pending stop orders orders and activate the ones that have a stop price higher or equal to the last price
func (book *orderBook) activateStopLossOrders(lastPrice uint64, events *[]model.Event, orders *[]model.Order) {
	if book.HighestLossPrice == 0 || book.HighestLossPrice < lastPrice {
		return
	}
	iterator := book.StopLossOrders.Seek(book.HighestLossPrice)
	if iterator == nil {
		return
	}

	for lastPrice <= book.HighestLossPrice {
		price := iterator.Key()
		pricePoint := iterator.Value()
		*orders = append(*orders, pricePoint.Entries...)
		for _, order := range pricePoint.Entries {
			book.LastEventSeqID++
			*events = append(*events, model.NewOrderActivatedEvent(book.LastEventSeqID, order.Market, order.Type, order.Side, order.ID, order.OwnerID, order.Price, order.Amount, order.Funds, order.Status))
		}
		book.StopLossOrders.Delete(price)

		if ok := iterator.Previous(); ok {
			book.HighestLossPrice = iterator.Key()
		} else {
			book.HighestLossPrice = 0
			break
		}
	}
	iterator.Close()
}

/**
 * 7. Cancel stop order
 */

// Cancel a limit order based on a given order ID and set price
// func (book *orderBook) cancelStopOrder(order model.Order, events *[]model.Event) {
// 	if order.Side == model.MarketSide_Buy {
// 		if pricePoint, ok := book.BuyEntries.Get(order.Price); ok {
// 			for i := 0; i < len(pricePoint.Entries); i++ {
// 				if pricePoint.Entries[i].ID == order.ID {
// 					ord := pricePoint.Entries[i]
// 					ord.SetStatus(model.OrderStatus_Cancelled)
// 					*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
// 					book.removeBuyBookEntry(&ord, pricePoint, i)
// 					// adjust highest bid
// 					if len(pricePoint.Entries) == 0 && book.HighestBid == ord.Price {
// 						iterator := book.BuyEntries.Seek(book.HighestBid)
// 						book.closeBidIterator(iterator)
// 					}
// 					return
// 				}
// 			}
// 		}
// 		return
// 	}
// 	if pricePoint, ok := book.SellEntries.Get(order.Price); ok {
// 		for i := 0; i < len(pricePoint.Entries); i++ {
// 			if pricePoint.Entries[i].ID == order.ID {
// 				ord := pricePoint.Entries[i]
// 				ord.SetStatus(model.OrderStatus_Cancelled)
// 				*events = append(*events, model.NewOrderStatusEvent(book.MarketID, ord.ID, ord.Amount, ord.Funds, ord.Status))
// 				book.removeSellBookEntry(&ord, pricePoint, i)
// 				// adjust lowest ask
// 				if len(pricePoint.Entries) == 0 && book.LowestAsk == ord.Price {
// 					iterator := book.SellEntries.Seek(book.LowestAsk)
// 					book.closeAskIterator(iterator)
// 				}
// 				return
// 			}
// 		}
// 	}
// }
