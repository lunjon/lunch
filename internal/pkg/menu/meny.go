package menu

import (
	"fmt"
	"io"
	"log"
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

// Day represents the part of a menu on a particular da.
type Day struct {
	weekDay string
	courses []*Course
}

func EmptyDay() *Day {
	return &Day{}
}

func (day *Day) SetWeekday(s string) {
	day.weekDay = s
}

func (day *Day) AddCourse(c *Course) {
	day.courses = append(day.courses, c)
}

// NewDay creates a new day.
func NewDay(wd string, courses ...*Course) *Day {
	return &Day{
		weekDay: strings.TrimSpace(wd),
		courses: courses,
	}
}

// Copy returns the name of this day.
func (day *Day) Copy() *Day {
	courses := make([]*Course, len(day.courses))
	copy(courses, day.courses)

	return &Day{
		weekDay: day.weekDay,
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
	translated := day.Translate()
	return translated == weekdayEnglish
}

var translations = map[string][]string{
	"monday":    {"måndag", "monday"},
	"tuesday":   {"tisdag", "tuesday"},
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

	log.Printf("No translation found, returning weekday: %s", day.weekDay)
	return day.weekDay
}

// GetMenu is the complete collection of courses
// over a week.
type Menu struct {
	name   string
	Days   []*Day
	output io.Writer
	// TODO: add week number
	number int
}

// NewMenu constructs a menu consisting of the given days.
func NewMenu(name string, days ...*Day) *Menu {
	return &Menu{
		Days:   days,
		output: os.Stdout,
	}
}

// Render the menu to its output.
func (menu *Menu) Render() {
	type row struct {
		bold bool
		data []string
	}

	var rows []row
	for _, day := range menu.Days {
		for _, c := range day.courses {
			r := row{
				bold: day.Today(),
				data: []string{day.weekDay, c.String()},
			}
			rows = append(rows, r)
		}
	}

	log.Printf("Found %d row(s)", len(rows))

	table := tablewriter.NewWriter(menu.output)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetHeader([]string{"Day", "Courses"})
	table.SetColWidth(80)

	for _, row := range rows {
		if row.bold {
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
}

// Count returns the number of days in this menu.
func (menu *Menu) Count() int {
	return len(menu.Days)
}

// SetOutput changes the default output of os.Stdout
// to o instead.
func (menu *Menu) SetOutput(o io.Writer) {
	menu.output = o
}

// FilterDay runs the predicate for each day in the menu
// and keeps only those that returns true on the predicate
// function.
func (menu *Menu) FilterDay(pred func(*Day) bool) {
	filtered := make([]*Day, 0)
	for _, day := range menu.Days {
		if pred(day) {
			filtered = append(filtered, day)
		}
	}

	menu.Days = filtered
}

func (menu *Menu) Name() string {
	return menu.name
}

// Today is a function that can be used to
// run FilterDay with to filter today's day
// in the menu.
func Today(day *Day) bool {
	today := time.Now().Weekday().String()
	return strings.ToLower(today) == day.Translate()

}
