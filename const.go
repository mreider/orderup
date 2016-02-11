package main

// Some application constants.

const (
	RESTAURANTS = "restaurants" // Restaurants database bucket
	ORDERLIST   = "orders"      // Orders list database bucket
	HISTORY     = "history"     // History database bucket

	V1 = "v1" // Current API version

	// Command set
	CREATE_Q_CMD     = "create-restaurant"
	DELETE_Q_CMD     = "delete-restaurant"
	CREATE_ORDER_CMD = "create-order"
	FINISH_ORDER_CMD = "finish-order"
	LIST_CMD         = "list"
	HISTORY_CMD      = "history"
)
