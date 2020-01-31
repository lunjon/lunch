package menu_test

import (
	"bytes"
	"github.com/lunjon/lunch/internal/pkg/menu"
	"strings"
	"testing"
)

func TestTableOutputFormat(t *testing.T) {
	m := testMeny
	output := bytes.NewBufferString("")
	m.SetOutput(output)
	m.Render()

	preds := []func(s string) bool{
		func(s string) bool { return s != "" },
		func(s string) bool { return strings.Contains(s, "Monday") },
		func(s string) bool { return strings.Contains(s, "Veg") },
		func(s string) bool { return strings.Contains(s, "Tuesday") },
		func(s string) bool { return !strings.Contains(s, "LOL") },
	}

	rendered := output.String()
	for i, pred := range preds {
		if !pred(rendered) {
			t.Fatalf("Predicate number %d failed", i+1)
		}
	}
}

func TestFilterDayMonday(t *testing.T) {
	m := testMeny
	m.FilterDay(func(d *menu.Day) bool {
		if d.Weekday() == "Monday" {
			return true
		}
		return false
	})

	if m.Count() != 1 {
		t.Fail()
	}
}

var monday = menu.NewDay(
	"Monday",
	menu.NewCourse("Veg", "Pancakes"),
	menu.NewCourse("Local", "Köttbullar"),
)

var tuesday = menu.NewDay(
	"Tuesday",
	menu.NewCourse("Veg", "Soup"),
	menu.NewCourse("Local", "Pizza"),
)

var testMeny = menu.NewMenu(
	"fake",
	monday,
	tuesday,
)
