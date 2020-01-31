package lunch

import (
	"github.com/lunjon/lunch/internal/pkg/menu"
)

// Lunch is an interface for types that can
// fetch menus from different restaurants.
type Lunch interface {
	GetMenu() (*menu.Menu, error)
}

var WeekDays = []string{
	"monday",
	"tuesday",
	"wednesday",
	"thursday",
	"friday",
}