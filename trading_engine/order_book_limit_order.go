package trading_engine

func (book *orderBook) processLimitBuy(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	if book.LowestAsk <= order.Price && book.LowestAsk != 0 {
		iterator := book.PricePoints.Seek(book.LowestAsk)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			defer iterator.Close()
			for order.Price >= book.LowestAsk {
				pricePoint := iterator.Value().(*PricePoint)
				for index := 0; index < len(pricePoint.SellBookEntries); index++ {
					sellEntry := &pricePoint.SellBookEntries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					if sellEntry.Amount >= order.Amount {
						trades = append(trades, NewTrade(order.ID, sellEntry.Order.ID, order.Amount, sellEntry.Order.Price))
						sellEntry.Amount -= order.Amount
						if sellEntry.Amount == 0 {
							// book.removeBookEntry(sellEntry)
							book.removeSellBookEntry(sellEntry, pricePoint, index)
						}
						return trades
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order
					if sellEntry.Amount < order.Amount {
						trades = append(trades, NewTrade(order.ID, sellEntry.Order.ID, sellEntry.Amount, sellEntry.Price))
						order.Amount -= sellEntry.Amount
						// book.removeBookEntry(sellEntry)
						book.removeSellBookEntry(sellEntry, pricePoint, index)
						index--
						continue
					}
				}

				if ok := iterator.Next(); ok {
					if len(iterator.Value().(*PricePoint).SellBookEntries) > 0 {
						book.LowestAsk = iterator.Key().(uint64)
					}
				} else {
					book.LowestAsk = 0
					break
				}
			}
		}
	}

	// if there are no more ordes just add the buy order to the list
	book.addBookEntry(BookEntry{Price: order.Price, Amount: order.Amount, Order: order})
	if book.HighestBid < order.Price || book.HighestBid == 0 {
		book.HighestBid = order.Price
	}

	return trades
}

func (book *orderBook) processLimitSell(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	if book.HighestBid >= order.Price && book.HighestBid != 0 {
		iterator := book.PricePoints.Seek(book.HighestBid)

		// traverse orders to find a matching one based on the sell order list
		if iterator != nil {
			defer iterator.Close()
			for order.Price <= book.HighestBid {
				pricePoint := iterator.Value().(*PricePoint)
				for index := 0; index < len(pricePoint.BuyBookEntries); index++ {
					buyEntry := &pricePoint.BuyBookEntries[index]
					// if we can fill the trade instantly then we add the trade and complete the order
					if buyEntry.Amount >= order.Amount {
						trades = append(trades, NewTrade(order.ID, buyEntry.Order.ID, order.Amount, buyEntry.Order.Price))
						buyEntry.Amount -= order.Amount
						if buyEntry.Amount == 0 {
							book.removeBuyBookEntry(buyEntry, pricePoint, index)
							// book.removeBookEntry(buyEntry)
						}
						return trades
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order
					if buyEntry.Amount < order.Amount {
						trades = append(trades, NewTrade(order.ID, buyEntry.Order.ID, buyEntry.Amount, buyEntry.Price))
						order.Amount -= buyEntry.Amount
						// book.removeBookEntry(buyEntry)
						book.removeBuyBookEntry(buyEntry, pricePoint, index)
						index--
						continue
					}
				}

				if ok := iterator.Previous(); ok {
					if len(iterator.Value().(*PricePoint).BuyBookEntries) > 0 {
						book.HighestBid = iterator.Key().(uint64)
					}
				} else {
					book.HighestBid = 0
					break
				}
			}
		}
	}

	// if there are no more ordes just add the buy order to the list
	book.addBookEntry(BookEntry{Price: order.Price, Amount: order.Amount, Order: order})
	if book.LowestAsk > order.Price || book.LowestAsk == 0 {
		book.LowestAsk = order.Price
	}

	return trades
}
