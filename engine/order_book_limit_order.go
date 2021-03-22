package engine

import (
	"gitlab.com/around25/products/matching-engine/model"
	"gitlab.com/around25/products/matching-engine/utils"
)

// append a new order status to the list of events with the current state of the order
func (book *orderBook) appendOrderStatusEvent(events *[]model.Event, order model.Order) {
	book.LastEventSeqID++
	*events = append(*events, model.NewOrderStatusEvent(
		book.LastEventSeqID,
		order.Market,
		order.Type,
		order.Side,
		order.ID,
		order.OwnerID,
		order.Price,
		order.Amount,
		order.Funds,
		order.Status,
		order.FilledAmount,
		order.UsedFunds,
	))
}

func (book *orderBook) appendTradeEvent(events *[]model.Event, order, existingOrder model.Order, amount uint64) {
	book.LastEventSeqID++
	book.LastTradeSeqID++
	if order.Side == model.MarketSide_Buy {
		*events = append(*events, model.NewTradeEvent(
			book.LastEventSeqID,
			order.Market,
			book.LastTradeSeqID,
			order.Side,
			existingOrder.ID,
			order.ID,
			existingOrder.OwnerID,
			order.OwnerID,
			amount,
			existingOrder.Price,
		))
	} else {
		*events = append(*events, model.NewTradeEvent(
			book.LastEventSeqID,
			order.Market,
			book.LastTradeSeqID,
			order.Side,
			order.ID,
			existingOrder.ID,
			order.OwnerID,
			existingOrder.OwnerID,
			amount,
			existingOrder.Price,
		))
	}
}

func (book *orderBook) processLimitBuy(order model.Order, events *[]model.Event) {
	if book.LowestAsk <= order.Price && book.LowestAsk != 0 {
		iterator := book.SellEntries.Seek(book.LowestAsk)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			for order.Price >= iterator.Key() {
				pricePoint := iterator.Value()
				complete := false
				for index := 0; index < len(pricePoint.Entries); index++ {
					sellEntry := &pricePoint.Entries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					orderUnfilledAmount := order.GetUnfilledAmount()
					sellEntryUnfilledAmount := sellEntry.GetUnfilledAmount()
					if sellEntryUnfilledAmount >= orderUnfilledAmount {
						funds := utils.Multiply(orderUnfilledAmount, sellEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
						book.appendTradeEvent(events, order, *sellEntry, orderUnfilledAmount)
						sellEntry.FilledAmount += orderUnfilledAmount
						sellEntry.UsedFunds += funds
						order.FilledAmount += orderUnfilledAmount
						order.UsedFunds += funds
						order.SetStatus(model.OrderStatus_Filled)

						// Add updates to the events for the filled orders
						book.appendOrderStatusEvent(events, order) // order is filled

						if sellEntry.GetUnfilledAmount() == 0 {
							sellEntry.SetStatus(model.OrderStatus_Filled)
							book.appendOrderStatusEvent(events, *sellEntry) // order is filled or partially filled

							book.removeSellBookEntry(sellEntry.Price, pricePoint, index)
						} else {
							sellEntry.SetStatus(model.OrderStatus_PartiallyFilled)
							book.appendOrderStatusEvent(events, *sellEntry) // order is filled or partially filled
						}

						complete = true
						break
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order
					funds := utils.Multiply(sellEntryUnfilledAmount, sellEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
					book.appendTradeEvent(events, order, *sellEntry, sellEntryUnfilledAmount)
					order.FilledAmount += sellEntryUnfilledAmount
					order.UsedFunds += funds
					order.SetStatus(model.OrderStatus_PartiallyFilled)
					sellEntry.FilledAmount += sellEntryUnfilledAmount
					sellEntry.UsedFunds += funds
					sellEntry.SetStatus(model.OrderStatus_Filled)

					// Add updates to the events for the filled orders
					book.appendOrderStatusEvent(events, *sellEntry) // order is filled

					book.removeSellBookEntry(sellEntry.Price, pricePoint, index)
					index--
				}

				if complete {
					book.closeAskIterator(iterator)
					return
				}

				if ok := iterator.Next(); ok {
					book.LowestAsk = iterator.Key()
				} else {
					book.LowestAsk = 0
					break
				}
			}

			iterator.Close()
		}
	}

	// if there are no more orders just add the buy order to the list
	book.addBuyBookEntry(order)
	// Add updates to the events for the added order
	if order.Status != model.OrderStatus_Untouched {
		book.appendOrderStatusEvent(events, order) // order is partially filled
	}

	if book.HighestBid < order.Price || book.HighestBid == 0 {
		book.HighestBid = order.Price
	}
}

func (book *orderBook) processLimitSell(order model.Order, events *[]model.Event) {
	if book.HighestBid >= order.Price && book.HighestBid != 0 {
		iterator := book.BuyEntries.Seek(book.HighestBid)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			for order.Price <= iterator.Key() {
				pricePoint := iterator.Value()
				complete := false
				for index := 0; index < len(pricePoint.Entries); index++ {
					buyEntry := &pricePoint.Entries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					orderUnfilledAmount := order.GetUnfilledAmount()
					buyEntryUnfilledAmount := buyEntry.GetUnfilledAmount()
					if buyEntryUnfilledAmount >= orderUnfilledAmount {
						funds := utils.Multiply(orderUnfilledAmount, buyEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
						book.appendTradeEvent(events, order, *buyEntry, orderUnfilledAmount)
						buyEntry.FilledAmount += orderUnfilledAmount
						buyEntry.UsedFunds += funds
						order.FilledAmount += orderUnfilledAmount
						order.UsedFunds += funds
						order.SetStatus(model.OrderStatus_Filled)
						// Add updates to the events for the filled orders
						book.appendOrderStatusEvent(events, order) // order is filled

						if buyEntry.GetUnfilledAmount() == 0 {
							buyEntry.SetStatus(model.OrderStatus_Filled)
							book.appendOrderStatusEvent(events, *buyEntry) // order is filled or partially filled
							book.removeBuyBookEntry(buyEntry.Price, pricePoint, index)
						} else {
							buyEntry.SetStatus(model.OrderStatus_PartiallyFilled)
							book.appendOrderStatusEvent(events, *buyEntry) // order is filled or partially filled
						}

						complete = true
						break
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order

					funds := utils.Multiply(buyEntryUnfilledAmount, buyEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
					book.appendTradeEvent(events, order, *buyEntry, buyEntryUnfilledAmount)
					order.FilledAmount += buyEntryUnfilledAmount
					order.UsedFunds += funds
					order.SetStatus(model.OrderStatus_PartiallyFilled)
					buyEntry.FilledAmount += buyEntryUnfilledAmount
					buyEntry.UsedFunds += funds
					buyEntry.SetStatus(model.OrderStatus_Filled)
					book.appendOrderStatusEvent(events, *buyEntry) // order is filled
					book.removeBuyBookEntry(buyEntry.Price, pricePoint, index)
					index--
				}

				if complete {
					book.closeBidIterator(iterator)
					return
				}

				if ok := iterator.Previous(); ok {
					book.HighestBid = iterator.Key()
				} else {
					book.HighestBid = 0
					break
				}
			}
			iterator.Close()
		}
	}

	// if there are no more orders just add the buy order to the list
	book.addSellBookEntry(order)

	// Add updates to the events for the added order
	if order.Status != model.OrderStatus_Untouched {
		book.appendOrderStatusEvent(events, order) // order is partially filled
	}

	if book.LowestAsk > order.Price || book.LowestAsk == 0 {
		book.LowestAsk = order.Price
	}
}
