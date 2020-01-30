package meny

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Course defines a meal type with a description.
type Course struct {
	courseType  string
	description string
}

// String returns a string representation of this course.
func (course *Course) String() string {
	return fmt.Sprintf("%s: %s", course.courseType, course.description)
}

// NewCourse creates a new course.
func NewCourse(t, d string) *Course {
	return &Course{
		courseType:  t,
		description: d,
	}
}

// Day represents the part of a meny on a particular da.
type Day struct {
	weekDay string
	courses []*Course
}

// NewDay creates a new day.
func NewDay(wd string, courses ...*Course) *Day {
	return &Day{
		weekDay: wd,
		courses: courses,
	}
}

// Weekday returns the name of this day.
func (day *Day) Weekday() string {
	return day.weekDay
}

// Today returns true if this day is todays date.
func (day *Day) Today() bool {
	weekdayEnglish := strings.ToLower(time.Now().Weekday().String())
	return day.Translate() == weekdayEnglish
}

var translations = map[string][]string{
	"monday":    {"måndag", "monday"},
	"tuesday":   {"tisday", "tuesdag"},
	"wednesday": {"onsdag", "wednesday"},
	"thursday":  {"torsdag", "thursday"},
	"friday":    {"fredag", "friday"},
	"saturday":  {"lördag", "saturday"},
	"sunday":    {"söndag", "sunday"},
}

// Translate this days name to english.
func (day *Day) Translate() string {
	d := strings.ToLower(day.weekDay)

	for eng, trans := range translations {
		for _, name := range trans {
			if d == name {
				return eng
			}
		}
	}

	return day.weekDay
}

// Meny is the complete collection of courses
// over a week.
type Meny struct {
	Days   []*Day
	output io.Writer
	// TODO: add week number
	number int
}

// NewMeny constructs a meny consisting of the given days.
func NewMeny(days ...*Day) *Meny {
	return &Meny{
		Days:   days,
		output: os.Stdout,
	}
}

// Render the meny to its output.
func (meny *Meny) Render() {
	weekday := strings.ToLower(time.Now().Weekday().String())

	type row struct {
		bold bool
		data []string
	}

	var rows []row
	for _, day := range meny.Days {
		for _, c := range day.courses {
			r := row{
				bold: day.Today(),
				data: []string{day.weekDay, c.String()},
			}
			rows = append(rows, r)
		}
	}

	table := tablewriter.NewWriter(meny.output)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetHeader([]string{"Day", "Courses"})
	table.SetColWidth(80)

	for _, row := range rows {
		if row.bold {
			// table.Rich(row.data, []tablewriter.Colors{tablewriter.Colors{tablewriter.Bold}})
			colors := []tablewriter.Colors{
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
			}
			table.Rich(row.data, colors)
		} else {
			table.Append(row.data)
		}
	}

	table.Render()
	fmt.Println(weekday)
}

// Count returns the number of days in this meny.
func (meny *Meny) Count() int {
	return len(meny.Days)
}

// SetOutput changes the default output of os.Stdout
// to o instead.
func (meny *Meny) SetOutput(o io.Writer) {
	meny.output = o
}

// FilterDay runs the predicate for each day in the meny
// and keeps only those that returns true on the predicate
// function.
func (meny *Meny) FilterDay(pred func(*Day) bool) {
	filtered := make([]*Day, 0)
	for _, day := range meny.Days {
		if pred(day) {
			filtered = append(filtered, day)
		}
	}

	meny.Days = filtered
}

// Today is a function that can be used to
// run FilterDay with to filter todays day
// in the meny.
func Today(day *Day) bool {
	today := time.Now().Weekday().String()
	return strings.ToLower(today) == strings.ToLower(day.weekDay)
}
