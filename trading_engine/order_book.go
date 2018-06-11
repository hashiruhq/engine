package trading_engine

// Process a new received order and return a list of trades make
func (orderBook *OrderBook) Process(order *Order) ([]*Trade, error) {
	var trades []*Trade
	var err error
	if order.Side == BUY {
		orderBook.mtx.Lock()
		trades, err = orderBook.processLimitBuy(order)
		orderBook.mtx.Unlock()

	} else {
		orderBook.mtx.Lock()
		trades, err = orderBook.processLimitSell(order)
		orderBook.mtx.Unlock()
	}
	return trades, err
}

func (orderBook *OrderBook) processLimitBuy(order *Order) ([]*Trade, error) {
	trades := make([]*Trade, 0, 5)

	// if there are no more ordes just add the buy order to the list
	n := len(orderBook.SellOrders)
	if n == 0 {
		orderBook.addBookEntry(order, 0)
		return trades, nil
	}
	i := 0
	// traverse orders to find a matching one based on the sell order list
	for {
		if i >= n {
			orderBook.addBookEntry(order, i)
			break
		}
		sellOrder := orderBook.SellOrders[i]
		// when we get to an order which is higher than our buy price we stop
		// we add the remaining order to the order book and return the trades made
		if sellOrder.Price > order.Price {
			orderBook.addBookEntry(order, i)
			break
		}
		// if we can fill the trade instantly then we add the trade and complete the order
		if sellOrder.Amount >= order.Amount {
			trades = append(trades, NewTrade(order, sellOrder, order.Amount, order.Price))
			sellOrder.Amount -= order.Amount
			if sellOrder.Amount == 0 {
				orderBook.removeBookEntry(SELL, i)
			}
			break
		}

		// if the sell order has a lower amount that what the buy order is then we fill only what we can from the sell order,
		// we complete the sell order and we move to the next order
		if sellOrder.Amount < order.Amount {
			trades = append(trades, NewTrade(order, sellOrder, sellOrder.Amount, sellOrder.Price))
			order.Amount -= sellOrder.Amount
			orderBook.removeBookEntry(SELL, i)
			n--
			continue
		}
		i++
	}

	return trades, nil
}

func (orderBook *OrderBook) processLimitSell(order *Order) ([]*Trade, error) {
	// bookEntry := NewBookEntry(order)
	trades := make([]*Trade, 0, 5)

	// if there are no more ordes just add the sell order to the list
	n := len(orderBook.BuyOrders)
	if n == 0 {
		orderBook.addBookEntry(order, 0)
		return trades, nil
	}
	i := 0
	// traverse orders to find a matching one based on the buy order list
	for {
		if i >= n {
			orderBook.addBookEntry(order, i)
			break
		}
		buyEntry := orderBook.BuyOrders[i]
		// when we get to an order which is higher than our buy price we stop
		// we add the remaining order to the order book and return the trades made
		if buyEntry.Price < order.Price {
			orderBook.addBookEntry(order, i)
			break
		}
		// if we can fill the trade instantly then we add the trade and complete the order
		if buyEntry.Amount >= order.Amount {
			trades = append(trades, NewTrade(order, buyEntry, order.Amount, order.Price))
			buyEntry.Amount -= order.Amount
			if buyEntry.Amount == 0 {
				orderBook.removeBookEntry(BUY, i)
			}
			break
		}

		// if the buy order has a lower amount that what the sell order is then we fill only what we can from the buy order,
		// we complete the buy order and we move to the next order
		if buyEntry.Amount < order.Amount {
			trades = append(trades, NewTrade(order, buyEntry, buyEntry.Amount, buyEntry.Price))
			order.Amount -= buyEntry.Amount
			orderBook.removeBookEntry(BUY, i)
			n--
			continue
		}
		i++
	}

	return trades, nil
}

// @todo ensure that the orders are sorted by price when they are added to the order book
func (orderBook *OrderBook) addBookEntry(order *Order, index int) error {
	if order.Side == BUY {
		orders := append(orderBook.BuyOrders[:index], order)
		orderBook.BuyOrders = append(orders, orderBook.BuyOrders[index:]...)
	} else {
		orders := append(orderBook.SellOrders[:index], order)
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
