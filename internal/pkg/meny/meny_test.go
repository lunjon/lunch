package meny_test

import (
	"bytes"
	"github.com/lunjon/lunch/internal/pkg/meny"
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
	m.FilterDay(func(d *meny.Day) bool {
		if d.Weekday() == "Monday" {
			return true
		}
		return false
	})

	if m.Count() != 1 {
		t.Fail()
	}
}

var monday = meny.NewDay(
	"Monday",
	meny.NewCourse("Veg", "Pancakes"),
	meny.NewCourse("Local", "Köttbullar"),
)

var tuesday = meny.NewDay(
	"Tuesday",
	meny.NewCourse("Veg", "Soup"),
	meny.NewCourse("Local", "Pizza"),
)

var testMeny = meny.NewMeny(
	monday,
	tuesday,
)
