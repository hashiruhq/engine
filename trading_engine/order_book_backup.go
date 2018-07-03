package trading_engine

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
	pricePoints := book.PricePoints
	market := Market{
		LowestAsk:  book.LowestAsk,
		HighestBid: book.HighestBid,
		BuyOrders:  make([]BookEntry, 0, pricePoints.Len()),
		SellOrders: make([]BookEntry, 0, pricePoints.Len()),
	}

	if market.LowestAsk != 0 {
		iterator := pricePoints.Seek(market.LowestAsk)
		if iterator != nil {
			defer iterator.Close()
			for {
				pricePoint := iterator.Value().(*PricePoint)
				for _, entry := range pricePoint.SellBookEntries {
					market.SellOrders = append(market.SellOrders, entry)
				}
				if ok := iterator.Next(); !ok {
					break
				}
			}
		}
	}

	if market.HighestBid != 0 {
		iterator := pricePoints.Seek(market.HighestBid)
		if iterator != nil {
			defer iterator.Close()
			for {
				pricePoint := iterator.Value().(*PricePoint)
				for _, entry := range pricePoint.BuyBookEntries {
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
