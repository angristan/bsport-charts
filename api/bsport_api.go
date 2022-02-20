package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Booking struct {
	Name                  string      `json:"name"`
	RecurrenceRuleBooking interface{} `json:"recurrence_rule_booking"`
	IsDiscardable         bool        `json:"is_discardable"`
	OfferDateStart        time.Time   `json:"offer_date_start"`
	OfferDurationMinute   int         `json:"offer_duration_minute"`
	Attendance            bool        `json:"attendance"`
	BookingStatusCode     int         `json:"booking_status_code"`
	Consumer              int         `json:"consumer"`
	ConsumerPaymentPack   int         `json:"consumer_payment_pack"`
	AttendanceDateUpdated interface{} `json:"attendance_date_updated"`
	Date                  time.Time   `json:"date"`
	FirstInCompany        bool        `json:"first_in_company"`
	ID                    int         `json:"id"`
	IsDeleted             bool        `json:"is_deleted"`
	Offer                 int         `json:"offer"`
	Source                int         `json:"source"`
	WasRefunded           bool        `json:"was_refunded"`
	Member                int         `json:"member"`
	CreditConsumed        int         `json:"credit_consumed"`
	Establishment         int         `json:"establishment"`
	Coach                 int         `json:"coach"`
	Level                 int         `json:"level"`
	CoachOverride         interface{} `json:"coach_override"`
	MetaActivity          int         `json:"meta_activity"`
	DateCanceled          interface{} `json:"date_canceled"`
	SpotID                interface{} `json:"spot_id"`
}

type bookingsResponse struct {
	Links struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"links"`
	NextPage int       `json:"next_page"`
	Count    int       `json:"count"`
	Results  []Booking `json:"results"`
}

type bookingsReqParams struct {
	page     int
	pageSize int
	member   string
	token    string
}

func doBookingsRequests(reqParams bookingsReqParams) (bookingsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.production.bsport.io/api/v1/booking/", nil)
	if err != nil {
		return bookingsResponse{}, err
	}

	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", reqParams.page))
	q.Add("page_size", fmt.Sprintf("%d", reqParams.pageSize))
	q.Add("mine", "true")
	q.Add("member", reqParams.member)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("authorization", "Token "+reqParams.token)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return bookingsResponse{}, err
	}
	defer resp.Body.Close()

	var bookingsResp bookingsResponse
	err = json.NewDecoder(resp.Body).Decode(&bookingsResp)
	if err != nil {
		return bookingsResponse{}, err
	}

	return bookingsResp, nil
}

func GetBookings(memberID string, token string) ([]Booking, error) {
	var bookings []Booking
	for page := 1; ; page++ {
		bookingsResp, err := doBookingsRequests(bookingsReqParams{
			page:     page,
			pageSize: 50,
			member:   memberID,
			token:    token,
		})
		if err != nil {
			return bookings, err
		}

		bookings = append(bookings, bookingsResp.Results...)

		if bookingsResp.Links.Next == "" {
			break
		}
	}

	// Skip canceled bookings
	var bookingsFiltered []Booking
	for _, b := range bookings {
		if b.DateCanceled != nil {
			continue
		}
		bookingsFiltered = append(bookingsFiltered, b)
	}

	return bookingsFiltered, nil
}
