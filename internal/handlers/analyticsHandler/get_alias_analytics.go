package analyticsHandler

import (
	"encoding/json"
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/services/analytics"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

func GetAliasAnalytics(w http.ResponseWriter, r *http.Request) {
	aliasAnalytics, err, statusCode := analytics.GetAliasAnalytics(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(aliasAnalytics)
}
