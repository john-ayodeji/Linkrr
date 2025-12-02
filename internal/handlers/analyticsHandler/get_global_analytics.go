package analyticsHandler

import (
	"encoding/json"
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/services/analytics"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

func GetGlobalAnalytics(w http.ResponseWriter, r *http.Request) {
	globalAnalytics, err, statusCode := analytics.GetGlobalAnalytics(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(globalAnalytics)
}
