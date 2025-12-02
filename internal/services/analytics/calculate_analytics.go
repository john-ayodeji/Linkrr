package analytics

import (
	"context"
	"log"
	"time"

	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

// AggregateAnalytics listens on the channel and aggregates analytics data into the DB
func AggregateAnalytics(urlAlias <-chan AnalyticsData) {
	for data := range urlAlias {
		go processAnalyticsEvent(data)
	}
}

func processAnalyticsEvent(data AnalyticsData) {
	ctx := context.Background()

	// Get the user_id who owns this URL
	userID, err := Cfg.Db.GetUrlOwnerByShortCode(ctx, data.ShortCode)
	if err != nil {
		utils.LogError("failed to get URL owner: " + err.Error())
		return
	}

	today := time.Now().Truncate(24 * time.Hour)

	// === URL-level aggregations ===
	// Use incremental updates (add 1 per click) to avoid race conditions

	// Daily clicks per URL
	if err := Cfg.Db.UpsertURLDaily(ctx, database.UpsertURLDailyParams{
		ShortCode:      data.ShortCode,
		Date:           today,
		TotalClicks:    1,
		UniqueVisitors: 1,
	}); err != nil {
		log.Printf("failed to upsert URL daily: %v", err)
	}

	// Geo data per URL
	if err := Cfg.Db.UpsertURLGeo(ctx, database.UpsertURLGeoParams{
		ShortCode:   data.ShortCode,
		Country:     data.Country,
		City:        data.City,
		TotalClicks: 1,
	}); err != nil {
		log.Printf("failed to upsert URL geo: %v", err)
	}

	// Device per URL
	if err := Cfg.Db.UpsertURLDevice(ctx, database.UpsertURLDeviceParams{
		ShortCode:   data.ShortCode,
		Device:      data.Device,
		TotalClicks: 1,
	}); err != nil {
		log.Printf("failed to upsert URL device: %v", err)
	}

	// Browser per URL
	if err := Cfg.Db.UpsertURLBrowser(ctx, database.UpsertURLBrowserParams{
		ShortCode:   data.ShortCode,
		Browser:     data.Browser,
		TotalClicks: 1,
	}); err != nil {
		log.Printf("failed to upsert URL browser: %v", err)
	}

	// Referrer per URL
	if err := Cfg.Db.UpsertURLReferrer(ctx, database.UpsertURLReferrerParams{
		ShortCode:   data.ShortCode,
		Referrer:    data.Referrer,
		TotalClicks: 1,
	}); err != nil {
		log.Printf("failed to upsert URL referrer: %v", err)
	}

	// === Alias-level aggregations (if alias exists) ===
	if data.Alias != "" {
		if err := Cfg.Db.UpsertAliasDaily(ctx, database.UpsertAliasDailyParams{
			ShortCode:      data.ShortCode,
			Alias:          data.Alias,
			Date:           today,
			TotalClicks:    1,
			UniqueVisitors: 1,
		}); err != nil {
			log.Printf("failed to upsert alias daily: %v", err)
		}

		if err := Cfg.Db.UpsertAliasGeo(ctx, database.UpsertAliasGeoParams{
			ShortCode:   data.ShortCode,
			Alias:       data.Alias,
			Country:     data.Country,
			City:        data.City,
			TotalClicks: 1,
		}); err != nil {
			log.Printf("failed to upsert alias geo: %v", err)
		}

		if err := Cfg.Db.UpsertAliasDevice(ctx, database.UpsertAliasDeviceParams{
			ShortCode:   data.ShortCode,
			Alias:       data.Alias,
			Device:      data.Device,
			TotalClicks: 1,
		}); err != nil {
			log.Printf("failed to upsert alias device: %v", err)
		}

		if err := Cfg.Db.UpsertAliasBrowser(ctx, database.UpsertAliasBrowserParams{
			ShortCode:   data.ShortCode,
			Alias:       data.Alias,
			Browser:     data.Browser,
			TotalClicks: 1,
		}); err != nil {
			log.Printf("failed to upsert alias browser: %v", err)
		}

		if err := Cfg.Db.UpsertAliasReferrer(ctx, database.UpsertAliasReferrerParams{
			ShortCode:   data.ShortCode,
			Alias:       data.Alias,
			Referrer:    data.Referrer,
			TotalClicks: 1,
		}); err != nil {
			log.Printf("failed to upsert alias referrer: %v", err)
		}
	}

	// === User-level aggregations ===

	if err := Cfg.Db.UpsertUserDaily(ctx, database.UpsertUserDailyParams{
		UserID:         userID,
		Date:           today,
		TotalClicks:    1,
		UniqueVisitors: 1,
	}); err != nil {
		log.Printf("failed to upsert user daily: %v", err)
	}

	if err := Cfg.Db.UpsertUserGeo(ctx, database.UpsertUserGeoParams{
		UserID:      userID,
		Country:     data.Country,
		TotalClicks: 1,
	}); err != nil {
		log.Printf("failed to upsert user geo: %v", err)
	}

	if err := Cfg.Db.UpsertUserBrowser(ctx, database.UpsertUserBrowserParams{
		UserID:      userID,
		Browser:     data.Browser,
		TotalClicks: 1,
	}); err != nil {
		log.Printf("failed to upsert user browser: %v", err)
	}

	if err := Cfg.Db.UpsertUserTopLinks(ctx, database.UpsertUserTopLinksParams{
		UserID:      userID,
		ShortCode:   data.ShortCode,
		TotalClicks: 1,
	}); err != nil {
		log.Printf("failed to upsert user top links: %v", err)
	}
}
