package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/database"
)

type DailyClick struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type WeeklyClick struct {
	WeekStart string `json:"week_start"`
	Clicks    int64  `json:"clicks"`
}

type MonthlyClick struct {
	Month  string `json:"month"`
	Clicks int64  `json:"clicks"`
}

type CountryClick struct {
	Country string `json:"country"`
	Clicks  int64  `json:"clicks"`
}

type CityClick struct {
	City   string `json:"city"`
	Clicks int64  `json:"clicks"`
}

type DeviceClick struct {
	Device string `json:"device"`
	Clicks int64  `json:"clicks"`
}

type BrowserClick struct {
	Browser string `json:"browser"`
	Clicks  int64  `json:"clicks"`
}

type ReferrerClick struct {
	Referrer string `json:"referrer"`
	Clicks   int64  `json:"clicks"`
}

type UrlAnalytics struct {
	ShortCode      string          `json:"short_code"`
	TotalClicks    int64           `json:"total_clicks"`
	UniqueVisitors int64           `json:"unique_visitors"`
	DailyClicks    []DailyClick    `json:"daily_clicks"`
	WeeklyClicks   []WeeklyClick   `json:"weekly_clicks"`
	MonthlyClicks  []MonthlyClick  `json:"monthly_clicks"`
	TopCountries   []CountryClick  `json:"top_countries"`
	TopCities      []CityClick     `json:"top_cities"`
	TopDevices     []DeviceClick   `json:"top_devices"`
	TopBrowsers    []BrowserClick  `json:"top_browsers"`
	TopReferrers   []ReferrerClick `json:"top_referrers"`
}

func GetURLAnalytics(r *http.Request) (UrlAnalytics, error, int) {
	ctx := context.Background()
	shortCode := r.PathValue("urlCode")

	if shortCode == "" {
		return UrlAnalytics{}, fmt.Errorf("url code is required"), http.StatusBadRequest
	}

	authHeader := r.Header.Get("Authorization")
	parts := strings.Fields(authHeader)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return UrlAnalytics{}, fmt.Errorf("invalid authorization header"), http.StatusUnauthorized
	}

	token := parts[1]
	ok, claims, err := auth.ValidateJWTHelper(token, Cfg.JWTSecret)
	if !ok {
		return UrlAnalytics{}, err, http.StatusUnauthorized
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("something went wrong"), http.StatusInternalServerError
	}

	urlOwner, err := Cfg.Db.GetUrlOwnerByShortCode(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return UrlAnalytics{}, fmt.Errorf("URL not found"), http.StatusNotFound
		}
		return UrlAnalytics{}, fmt.Errorf("failed to verify URL ownership"), http.StatusInternalServerError
	}

	if urlOwner != userID {
		return UrlAnalytics{}, fmt.Errorf("forbidden"), http.StatusForbidden
	}

	analytics := UrlAnalytics{
		ShortCode: shortCode,
	}

	totalClicksResult, err := Cfg.Db.GetURLTotalClicks(ctx, shortCode)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("failed to get total clicks"), http.StatusInternalServerError
	}
	analytics.TotalClicks = totalClicksResult.(int64)

	uniqueVisitorsResult, err := Cfg.Db.GetURLUniqueVisitors(ctx, shortCode)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("failed to get unique visitors"), http.StatusInternalServerError
	}
	analytics.UniqueVisitors = uniqueVisitorsResult.(int64)

	dailyClicks, err := Cfg.Db.GetURLDailyClicks(ctx, shortCode)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("failed to get daily clicks"), http.StatusInternalServerError
	}
	for _, dc := range dailyClicks {
		analytics.DailyClicks = append(analytics.DailyClicks, DailyClick{
			Date:   dc.Date.Format("2006-01-02"),
			Clicks: int64(dc.TotalClicks),
		})
	}

	analytics.WeeklyClicks = aggregateToWeekly(dailyClicks)
	analytics.MonthlyClicks = aggregateToMonthly(dailyClicks)

	geoClicks, err := Cfg.Db.GetURLGeoClicks(ctx, shortCode)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("failed to get geo clicks"), http.StatusInternalServerError
	}
	for _, gc := range geoClicks {
		analytics.TopCountries = append(analytics.TopCountries, CountryClick{
			Country: gc.Country,
			Clicks:  int64(gc.TotalClicks),
		})
		analytics.TopCities = append(analytics.TopCities, CityClick{
			City:   gc.City,
			Clicks: int64(gc.TotalClicks),
		})
	}

	deviceClicks, err := Cfg.Db.GetURLDeviceClicks(ctx, shortCode)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("failed to get device clicks"), http.StatusInternalServerError
	}
	for _, device := range deviceClicks {
		analytics.TopDevices = append(analytics.TopDevices, DeviceClick{
			Device: device.Device,
			Clicks: int64(device.TotalClicks),
		})
	}

	browserClicks, err := Cfg.Db.GetURLBrowserClicks(ctx, shortCode)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("failed to get browser clicks"), http.StatusInternalServerError
	}
	for _, browser := range browserClicks {
		analytics.TopBrowsers = append(analytics.TopBrowsers, BrowserClick{
			Browser: browser.Browser,
			Clicks:  int64(browser.TotalClicks),
		})
	}

	referrerClicks, err := Cfg.Db.GetURLReferrerClicks(ctx, shortCode)
	if err != nil {
		return UrlAnalytics{}, fmt.Errorf("failed to get referrer clicks"), http.StatusInternalServerError
	}
	for _, ref := range referrerClicks {
		analytics.TopReferrers = append(analytics.TopReferrers, ReferrerClick{
			Referrer: ref.Referrer,
			Clicks:   int64(ref.TotalClicks),
		})
	}

	return analytics, nil, http.StatusOK
}

func aggregateToWeekly(dailyData []database.GetURLDailyClicksRow) []WeeklyClick {
	if len(dailyData) == 0 {
		return []WeeklyClick{}
	}

	weeklyMap := make(map[string]int64)

	for _, dc := range dailyData {
		year, week := dc.Date.ISOWeek()
		weekStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		weekStart = weekStart.AddDate(0, 0, (week-1)*7)

		for weekStart.Weekday() != time.Monday {
			weekStart = weekStart.AddDate(0, 0, -1)
		}

		weekKey := weekStart.Format("2006-01-02")
		weeklyMap[weekKey] += int64(dc.TotalClicks)
	}

	var weekly []WeeklyClick
	for weekStart, clicks := range weeklyMap {
		weekly = append(weekly, WeeklyClick{
			WeekStart: weekStart,
			Clicks:    clicks,
		})
	}

	return weekly
}

func aggregateToMonthly(dailyData []database.GetURLDailyClicksRow) []MonthlyClick {
	if len(dailyData) == 0 {
		return []MonthlyClick{}
	}

	monthlyMap := make(map[string]int64)

	for _, dc := range dailyData {
		monthKey := dc.Date.Format("2006-01")
		monthlyMap[monthKey] += int64(dc.TotalClicks)
	}

	var monthly []MonthlyClick
	for month, clicks := range monthlyMap {
		monthly = append(monthly, MonthlyClick{
			Month:  month,
			Clicks: clicks,
		})
	}

	return monthly
}
