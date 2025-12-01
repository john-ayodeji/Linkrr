package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/config"
)

var Cfg *config.ApiConfig

type user struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type url struct {
	ShortURL string   `json:"short_url"`
	LongURL  string   `json:"long_url"`
	Aliases  []string `json:"aliases"`
}

type MyLinks struct {
	User user  `json:"user"`
	Data []url `json:"data"`
}

func GetMyLinks(r *http.Request) (MyLinks, error, int) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Fields(authHeader)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return MyLinks{}, fmt.Errorf("invalid authorization header"), http.StatusUnauthorized
	}

	token := parts[1]
	ok, claims, err := auth.ValidateJWTHelper(token, Cfg.JWTSecret)
	if !ok {
		return MyLinks{}, err, http.StatusUnauthorized
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return MyLinks{}, fmt.Errorf("something went wrong"), http.StatusInternalServerError
	}

	data, err := Cfg.Db.GetURLsForUser(r.Context(), userId)
	if err != nil {
		return MyLinks{}, fmt.Errorf("no urls found"), http.StatusNotFound
	}

	if len(data) == 0 {
		return MyLinks{}, fmt.Errorf("no urls found"), http.StatusNotFound
	}

	userInfo := user{
		Id:   data[0].UserID,
		Name: data[0].Name,
	}

	var urlInfo []url

	for _, item := range data {
		aliases := []string{}
		if item.Alias.Valid {
			aliases = append(aliases, item.Alias.String)
		}

		urlInfo = append(urlInfo, url{
			ShortURL: r.Host + "/" + item.ShortCode,
			LongURL:  item.Url,
			Aliases:  aliases,
		})
	}

	return MyLinks{
		User: userInfo,
		Data: urlInfo,
	}, nil, http.StatusOK
}
