package daog

type Order struct {
	ColumnName string
	Desc       bool
}

// Pager 分页参数结构，PageSize 每一页的大小，PageNumber 页码，从0算起
type Pager struct {
	PageSize   int
	PageNumber int
}

type OrdersBuilder struct {
	orderItems []*Order
}

func NewOrdersBuilder() *OrdersBuilder {
	return &OrdersBuilder{}
}

func (orders *OrdersBuilder) NewOrder(columnName string) *OrdersBuilder {
	orders.orderItems = append(orders.orderItems, NewOrder(columnName))
	return orders
}

func (orders *OrdersBuilder) NewDescOrder(columnName string) *OrdersBuilder {
	orders.orderItems = append(orders.orderItems, NewDescOrder(columnName))
	return orders
}

func (orders *OrdersBuilder) Build() []*Order {
	return orders.orderItems
}

func NewPager(pageSize int, pageNumber int) *Pager {
	return &Pager{pageSize, pageNumber}
}

func NewOrder(columnName string) *Order {
	return &Order{columnName, false}
}

func NewDescOrder(columnName string) *Order {
	return &Order{columnName, true}
}
