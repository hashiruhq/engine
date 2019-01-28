package engine

func (book *orderBook) Load(market Market) error {
	book.LowestAsk = market.LowestAsk
	book.HighestBid = market.HighestBid

	for _, buyBookEntry := range market.BuyOrders {
		book.addBuyBookEntry(buyBookEntry)
	}

	for _, sellBookEntry := range market.SellOrders {
		book.addSellBookEntry(sellBookEntry)
	}

	return nil
}

// Backup the order book in another structure for exporting
func (book *orderBook) Backup() Market {
	market := Market{
		LowestAsk:  book.LowestAsk,
		HighestBid: book.HighestBid,
		BuyOrders:  make([]BookEntry, 0, book.BuyEntries.Len()),
		SellOrders: make([]BookEntry, 0, book.SellEntries.Len()),
	}

	if market.LowestAsk != 0 {
		iterator := book.SellEntries.Seek(market.LowestAsk)
		if iterator != nil {
			defer iterator.Close()
			for {
				pricePoint := iterator.Value()
				for _, entry := range pricePoint.Entries {
					market.SellOrders = append(market.SellOrders, entry)
				}
				if ok := iterator.Next(); !ok {
					break
				}
			}
		}
	}

	if market.HighestBid != 0 {
		iterator := book.BuyEntries.Seek(market.HighestBid)
		if iterator != nil {
			defer iterator.Close()
			for {
				pricePoint := iterator.Value()
				for _, entry := range pricePoint.Entries {
					market.BuyOrders = append(market.BuyOrders, entry)
				}
				if ok := iterator.Previous(); !ok {
					break
				}
			}
		}
	}
	return market
}
