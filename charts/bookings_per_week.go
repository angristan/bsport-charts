package charts

import (
	"fmt"
	"time"

	"github.com/angristan/bsport-charts/api"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func genDataForBookingsPerWeekLineChart(bookings []api.Booking, weeks []string) []opts.LineData {
	data := make(map[string]int)
	for _, b := range bookings {
		year, week := b.OfferDateStart.Local().ISOWeek()
		data[fmt.Sprintf("%d-%d", year, week)]++
	}

	items := make([]opts.LineData, 0)

	for _, week := range weeks {
		if val, ok := data[week]; ok {
			items = append(items, opts.LineData{Value: val, Name: week})
		} else {
			items = append(items, opts.LineData{Value: 0, Name: week})
		}
	}

	return items
}

func BookingsPerWeekLineChart(bookings []api.Booking) *charts.Line {
	oldestBooking := bookings[len(bookings)-1]

	weeks := []string{}
	date := oldestBooking.OfferDateStart.Local()
	for date.Before(time.Now()) {
		year, week := date.ISOWeek()
		weeks = append(weeks, fmt.Sprintf("%d-%d", year, week))
		date = date.AddDate(0, 0, 7)
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Number of bookings per week (n=%d)", len(bookings)),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Number of bookings",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Week number",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   100,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		// charts.WithInitializationOpts(opts.Initialization{Theme: "shine"}),
	)

	line.SetXAxis(weeks).
		AddSeries("bookings", genDataForBookingsPerWeekLineChart(bookings, weeks)).
		SetSeriesOptions(
			charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
				Name: "Avg",
				Type: "average",
			}),
			charts.WithLabelOpts(opts.Label{
				Show: true,
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.2,
			}),
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
		)
	return line
}
