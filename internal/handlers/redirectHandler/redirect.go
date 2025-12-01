package redirectHandler

import (
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/services/redirect"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	url, err, statusCode := redirect.Redirect(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	http.Redirect(w, r, url, statusCode)
}
