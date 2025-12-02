package analyticsHandler

import (
	"encoding/json"
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/services/analytics"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

func GetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	urlAnalytics, err, statusCode := analytics.GetURLAnalytics(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urlAnalytics)
}
