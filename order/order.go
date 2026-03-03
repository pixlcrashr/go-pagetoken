package order

import "fmt"

type Order bool

func (o Order) String() string {
	if o == Desc {
		return "desc"
	}

	return "asc"
}

const (
	Asc  Order = false
	Desc Order = true
)

func (o *Order) UnmarshalString(s string) error {
	switch s {
	case "asc":
		*o = Asc
	case "desc":
		*o = Desc
	default:
		return fmt.Errorf("invalid order: %s", s)
	}

	return nil
}
