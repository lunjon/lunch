package edison

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/lunjon/lunch/internal/pkg/lunch"
	"github.com/lunjon/lunch/internal/pkg/menu"
	"log"
)

var (
	itemChannel chan *menu.Day
)

type Edison struct {
	url string
}

func New() *Edison {
	return &Edison{
		url: "http://restaurangedison.se/lunch",
	}
}

// GetMenu ...
func (edison *Edison) GetMenu() (*menu.Menu, error) {
	itemChannel = make(chan *menu.Day, len(lunch.WeekDays))
	collector := colly.NewCollector()

	for _, id := range lunch.WeekDays {
		query := fmt.Sprintf("div[id=%s]", id)
		collector.OnHTML(query, parseDay)
	}

	collector.OnScraped(func(_ *colly.Response) {
		log.Print("Scraping done, closing channel")
		close(itemChannel)
	})

	err := collector.Visit(edison.url)
	if err != nil {
		return nil, err
	}

	items := make([]*menu.Day, 0)
	for d := range itemChannel {
		if d == nil {
			return nil, err
		}
		items = append(items, d)
	}

	return menu.NewMenu("Edison", items...), nil
}

func parseDay(element *colly.HTMLElement) {
	day := element.ChildText("h3")
	courses := make([]*menu.Course, 0)

	element.ForEach("tr", func(i int, element *colly.HTMLElement) {
		courseType := element.ChildText("td[class=course_type]")
		courseDescription := element.ChildText("td[class=course_description]")

		c := menu.NewCourse(courseType, courseDescription)
		courses = append(courses, c)
	})

	m := menu.NewDay(day, courses...)

	log.Print("Sending Day on channel")
	itemChannel <- m
}
