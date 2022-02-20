package charts

import (
	"github.com/angristan/bsport-charts/api"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	weekdaysShort = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
)

func genDataForBookingsPerWeekDayBarChart(bookings []api.Booking) []opts.BarData {
	data := make(map[int]int)
	for i := 0; i < 7; i++ {
		data[i] = 0
	}

	for _, b := range bookings {
		data[int(b.OfferDateStart.Local().Weekday())]++
	}

	items := make([]opts.BarData, 0)
	for i := 0; i < len(weekdaysShort); i++ {
		items = append(items, opts.BarData{Value: data[i], Name: weekdaysShort[i]})
	}
	return items
}

func BookingsPerWeekDayBarChart(bookings []api.Booking) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Number of bookings per week day"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
	)

	bar.SetXAxis(weekdaysShort).
		AddSeries("bookings", genDataForBookingsPerWeekDayBarChart(bookings)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)
	return bar
}
