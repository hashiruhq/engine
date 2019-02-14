package engine

func (book *orderBook) GetMarketDepth() *MarketDepth {
	depth := MarketDepth{
		Market: book.MarketID,
		Bid:    make([]*MarketDepthPriceLevel, 0, 100),
		Ask:    make([]*MarketDepthPriceLevel, 0, 100),
		// @todo set last price and volume with the last generated trade
		LastPrice:  0,
		LastVolume: 0,
	}

	// load ask price levels
	if book.LowestAsk != 0 {
		iterator := book.SellEntries.Seek(book.LowestAsk)
		if iterator != nil {
			for i := 0; i < 100; i++ {
				pricePoint := iterator.Value()
				volume := uint64(0)
				for _, entry := range pricePoint.Entries {
					volume += entry.Amount
				}
				depth.Ask = append(depth.Ask, &MarketDepthPriceLevel{Price: iterator.Key(), Volume: volume})
				if ok := iterator.Next(); !ok {
					break
				}
			}
			iterator.Close()
		}
	}

	// load bid price levels
	if book.HighestBid != 0 {
		iterator := book.BuyEntries.Seek(book.HighestBid)
		if iterator != nil {
			for i := 0; i < 100; i++ {
				pricePoint := iterator.Value()
				volume := uint64(0)
				for _, entry := range pricePoint.Entries {
					volume += entry.Amount
				}
				depth.Bid = append(depth.Bid, &MarketDepthPriceLevel{Price: iterator.Key(), Volume: volume})
				if ok := iterator.Previous(); !ok {
					break
				}
			}
			iterator.Close()
		}
	}
	return &depth
}
