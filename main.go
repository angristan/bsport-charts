package main

import (
	"flag"
	"io"
	"os"

	"github.com/angristan/bsport-charts/api"
	"github.com/angristan/bsport-charts/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/sirupsen/logrus"
)

var (
	token    string
	memberID string
)

func main() {
	flag.StringVar(&token, "token", "", "BSport API token")
	flag.StringVar(&memberID, "member", "", "BSport member ID")
	flag.Parse()

	if token == "" {
		logrus.Fatal("missing -token flag")
	}

	if memberID == "" {
		logrus.Fatal("missing -member flag")
	}

	bookings, err := api.GetBookings(memberID, token)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get bookings from API")
	}

	page := components.NewPage()
	page.AddCharts(
		charts.BookingsHeatMapChart(bookings),
		charts.BookingsPerWeekDayBarChart(bookings),
		charts.BookingsPerWeekLineChart(bookings),
	)

	page.PageTitle = "BSport Charts"

	f, err := os.Create("charts.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	logrus.Info("Generated charts.html")
}
