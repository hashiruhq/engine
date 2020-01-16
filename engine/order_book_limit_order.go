package engine

import (
	"gitlab.com/around25/products/matching-engine/model"
)

func (book *orderBook) processLimitBuy(order model.Order, events *[]model.Event) {
	if book.LowestAsk <= order.Price && book.LowestAsk != 0 {
		iterator := book.SellEntries.Seek(book.LowestAsk)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			for order.Price >= book.LowestAsk {
				pricePoint := iterator.Value()
				complete := false
				for index := 0; index < len(pricePoint.Entries); index++ {
					sellEntry := &pricePoint.Entries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					if sellEntry.Amount >= order.Amount {
						// funds := utils.Multiply(order.Amount, sellEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
						book.LastEventSeqID++
						book.LastTradeSeqID++
						*events = append(*events, model.NewTradeEvent(
							book.LastEventSeqID, 
							book.MarketID, 
							book.LastTradeSeqID, 
							model.MarketSide_Buy, 
							sellEntry.ID, 
							order.ID, 
							sellEntry.OwnerID, 
							order.OwnerID, 
							order.Amount, 
							sellEntry.Price,
						))
						sellEntry.Amount -= order.Amount
						order.SetStatus(model.OrderStatus_Filled)
						if sellEntry.Amount == 0 {
							sellEntry.SetStatus(model.OrderStatus_Filled)
							book.removeSellBookEntry(sellEntry.Price, pricePoint, index)
						} else {
							sellEntry.SetStatus(model.OrderStatus_PartiallyFilled)
						}

						complete = true
						break
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order

					// funds := utils.Multiply(sellEntry.Amount, sellEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
					book.LastEventSeqID++
					book.LastTradeSeqID++
					*events = append(*events, model.NewTradeEvent(
						book.LastEventSeqID, 
						book.MarketID, 
						book.LastTradeSeqID, 
						model.MarketSide_Buy, 
						sellEntry.ID, 
						order.ID, 
						sellEntry.OwnerID, 
						order.OwnerID, 
						sellEntry.Amount, 
						sellEntry.Price,
					))
					order.Amount -= sellEntry.Amount
					order.SetStatus(model.OrderStatus_PartiallyFilled)
					sellEntry.SetStatus(model.OrderStatus_Filled)
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
	if book.HighestBid < order.Price || book.HighestBid == 0 {
		book.HighestBid = order.Price
	}
}

func (book *orderBook) processLimitSell(order model.Order, events *[]model.Event) {
	if book.HighestBid >= order.Price && book.HighestBid != 0 {
		iterator := book.BuyEntries.Seek(book.HighestBid)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			for order.Price <= book.HighestBid {
				pricePoint := iterator.Value()
				complete := false
				for index := 0; index < len(pricePoint.Entries); index++ {
					buyEntry := &pricePoint.Entries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					if buyEntry.Amount >= order.Amount {
						// funds := utils.Multiply(order.Amount, buyEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
						book.LastEventSeqID++
						book.LastTradeSeqID++
						*events = append(*events, model.NewTradeEvent(
							book.LastEventSeqID,
							book.MarketID,
							book.LastTradeSeqID,
							model.MarketSide_Sell,
							order.ID,
							buyEntry.ID,
							order.OwnerID,
							buyEntry.OwnerID,
							order.Amount,
							buyEntry.Price,
						))
						buyEntry.Amount -= order.Amount

						order.SetStatus(model.OrderStatus_Filled)
						if buyEntry.Amount == 0 {
							buyEntry.SetStatus(model.OrderStatus_Filled)
							book.removeBuyBookEntry(buyEntry.Price, pricePoint, index)
						} else {
							buyEntry.SetStatus(model.OrderStatus_PartiallyFilled)
						}
						complete = true
						break
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order

					// funds := utils.Multiply(buyEntry.Amount, buyEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
					book.LastEventSeqID++
					book.LastTradeSeqID++
					*events = append(*events, model.NewTradeEvent(
						book.LastEventSeqID,
						book.MarketID,
						book.LastTradeSeqID,
						model.MarketSide_Sell,
						order.ID,
						buyEntry.ID,
						order.OwnerID,
						buyEntry.OwnerID,
						buyEntry.Amount,
						buyEntry.Price,
					))
					order.Amount -= buyEntry.Amount
					order.SetStatus(model.OrderStatus_PartiallyFilled)
					buyEntry.SetStatus(model.OrderStatus_Filled)
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
	if book.LowestAsk > order.Price || book.LowestAsk == 0 {
		book.LowestAsk = order.Price
	}
}
