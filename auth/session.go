package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/immanuel-254/blog/auth/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	// get data
	var data map[string]string
	GetData(data, w, r)

	key, code, err := AuthLogin(queries, ctx, data)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	resp := map[string]interface{}{"auth": key}
	SendData(resp, w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queries := models.New(DB)
	ctx := r.Context()

	session, err := queries.SessionRead(ctx, w.Header().Get("auth"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// delete session
	err = queries.SessionDelete(ctx, w.Header().Get("auth"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = queries.LogCreate(ctx, models.LogCreateParams{
		DbTable:   "session",
		Action:    "delete",
		ObjectID:  session.ID,
		UserID:    session.UserID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"message": "user logged out"}
	SendData(resp, w, r)
}

func SessionList(w http.ResponseWriter, r *http.Request) {
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

	sessions, err := queries.SessionList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "session", "list", 0, authUser.ID, w, r)

	SendData(map[string]interface{}{"sessions": sessions}, w, r)
}
