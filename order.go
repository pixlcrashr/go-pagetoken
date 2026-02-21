package pagetoken

import "fmt"

type Order bool

func (o Order) String() string {
	if o == OrderAsc {
		return "1"
	}

	return "0"
}

const (
	OrderAsc  Order = true
	OrderDesc Order = false
)

func ParseOrder(s string) (Order, error) {
	switch s {
	case "1":
		return OrderAsc, nil
	case "0":
		return OrderDesc, nil
	default:
		return OrderDesc, fmt.Errorf("invalid order: %s", s)
	}
}
