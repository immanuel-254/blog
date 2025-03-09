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

	"github.com/immanuel-254/blog/blog/models"
	"github.com/immanuel-254/blog/database"
)

const CommentRouteGroup = "/comment"

var (
	CommentCreateView = View{
		Route:   fmt.Sprintf("%s/create", CommentRouteGroup),
		Handler: http.HandlerFunc(CommentCreate),
		Methods: []string{http.MethodPost},
	}

	CommentReadView = View{
		Route:   fmt.Sprintf("%s/read/", CommentRouteGroup),
		Handler: http.HandlerFunc(CommentRead),
		Methods: []string{http.MethodGet},
	}

	CommentListView = View{
		Route:   fmt.Sprintf("%s/list", CommentRouteGroup),
		Handler: http.HandlerFunc(CommentList),
		Methods: []string{http.MethodGet},
	}

	CommentUpdateView = View{
		Route:   fmt.Sprintf("%s/update/", CommentRouteGroup),
		Handler: http.HandlerFunc(CommentUpdate),
		Methods: []string{http.MethodPut},
	}

	CommentDeleteView = View{
		Route:   fmt.Sprintf("%s/delete/", CommentRouteGroup),
		Handler: http.HandlerFunc(CommentDelete),
		Methods: []string{http.MethodDelete},
	}
)

func CommentCreate(w http.ResponseWriter, r *http.Request) {
	// Entities To be Created; Comment
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	var input struct {
		userid int64
	}

	if _, ok := data["userid"].(string); ok {
		id, err := strconv.ParseInt(data["userid"].(string), 10, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input.userid = id
	}

	if _, ok := data["userid"].(float64); ok {
		input.userid = int64(data["userid"].(float64))
	}

	comment, err := queries.CommentCreate(ctx, models.CommentCreateParams{
		UserID:    sql.NullInt64{Int64: input.userid, Valid: true},
		Body:      data["body"].(string),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Comment models.Comment `json:"comment"`
	}

	output.Comment = comment

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CommentRead(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/read/", CategoryRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Entities To Read; Comment
	queries := models.New(database.DB)
	ctx := context.Background()

	comment, err := queries.CommentRead(ctx, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Comment models.CommentReadRow `json:"comment"`
	}

	output.Comment = comment

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CommentList(w http.ResponseWriter, r *http.Request) {
	// Entities To List; Comment
	queries := models.New(database.DB)
	ctx := context.Background()

	comments, err := queries.CommentList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Comments []models.CommentListRow `json:"comments"`
	}

	output.Comments = comments

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CommentUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/update/", CategoryRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Entities To Update; Comment
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	comment, err := queries.CommentUpdate(ctx, models.CommentUpdateParams{
		ID:        int64(id),
		Body:      data["body"],
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Comment models.CommentUpdateRow `json:"comment"`
	}

	output.Comment = comment

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CommentDelete(w http.ResponseWriter, r *http.Request) {
	// Entities To Delete; Comment
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/delete/", CategoryRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	err = queries.CommentDelete(ctx, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
