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

type Orders struct {
	orderItems []*Order
}

func (orders *Orders) NewOrder(columnName string) *Orders {
	orders.orderItems = append(orders.orderItems, NewOrder(columnName))
	return orders
}

func (orders *Orders) NewDescOrder(columnName string) *Orders {
	orders.orderItems = append(orders.orderItems, NewDescOrder(columnName))
	return orders
}

func (orders *Orders) Build() []*Order {
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
