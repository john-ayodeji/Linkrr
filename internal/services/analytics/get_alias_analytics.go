package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/database"
)

type AliasAnalytics struct {
	ShortCode      string          `json:"short_code"`
	Alias          string          `json:"alias"`
	TotalClicks    int64           `json:"total_clicks"`
	UniqueVisitors int64           `json:"unique_visitors"`
	DailyClicks    []DailyClick    `json:"daily_clicks"`
	TopCountries   []CountryClick  `json:"top_countries"`
	TopCities      []CityClick     `json:"top_cities"`
	TopDevices     []DeviceClick   `json:"top_devices"`
	TopBrowsers    []BrowserClick  `json:"top_browsers"`
	TopReferrers   []ReferrerClick `json:"top_referrers"`
}

func GetAliasAnalytics(r *http.Request) (AliasAnalytics, error, int) {
	ctx := context.Background()
	shortCode := r.PathValue("urlCode")
	alias := r.PathValue("alias")

	if shortCode == "" || alias == "" {
		return AliasAnalytics{}, fmt.Errorf("url code and alias are required"), http.StatusBadRequest
	}

	authHeader := r.Header.Get("Authorization")
	parts := strings.Fields(authHeader)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return AliasAnalytics{}, fmt.Errorf("invalid authorization header"), http.StatusUnauthorized
	}

	token := parts[1]
	ok, claims, err := auth.ValidateJWTHelper(token, Cfg.JWTSecret)
	if !ok {
		return AliasAnalytics{}, err, http.StatusUnauthorized
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("something went wrong"), http.StatusInternalServerError
	}

	// Verify user owns this URL
	urlOwner, err := Cfg.Db.GetUrlOwnerByShortCode(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return AliasAnalytics{}, fmt.Errorf("URL not found"), http.StatusNotFound
		}
		return AliasAnalytics{}, fmt.Errorf("failed to verify URL ownership"), http.StatusInternalServerError
	}

	if urlOwner != userID {
		return AliasAnalytics{}, fmt.Errorf("forbidden"), http.StatusForbidden
	}

	analytics := AliasAnalytics{
		ShortCode: shortCode,
		Alias:     alias,
	}

	totalClicksResult, err := Cfg.Db.GetAliasTotalClicks(ctx, database.GetAliasTotalClicksParams{
		ShortCode: shortCode,
		Alias:     alias,
	})
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("failed to get total clicks"), http.StatusInternalServerError
	}
	analytics.TotalClicks = totalClicksResult.(int64)

	uniqueVisitorsResult, err := Cfg.Db.GetAliasUniqueVisitors(ctx, database.GetAliasUniqueVisitorsParams{
		ShortCode: shortCode,
		Alias:     alias,
	})
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("failed to get unique visitors"), http.StatusInternalServerError
	}
	analytics.UniqueVisitors = uniqueVisitorsResult.(int64)

	dailyClicks, err := Cfg.Db.GetAliasDailyClicks(ctx, database.GetAliasDailyClicksParams{
		ShortCode: shortCode,
		Alias:     alias,
	})
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("failed to get daily clicks"), http.StatusInternalServerError
	}
	for _, dc := range dailyClicks {
		analytics.DailyClicks = append(analytics.DailyClicks, DailyClick{
			Date:   dc.Date.Format("2006-01-02"),
			Clicks: int64(dc.TotalClicks),
		})
	}

	geoClicks, err := Cfg.Db.GetAliasGeoClicks(ctx, database.GetAliasGeoClicksParams{
		ShortCode: shortCode,
		Alias:     alias,
	})
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("failed to get geo clicks"), http.StatusInternalServerError
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

	deviceClicks, err := Cfg.Db.GetAliasDeviceClicks(ctx, database.GetAliasDeviceClicksParams{
		ShortCode: shortCode,
		Alias:     alias,
	})
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("failed to get device clicks"), http.StatusInternalServerError
	}
	for _, device := range deviceClicks {
		analytics.TopDevices = append(analytics.TopDevices, DeviceClick{
			Device: device.Device,
			Clicks: int64(device.TotalClicks),
		})
	}

	browserClicks, err := Cfg.Db.GetAliasBrowserClicks(ctx, database.GetAliasBrowserClicksParams{
		ShortCode: shortCode,
		Alias:     alias,
	})
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("failed to get browser clicks"), http.StatusInternalServerError
	}
	for _, browser := range browserClicks {
		analytics.TopBrowsers = append(analytics.TopBrowsers, BrowserClick{
			Browser: browser.Browser,
			Clicks:  int64(browser.TotalClicks),
		})
	}

	referrerClicks, err := Cfg.Db.GetAliasReferrerClicks(ctx, database.GetAliasReferrerClicksParams{
		ShortCode: shortCode,
		Alias:     alias,
	})
	if err != nil {
		return AliasAnalytics{}, fmt.Errorf("failed to get referrer clicks"), http.StatusInternalServerError
	}
	for _, ref := range referrerClicks {
		analytics.TopReferrers = append(analytics.TopReferrers, ReferrerClick{
			Referrer: ref.Referrer,
			Clicks:   int64(ref.TotalClicks),
		})
	}

	return analytics, nil, http.StatusOK
}
