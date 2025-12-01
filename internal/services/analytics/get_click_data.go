package analytics

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/services/redirect"
	"github.com/john-ayodeji/Linkrr/internal/utils"
	"github.com/mssola/useragent"
)

var Cfg *config.ApiConfig

func GetClickData(urldata <-chan redirect.URLData) {
	for data := range urldata {
		codeAlias, err := Cfg.Db.GetShortCodeAndAlias(context.Background(), data.UrlCode)
		if err != nil {
			utils.LogError("failed to get shortcode & alias")
		}
		ua := useragent.New(data.UserAgent)
		browser, _ := ua.Browser()
		country, _ := GetIpLocation(data.IP)

		err1 := Cfg.Db.CreateAnalyticsData(context.Background(), database.CreateAnalyticsDataParams{
			ID:        uuid.New(),
			ShortCode: codeAlias.Code,
			Alias:     codeAlias.Alias,
			ClickedAt: data.ClickedAt,
			Ip:        data.IP,
			Country:   country,
			Referrer:  sql.NullString{String: data.Referer},
			Device:    ua.Platform(),
			Os:        ua.OS(),
			Browser:   browser,
		})
		if err1 != nil {
			utils.LogError("failed to save analytics data")
		}
	}
}

func GetIpLocation(ip string) (string, error) {
	type location struct {
		CountryCode string `json:"country_name"`
	}
	req, err := http.NewRequest("GET", Cfg.IpStackURl+ip+"?access_key="+Cfg.IpStackApiKey+"&fields=country_name", nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch location")
	}

	client := http.Client{
		Timeout: 2 * time.Minute,
	}

	res, err := client.Do(req)

	defer res.Body.Close()

	var country location
	decoded := json.NewDecoder(res.Body)
	if err := decoded.Decode(&country); err != nil {
		return "", fmt.Errorf("failed to decode response body")
	}

	return country.CountryCode, nil
}
