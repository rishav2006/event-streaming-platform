package controllers

import "sync"

type Demo struct {
	mu                  sync.Mutex
	LastFileNumOrder    int
	LastFileNumPayment  int
	Checker             bool
	OffsetSliceOrders   []int
	OffsetSlicePayments []int
	ExsOffset           int
	CounterOrder        int
	CounterPayment      int
}

type AnswerNoGroup struct {
	Message   string
	Partition string
	Topic     string
}

type Answer struct {
	Group     string
	Consumer  string
	Message   string
	Partition string
}

// Check for checker, if it's false then the offset has not been calculated yet, so calculate it

// path for orders requests
var orders = []string{
	"internals/files/folders/orders/o0.log",
	"internals/files/folders/orders/o1.log",
	"internals/files/folders/orders/o2.log",
}

// path for payments requests
var payments = []string{
	"internals/files/folders/payments/p0.log",
	"internals/files/folders/payments/p1.log",
}

// path for default requests
var defaultPath string = "internals/files/folders/default/default.log"
