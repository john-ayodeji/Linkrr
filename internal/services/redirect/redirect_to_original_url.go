package redirect

import (
	"fmt"
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/database"
)

var Cfg *config.ApiConfig

func Redirect(r *http.Request) (string, error, int) {
	urlCode := r.PathValue("urlCode")
	if urlCode == "" {
		return "", fmt.Errorf("page not found"), http.StatusNotFound
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
