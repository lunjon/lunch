package edison

import (
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/lunjon/lunch/internal/pkg/meny"
)

var (
	itemChannel chan *meny.Day
)

// Collect ...
func Collect() (*meny.Meny, error) {
	days := []string{"monday", "tuesday", "wednesday", "thursday", "friday"}

	itemChannel = make(chan *meny.Day, len(days))
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

	items := make([]*meny.Day, 0)
	for d := range itemChannel {
		if d == nil {
			fmt.Println("Failed to parse Edison menu...")
			os.Exit(1)
		}
		items = append(items, d)
	}

	return meny.NewMeny(items...), nil
}

func parseDay(element *colly.HTMLElement) {
	day := element.ChildText("h3")
	courses := make([]*meny.Course, 0)

	element.ForEach("tr", func(i int, element *colly.HTMLElement) {
		courseType := element.ChildText("td[class=course_type]")
		courseDescription := element.ChildText("td[class=course_description]")

		c := meny.NewCourse(courseType, courseDescription)
		courses = append(courses, c)
	})

	m := meny.NewDay(day, courses...)
	itemChannel <- m
}
