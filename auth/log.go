package auth

import (
	"net/http"

	"github.com/immanuel-254/blog/auth/models"
)

func LogList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queries := models.New(DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	logs, err := queries.LogList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "log", "list", 0, authUser.ID, w, r)

	SendData(map[string]interface{}{"logs": logs}, w, r)
}
