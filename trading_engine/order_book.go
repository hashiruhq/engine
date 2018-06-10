package trading_engine

// Process a new received order and return a list of trades make
func (orderBook *OrderBook) Process(order *Order) ([]*Trade, error) {
	if order.Side == BUY {
		return orderBook.processLimitBuy(order)
	}
	return orderBook.processLimitSell(order)
}

func (orderBook *OrderBook) processLimitBuy(order *Order) ([]*Trade, error) {
	bookEntry := NewBookEntry(order)
	trades := make([]*Trade, 0, 5)

	// if there are no more ordes just add the buy order to the list
	n := len(orderBook.SellOrders)
	if n == 0 {
		orderBook.addBookEntry(bookEntry, 0)
		return trades, nil
	}
	i := 0
	// traverse orders to find a matching one based on the sell order list
	for {
		if i >= n {
			orderBook.addBookEntry(bookEntry, i)
			break
		}
		sellOrder := orderBook.SellOrders[i]
		// when we get to an order which is higher than our buy price we stop
		// we add the remaining order to the order book and return the trades made
		if sellOrder.Price > bookEntry.Price {
			orderBook.addBookEntry(bookEntry, i)
			break
		}
		// if we can fill the trade instantly then we add the trade and complete the order
		if sellOrder.Amount >= bookEntry.Amount {
			trades = append(trades, NewTrade(bookEntry.Order, sellOrder.Order, bookEntry.Amount, bookEntry.Price))
			sellOrder.Amount -= bookEntry.Amount
			if sellOrder.Amount == 0 {
				orderBook.removeBookEntry(SELL, i)
			}
			break
		}

		// if the sell order has a lower amount that what the buy order is then we fill only what we can from the sell order,
		// we complete the sell order and we move to the next order
		if sellOrder.Amount < bookEntry.Amount {
			trades = append(trades, NewTrade(bookEntry.Order, sellOrder.Order, sellOrder.Amount, sellOrder.Price))
			bookEntry.Amount -= sellOrder.Amount
			orderBook.removeBookEntry(SELL, i)
			n--
			continue
		}
		i++
	}

	return trades, nil
}

func (orderBook *OrderBook) processLimitSell(order *Order) ([]*Trade, error) {
	bookEntry := NewBookEntry(order)
	trades := make([]*Trade, 0, 5)

	// if there are no more ordes just add the sell order to the list
	n := len(orderBook.BuyOrders)
	if n == 0 {
		orderBook.addBookEntry(bookEntry, 0)
		return trades, nil
	}
	i := 0
	// traverse orders to find a matching one based on the buy order list
	for {
		if i >= n {
			orderBook.addBookEntry(bookEntry, i)
			break
		}
		buyEntry := orderBook.BuyOrders[i]
		// when we get to an order which is higher than our buy price we stop
		// we add the remaining order to the order book and return the trades made
		if buyEntry.Price < bookEntry.Price {
			orderBook.addBookEntry(bookEntry, i)
			break
		}
		// if we can fill the trade instantly then we add the trade and complete the order
		if buyEntry.Amount >= bookEntry.Amount {
			trades = append(trades, NewTrade(bookEntry.Order, buyEntry.Order, bookEntry.Amount, bookEntry.Price))
			buyEntry.Amount -= bookEntry.Amount
			if buyEntry.Amount == 0 {
				orderBook.removeBookEntry(BUY, i)
			}
			break
		}

		// if the buy order has a lower amount that what the sell order is then we fill only what we can from the buy order,
		// we complete the buy order and we move to the next order
		if buyEntry.Amount < bookEntry.Amount {
			trades = append(trades, NewTrade(bookEntry.Order, buyEntry.Order, buyEntry.Amount, buyEntry.Price))
			bookEntry.Amount -= buyEntry.Amount
			orderBook.removeBookEntry(BUY, i)
			n--
			continue
		}
		i++
	}

	return trades, nil
}

// @todo ensure that the orders are sorted by price when they are added to the order book
func (orderBook *OrderBook) addBookEntry(bookEntry *BookEntry, index int) error {
	if bookEntry.Order.Side == BUY {
		orders := append(orderBook.BuyOrders[:index], bookEntry)
		orderBook.BuyOrders = append(orders, orderBook.BuyOrders[index:]...)
	} else {
		orders := append(orderBook.SellOrders[:index], bookEntry)
		orderBook.SellOrders = append(orders, orderBook.SellOrders[index:]...)
	}
	return nil
}

func (orderBook *OrderBook) removeBookEntry(side OrderSide, index int) error {
	if side == BUY {
		orderBook.BuyOrders = append(orderBook.BuyOrders[:index], orderBook.BuyOrders[index+1:]...)
	} else {
		orderBook.SellOrders = append(orderBook.SellOrders[:index], orderBook.SellOrders[index+1:]...)
	}
	return nil
}
