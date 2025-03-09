package views

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	auth "github.com/immanuel-254/blog/auth/models"
	"github.com/immanuel-254/blog/blog/models"
	"github.com/immanuel-254/blog/database"
)

const ProfileRouteGroup = "/profile"

var (
	ProfileReadView = View{
		Route:   fmt.Sprintf("%s/read/", ProfileRouteGroup),
		Handler: http.HandlerFunc(ProfileRead),
		Methods: []string{http.MethodGet},
	}

	ProfileListView = View{
		Route:   fmt.Sprintf("%s/list", ProfileRouteGroup),
		Handler: http.HandlerFunc(ProfileList),
		Methods: []string{http.MethodGet},
	}

	ProfileUpdateView = View{
		Route:   fmt.Sprintf("%s/update/", ProfileRouteGroup),
		Handler: http.HandlerFunc(ProfileUpdate),
		Methods: []string{http.MethodPut},
	}
)

func ProfileRead(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/read/", ProfileRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	// Entities To Read; Profile, User
	authqueries := auth.New(database.DB)
	queries := models.New(database.DB)
	ctx := context.Background()

	profile, err := queries.ProfileRead(ctx, int64(id))

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	user, err := authqueries.UserRead(ctx, profile.UserID.Int64)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	var output struct {
		User    auth.UserReadRow `json:"user"`
		Profile models.Profile   `json:"profile"`
	}

	output.User = user
	output.Profile = profile

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func ProfileList(w http.ResponseWriter, r *http.Request) {
	authqueries := auth.New(database.DB)
	queries := models.New(database.DB)
	ctx := context.Background()

	// Fetch user list
	userlist, err := authqueries.UserList(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Fetch profile list
	profilelist, err := queries.ProfileList(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Combine users and profiles
	type useroutput struct {
		User    auth.UserListRow `json:"user"`
		Profile models.Profile   `json:"profile"`
	}
	var output []useroutput

	for _, profile := range profilelist {
		for _, user := range userlist {
			if profile.UserID.Int64 == user.ID {
				output = append(output, useroutput{
					User:    user,
					Profile: profile,
				})
			}
		}
	}

	// Respond with the combined output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/update/", ProfileRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
		return
	}

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Entities To Update; User
	queries := models.New(database.DB)
	ctx := context.Background()

	profile, err := queries.ProfileUpdate(ctx, models.ProfileUpdateParams{
		ID:        int64(id),
		Username:  data["username"],
		Image:     sql.NullString{String: data["image"], Valid: true},
		Bio:       sql.NullString{String: data["bio"], Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	var output struct {
		Profile models.Profile `json:"profile"`
	}

	output.Profile = profile

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
