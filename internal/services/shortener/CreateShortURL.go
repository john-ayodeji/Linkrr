package shortener

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

var Cfg *config.ApiConfig

type Request struct {
	LongURL string `json:"long_url"`
}

type Response struct {
	ID        uuid.UUID `json:"id"`
	ShortUrl  string    `json:"short_url"`
	LongUrl   string    `json:"long_url"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateShortURL(r *http.Request) (Response, error, int) {
	jwtToken, err := utils.GetBearerToken(r.Header)
	if err != nil {
		return Response{}, err, http.StatusUnauthorized
	}
	ok, claims, err := auth.ValidateJWTHelper(jwtToken, Cfg.JWTSecret)
	if !ok {
		return Response{}, fmt.Errorf("%v", err), http.StatusUnauthorized
	}

	var req Request
	decoded := json.NewDecoder(r.Body)
	if err := decoded.Decode(&req); err != nil {
		return Response{}, fmt.Errorf("invalid request body"), http.StatusInternalServerError
	}

	input := strings.TrimSpace(req.LongURL)
	if input == "" {
		return Response{}, fmt.Errorf("empty url"), http.StatusBadRequest
	}
	if !strings.Contains(input, "://") {
		input = "https://" + input
	}

	url_short_code, err := RandomIDwithRetry(r)
	if err != nil {
		return Response{}, err, http.StatusInternalServerError
	}

	userid, _ := uuid.Parse(claims.Subject)

	urlData, err := Cfg.Db.CreateURL(r.Context(), database.CreateURLParams{
		ID:        uuid.New(),
		ShortCode: url_short_code,
		Url:       input,
		UserID:    userid,
	})

	resp := Response{
		ID:        urlData.ID,
		ShortUrl:  r.Host + "/" + urlData.ShortCode,
		LongUrl:   urlData.Url,
		UserID:    urlData.UserID,
		CreatedAt: urlData.CreatedAt,
	}

	return resp, nil, http.StatusCreated
}
