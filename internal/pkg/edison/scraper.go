package edison

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/olekukonko/tablewriter"
	"os"
)

type course struct {
	Type        string
	Description string
}

type menyItem struct {
	day     string
	courses []*course
}

type Meny struct {
	items []*menyItem
}

func (meny *Meny) Render() {
	var data [][]string
	for _, item := range meny.items {
		for _, c := range item.courses {
			data = append(data, []string{item.day, c.Description})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetHeader([]string{"Day", "Courses"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

var (
	itemChannel chan *menyItem
)

// Collect ...
func Collect() (*Meny, error) {
	days := []string{"monday", "tuesday", "wednesday", "thursday", "friday"}

	itemChannel = make(chan *menyItem, len(days))
	collector := colly.NewCollector()

	for _, id := range days {
		query := fmt.Sprintf("div[id=%s]", id)
		collector.OnHTML(query, parseDay)
	}

	collector.OnScraped(func(_ *colly.Response) {
		close(itemChannel)
	})

	err := collector.Visit("http://restaurangedison.se/lunch")
	if err != nil {
		return nil, err
	}

	items := make([]*menyItem, 0)
	for d := range itemChannel {
		if d == nil {
			fmt.Println("Failed to parse Edison menu...")
			os.Exit(1)
		}
		items = append(items, d)
	}

	return &Meny{items: items}, nil
}

func parseDay(element *colly.HTMLElement) {
	day := element.ChildText("h3")
	courses := make([]*course, 0)

	element.ForEach("tr", func(i int, element *colly.HTMLElement) {
		courseType := element.ChildText("td[class=course_type]")
		courseDescription := element.ChildText("td[class=course_description]")

		c := &course{
			Type:        courseType,
			Description: courseDescription,
		}
		courses = append(courses, c)
	})

	m := &menyItem{
		day:     day,
		courses: courses,
	}

	itemChannel <- m
}
