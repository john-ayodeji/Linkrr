package analytics

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
)

type TopLink struct {
	ShortCode string `json:"short_code"`
	Clicks    int64  `json:"clicks"`
}

type GlobalAnalytics struct {
	TotalClicks    int64          `json:"total_clicks"`
	UniqueVisitors int64          `json:"unique_visitors"`
	DailyClicks    []DailyClick   `json:"daily_clicks"`
	TopCountries   []CountryClick `json:"top_countries"`
	TopBrowsers    []BrowserClick `json:"top_browsers"`
	TopLinks       []TopLink      `json:"top_links"`
}

func GetGlobalAnalytics(r *http.Request) (GlobalAnalytics, error, int) {
	ctx := context.Background()

	authHeader := r.Header.Get("Authorization")
	parts := strings.Fields(authHeader)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return GlobalAnalytics{}, fmt.Errorf("invalid authorization header"), http.StatusUnauthorized
	}

	token := parts[1]
	ok, claims, err := auth.ValidateJWTHelper(token, Cfg.JWTSecret)
	if !ok {
		return GlobalAnalytics{}, err, http.StatusUnauthorized
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return GlobalAnalytics{}, fmt.Errorf("something went wrong"), http.StatusInternalServerError
	}

	analytics := GlobalAnalytics{}

	totalClicksResult, err := Cfg.Db.GetUserTotalClicks(ctx, userID)
	if err != nil {
		return GlobalAnalytics{}, fmt.Errorf("failed to get total clicks"), http.StatusInternalServerError
	}
	analytics.TotalClicks = totalClicksResult.(int64)

	uniqueVisitorsResult, err := Cfg.Db.GetUserUniqueVisitors(ctx, userID)
	if err != nil {
		return GlobalAnalytics{}, fmt.Errorf("failed to get unique visitors"), http.StatusInternalServerError
	}
	analytics.UniqueVisitors = uniqueVisitorsResult.(int64)

	dailyClicks, err := Cfg.Db.GetUserDailyClicks(ctx, userID)
	if err != nil {
		return GlobalAnalytics{}, fmt.Errorf("failed to get daily clicks"), http.StatusInternalServerError
	}
	for _, dc := range dailyClicks {
		analytics.DailyClicks = append(analytics.DailyClicks, DailyClick{
			Date:   dc.Date.Format("2006-01-02"),
			Clicks: int64(dc.TotalClicks),
		})
	}

	countryClicks, err := Cfg.Db.GetUserGeoClicks(ctx, userID)
	if err != nil {
		return GlobalAnalytics{}, fmt.Errorf("failed to get country clicks"), http.StatusInternalServerError
	}
	for _, cc := range countryClicks {
		analytics.TopCountries = append(analytics.TopCountries, CountryClick{
			Country: cc.Country,
			Clicks:  int64(cc.TotalClicks),
		})
	}

	browserClicks, err := Cfg.Db.GetUserBrowserClicks(ctx, userID)
	if err != nil {
		return GlobalAnalytics{}, fmt.Errorf("failed to get browser clicks"), http.StatusInternalServerError
	}
	for _, browser := range browserClicks {
		analytics.TopBrowsers = append(analytics.TopBrowsers, BrowserClick{
			Browser: browser.Browser,
			Clicks:  int64(browser.TotalClicks),
		})
	}

	topLinks, err := Cfg.Db.GetUserTopLinks(ctx, userID)
	if err != nil {
		return GlobalAnalytics{}, fmt.Errorf("failed to get top links"), http.StatusInternalServerError
	}
	for _, link := range topLinks {
		analytics.TopLinks = append(analytics.TopLinks, TopLink{
			ShortCode: link.ShortCode,
			Clicks:    int64(link.TotalClicks),
		})
	}

	return analytics, nil, http.StatusOK
}
