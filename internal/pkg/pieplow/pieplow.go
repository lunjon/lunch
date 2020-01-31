package pieplow

import (
	"github.com/gocolly/colly/v2"
	"github.com/lunjon/lunch/internal/pkg/lunch"
	"github.com/lunjon/lunch/internal/pkg/menu"
	"log"
	"strconv"
	"strings"
)

var (
	itemChannel chan *menu.Day
)

type Pieplow struct {
	url string
}

func New() *Pieplow {
	return &Pieplow{
		url: "http://lund.pieplowsmat.se/street-food/",
	}
}

// GetMenu ...
func (pieplow *Pieplow) GetMenu() (*menu.Menu, error) {
	itemChannel = make(chan *menu.Day, len(lunch.WeekDays))
	collector := colly.NewCollector()

	weekNumber := 0

	// Find the weekday
	collector.OnHTML("h3", func(element *colly.HTMLElement) {
		text := strings.ToLower(element.Text)
		if !strings.Contains(text, "week") {
			return
		}

		split := strings.Split(text, " ")
		if len(split) != 2 {
			return
		}

		d, err := strconv.Atoi(split[1])
		if err != nil {
			return
		}

		weekNumber = d
	})

	// The outer wrapper
	collector.OnHTML("div[class=wpb_wrapper]", func(element *colly.HTMLElement) {
		element.ForEachWithBreak("div[class=wpb_wrapper]", innerParser)
	})
	collector.OnScraped(func(_ *colly.Response) {
		close(itemChannel)
	})

	err := collector.Visit(pieplow.url)
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

	log.Printf("Found week number %d", weekNumber)

	return menu.NewMenu("Pieplow", items...), nil
}

func innerParser(i int, outerElement *colly.HTMLElement) bool {
	// The level at which the lunch menu appear
	if i != 3 {
		return true
	}

	test := strings.ToLower(outerElement.ChildText("strong"))
	if !strings.Contains(test, "monday") {
		return false
	}

	index := 0
	day := menu.EmptyDay()
	dayCount := 0

	outerElement.ForEachWithBreak("p", func(i int, element *colly.HTMLElement) bool {
		if dayCount == 5 {
			return false
		}

		switch index {
		case 0:
			wd := element.ChildText("strong")
			day.SetWeekday(wd)
		case 1:
			day.AddCourse(menu.NewCourse("Meat", element.Text))
		case 2:
			day.AddCourse(menu.NewCourse("Veg", element.Text))
		}

		if index == 2 {
			// Reset
			index = 0
			itemChannel <- day.Copy()
			day = menu.EmptyDay()
			dayCount++
		} else {
			index++

		}

		return true
	})

	return true
}
