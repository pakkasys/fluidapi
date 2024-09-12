package order

type IOrderable interface {
	GetOrders() []Order
}
