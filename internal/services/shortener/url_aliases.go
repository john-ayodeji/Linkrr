package shortener

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

type AliasRequest struct {
	URLCode string `json:"url_code"`
	Alias   string `json:"alias"`
}

type AliasResponse struct {
	ID          uuid.UUID `json:"id"`
	Alias       string    `json:"alias"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateAlias(r *http.Request) (AliasResponse, error, int) {
	jwtToken, err := utils.GetBearerToken(r.Header)
	if err != nil {
		return AliasResponse{}, err, http.StatusUnauthorized
	}
	ok, _, err := auth.ValidateJWTHelper(jwtToken, Cfg.JWTSecret)
	if !ok {
		return AliasResponse{}, err, http.StatusUnauthorized
	}

	var aliasData AliasRequest
	decoded := json.NewDecoder(r.Body)
	if err := decoded.Decode(&aliasData); err != nil {
		return AliasResponse{}, fmt.Errorf("fill all input fields"), http.StatusInternalServerError
	}

	_, err1 := Cfg.Db.GetURL(r.Context(), aliasData.URLCode)
	if err1 != nil {
		return AliasResponse{}, fmt.Errorf("url does not exist"), http.StatusNotFound
	}

	data, err := Cfg.Db.CreateAlias(r.Context(), database.CreateAliasParams{
		ID:      uuid.New(),
		Alias:   aliasData.Alias,
		UrlCode: aliasData.URLCode,
	})
	if err != nil {
		return AliasResponse{}, fmt.Errorf("alias not available"), http.StatusInternalServerError
	}

	resp := AliasResponse{
		ID:          data.ID,
		Alias:       r.Host + "/" + data.Alias,
		OriginalURL: r.Host + "/" + data.UrlCode,
		CreatedAt:   data.CreatedAt,
	}

	return resp, nil, http.StatusCreated
}
