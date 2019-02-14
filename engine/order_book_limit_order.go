package engine

func (book *orderBook) processLimitBuy(order Order, trades *[]Trade) {
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
						*trades = append(*trades, NewTrade(book.MarketID, MarketSide_Sell, sellEntry.ID, order.ID, sellEntry.OwnerID, order.OwnerID, order.Amount, sellEntry.Price))
						sellEntry.Amount -= order.Amount
						if sellEntry.Amount == 0 {
							book.removeSellBookEntry(sellEntry, pricePoint, index)
						}

						complete = true
						break
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order
					if sellEntry.Amount < order.Amount {
						// funds := utils.Multiply(sellEntry.Amount, sellEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
						*trades = append(*trades, NewTrade(book.MarketID, MarketSide_Sell, sellEntry.ID, order.ID, sellEntry.OwnerID, order.OwnerID, sellEntry.Amount, sellEntry.Price))
						order.Amount -= sellEntry.Amount
						book.removeSellBookEntry(sellEntry, pricePoint, index)
						index--
						continue
					}
				}

				if complete {
					if len(pricePoint.Entries) != 0 {
						iterator.Close()
						return
					}
					if ok := iterator.Next(); ok {
						book.LowestAsk = iterator.Key()
						iterator.Close()
						return
					}
					book.LowestAsk = 0
					iterator.Close()
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

	// if there are no more ordes just add the buy order to the list
	book.addBuyBookEntry(order)
	if book.HighestBid < order.Price || book.HighestBid == 0 {
		book.HighestBid = order.Price
	}
}

func (book *orderBook) processLimitSell(order Order, trades *[]Trade) {
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
						*trades = append(*trades, NewTrade(book.MarketID, MarketSide_Sell, order.ID, buyEntry.ID, order.OwnerID, buyEntry.OwnerID, order.Amount, buyEntry.Price))
						buyEntry.Amount -= order.Amount
						if buyEntry.Amount == 0 {
							book.removeBuyBookEntry(buyEntry, pricePoint, index)
						}
						complete = true
						break
					}

					// if the sell order has a lower amount than what the buy order is then we fill only what we can from the sell order,
					// we complete the sell order and we move to the next order
					if buyEntry.Amount < order.Amount {
						// funds := utils.Multiply(buyEntry.Amount, buyEntry.Price, book.VolumePrecision, book.PricePrecision, book.PricePrecision)
						*trades = append(*trades, NewTrade(book.MarketID, MarketSide_Sell, order.ID, buyEntry.ID, order.OwnerID, buyEntry.OwnerID, buyEntry.Amount, buyEntry.Price))
						order.Amount -= buyEntry.Amount
						book.removeBuyBookEntry(buyEntry, pricePoint, index)
						index--
						continue
					}
				}

				if complete {
					if len(pricePoint.Entries) != 0 {
						iterator.Close()
						return
					}
					if ok := iterator.Previous(); ok {
						book.HighestBid = iterator.Key()
						iterator.Close()
						return
					}
					book.HighestBid = 0
					iterator.Close()
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

	// if there are no more ordes just add the buy order to the list
	book.addSellBookEntry(order)
	if book.LowestAsk > order.Price || book.LowestAsk == 0 {
		book.LowestAsk = order.Price
	}
}