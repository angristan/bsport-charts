package charts

import (
	"time"

	"github.com/angristan/bsport-charts/api"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	weekDays = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	dayHrs   = [...]string{
		"08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
	}
)

type Slot struct {
	weekday time.Weekday
	hour    int
}

func genHeatMapData(bookings []api.Booking) []opts.HeatMapData {
	// Init slots
	slots := make(map[Slot]int)
	for i := 0; i < 7; i++ {
		for j := 8; j < 22; j++ {
			slots[Slot{time.Weekday(i), j}] = 0
		}
	}

	for _, b := range bookings {
		for hour := b.OfferDateStart.Local().Hour(); hour < b.OfferDateStart.Add(time.Duration(b.OfferDurationMinute*int(time.Minute))).Hour(); hour++ {
			slots[Slot{b.OfferDateStart.Local().Weekday(), hour}]++
		}
	}

	items := make([]opts.HeatMapData, 0)

	for i := 0; i < 7; i++ {
		for j := 8; j < 22; j++ {
			items = append(items, opts.HeatMapData{Value: [3]interface{}{j, time.Weekday(i).String(), slots[Slot{time.Weekday(i), j}]}})
		}
	}

	return items
}

func genDataForBookingsHeatMapChart(bookings []api.Booking) int {
	// Init slots
	slots := make(map[Slot]int)
	for i := 0; i < 7; i++ {
		for j := 8; j < 22; j++ {
			slots[Slot{time.Weekday(i), j}] = 0
		}
	}

	max := 0
	for _, b := range bookings {
		for hour := b.OfferDateStart.Local().Hour(); hour < b.OfferDateStart.Add(time.Duration(b.OfferDurationMinute*int(time.Minute))).Hour(); hour++ {
			slots[Slot{b.OfferDateStart.Local().Weekday(), hour}]++
			if slots[Slot{b.OfferDateStart.Local().Weekday(), hour}] > max {
				max = slots[Slot{b.OfferDateStart.Local().Weekday(), hour}]
			}
		}
	}

	return max
}

func BookingsHeatMapChart(bookings []api.Booking) *charts.HeatMap {
	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Heatmap of bookings",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      weekDays,
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        float32(genDataForBookingsHeatMapChart(bookings)),
			InRange: &opts.VisualMapInRange{
				Color: []string{"#50a3ba", "#eac736", "#d94e5d"},
			},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
	)

	hm.SetXAxis(dayHrs).AddSeries("heatmap", genHeatMapData(bookings))
	return hm
}
