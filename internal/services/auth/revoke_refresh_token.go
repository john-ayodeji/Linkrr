package authService

import (
	"fmt"
	"net/http"

	"github.com/john-ayodeji/Linkrr/utils"
)

type Status struct {
	Status string `json:"status"`
}

func RevokeRefreshToken(r *http.Request) (Status, error, int) {
	token, err := utils.GetBearerToken(r.Header)
	if err != nil {
		return Status{}, err, http.StatusNotFound
	}

	_, err1 := Cfg.Db.VerifyRefreshToken(r.Context(), token)
	if err1 != nil {
		return Status{}, fmt.Errorf("invalid refresh token"), http.StatusUnauthorized
	}

	Err := Cfg.Db.RevokeRefreshToken(r.Context(), token)
	if Err != nil {
		return Status{}, fmt.Errorf("something went wrong"), http.StatusInternalServerError
	}

	return Status{Status: "success"}, nil, http.StatusOK
}
