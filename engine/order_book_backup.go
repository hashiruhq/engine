package engine

func (book *orderBook) Load(market MarketBackup) error {
	book.LowestAsk = market.LowestAsk
	book.HighestBid = market.HighestBid

	for _, buyBookEntry := range market.BuyOrders {
		book.addBuyBookEntry(*buyBookEntry)
	}

	for _, sellBookEntry := range market.SellOrders {
		book.addSellBookEntry(*sellBookEntry)
	}

	return nil
}

// Backup the order book in another structure for exporting
func (book *orderBook) Backup() MarketBackup {
	market := MarketBackup{
		LowestAsk:  book.LowestAsk,
		HighestBid: book.HighestBid,
		BuyOrders:  make([]*Order, 0, book.BuyEntries.Len()),
		SellOrders: make([]*Order, 0, book.SellEntries.Len()),
	}

	if market.LowestAsk != 0 {
		iterator := book.SellEntries.Seek(market.LowestAsk)
		if iterator != nil {
			for {
				pricePoint := iterator.Value()
				for _, entry := range pricePoint.Entries {
					market.SellOrders = append(market.SellOrders, &entry)
				}
				if ok := iterator.Next(); !ok {
					break
				}
			}
			iterator.Close()
		}
	}

	if market.HighestBid != 0 {
		iterator := book.BuyEntries.Seek(market.HighestBid)
		if iterator != nil {
			for {
				pricePoint := iterator.Value()
				for _, entry := range pricePoint.Entries {
					market.BuyOrders = append(market.BuyOrders, &entry)
				}
				if ok := iterator.Previous(); !ok {
					break
				}
			}
			iterator.Close()
		}
	}
	return market
}
