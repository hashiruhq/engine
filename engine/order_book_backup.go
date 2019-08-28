package engine

import (
	"log"

	"gitlab.com/around25/products/matching-engine/model"
)

// Load the full order book from the backup object
func (book *orderBook) Load(market model.MarketBackup) error {
	book.MarketID = market.MarketID
	book.PricePrecision = int(market.PricePrecision)
	book.VolumePrecision = int(market.VolumePrecision)
	book.LowestAsk = market.LowestAsk
	book.HighestBid = market.HighestBid
	book.LowestEntryPrice = market.LowestEntryPrice
	book.HighestLossPrice = market.HighestLossPrice

	// load limit orders
	for _, buyBookEntry := range market.BuyOrders {
		book.addBuyBookEntry(*buyBookEntry)
	}
	for _, sellBookEntry := range market.SellOrders {
		book.addSellBookEntry(*sellBookEntry)
	}

	// load market orders
	book.BuyMarketEntries = make([]model.Order, len(market.BuyMarketEntries))
	for i, order := range market.BuyMarketEntries {
		book.BuyMarketEntries[i] = *order
	}
	book.SellMarketEntries = make([]model.Order, len(market.SellMarketEntries))
	for i, order := range market.SellMarketEntries {
		book.SellMarketEntries[i] = *order
	}

	// load stop orders
	for _, order := range market.StopEntryOrders {
		book.StopEntryOrders.addOrder(order.StopPrice, *order)
	}
	for _, order := range market.StopLossOrders {
		book.StopLossOrders.addOrder(order.StopPrice, *order)
	}

	return nil
}

// Backup the order book in another structure for exporting
func (book *orderBook) Backup() model.MarketBackup {
	market := model.MarketBackup{
		MarketID:          book.MarketID,
		PricePrecision:    int32(book.PricePrecision),
		VolumePrecision:   int32(book.VolumePrecision),
		LowestAsk:         book.LowestAsk,
		HighestBid:        book.HighestBid,
		LowestEntryPrice:  book.LowestEntryPrice,
		HighestLossPrice:  book.HighestLossPrice,
		BuyOrders:         make([]*model.Order, 0, book.BuyEntries.Len()),
		SellOrders:        make([]*model.Order, 0, book.SellEntries.Len()),
		BuyMarketEntries:  make([]*model.Order, len(book.BuyMarketEntries)),
		SellMarketEntries: make([]*model.Order, len(book.SellMarketEntries)),
		StopEntryOrders:   make([]*model.Order, 0, 0),
		StopLossOrders:    make([]*model.Order, 0, 0),
	}
	log.Println("Sell market len:", book.MarketID, len(book.SellMarketEntries))

	// backup limit orders
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

	// backup market orders
	for i, order := range book.BuyMarketEntries {
		market.BuyMarketEntries[i] = &order
	}

	log.Println("Sell market len:", book.MarketID, len(book.SellMarketEntries))
	for i, order := range book.SellMarketEntries {
		log.Printf("Sell market item: %d %#v\n", i, order)
		market.SellMarketEntries[i] = &order
	}

	// backup stop orders
	if market.LowestEntryPrice != 0 {
		iterator := book.StopEntryOrders.Seek(market.LowestEntryPrice)
		if iterator != nil {
			for {
				pricePoint := iterator.Value()
				for _, entry := range pricePoint.Entries {
					market.StopEntryOrders = append(market.StopEntryOrders, &entry)
				}
				if ok := iterator.Next(); !ok {
					break
				}
			}
			iterator.Close()
		}
	}

	if market.HighestLossPrice != 0 {
		iterator := book.BuyEntries.Seek(market.HighestLossPrice)
		if iterator != nil {
			for {
				pricePoint := iterator.Value()
				for _, entry := range pricePoint.Entries {
					market.StopLossOrders = append(market.StopLossOrders, &entry)
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
