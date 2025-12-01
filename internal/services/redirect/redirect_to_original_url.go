package redirect

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/database"
)

var Cfg *config.ApiConfig

type URLData struct {
	IP        string
	Referer   string
	UserAgent string
	ClickedAt time.Time
	UrlCode   string
}

var RedirectEvent = make(chan URLData, 100)

func Redirect(r *http.Request) (string, error, int) {
	var ip string
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ip = strings.Split(ip, ",")[0]
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	referrer := r.Header.Get("Referer")
	userAgent := r.Header.Get("User-Agent")
	clickedAt := time.Now().UTC()

	urlCode := r.PathValue("urlCode")
	if urlCode == "" {
		return "", fmt.Errorf("page not found"), http.StatusNotFound
	}

	RedirectEvent <- URLData{
		IP:        ip,
		Referer:   referrer,
		UserAgent: userAgent,
		ClickedAt: clickedAt,
		UrlCode:   urlCode,
	}

	url, err := Cfg.Db.GetOriginalUrl(r.Context(), database.GetOriginalUrlParams{
		ShortCode: urlCode,
		Alias:     urlCode,
	})
	if err != nil {
		return "", fmt.Errorf("page not foundx"), http.StatusNotFound
	}

	return url, nil, http.StatusPermanentRedirect
}
