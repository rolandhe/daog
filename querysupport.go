// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog



// Pager 分页参数结构，PageSize 每一页的大小，PageNumber 页码，从1算起
type Pager struct {
	PageSize   int
	PageNumber int
}

func NewPager(pageSize int, pageNumber int) *Pager {
	return &Pager{pageSize, pageNumber}
}


// Order 描述sql中的单个 order 条件
type Order struct {
	ColumnName string
	Desc       bool
}


func NewOrder(columnName string) *Order {
	return &Order{columnName, false}
}

func NewDescOrder(columnName string) *Order {
	return &Order{columnName, true}
}

// NewOrdersBuilder 构建 OrdersBuilder对象
func NewOrdersBuilder() *OrdersBuilder {
	return &OrdersBuilder{}
}

// OrdersBuilder 构建order by 条件工具
type OrdersBuilder struct {
	orderItems []*Order
}

// NewOrder 增加一个升序的条件
func (orders *OrdersBuilder) NewOrder(columnName string) *OrdersBuilder {
	orders.orderItems = append(orders.orderItems, NewOrder(columnName))
	return orders
}

// NewDescOrder 增加一个降序的条件
func (orders *OrdersBuilder) NewDescOrder(columnName string) *OrdersBuilder {
	orders.orderItems = append(orders.orderItems, NewDescOrder(columnName))
	return orders
}

// Build 构建出最终的 order by sql  片段
func (orders *OrdersBuilder) Build() []*Order {
	return orders.orderItems
}


