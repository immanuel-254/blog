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

const CategoryRouteGroup = "/category"

var (
	CategoryCreateView = View{
		Route:   fmt.Sprintf("%s/create", CategoryRouteGroup),
		Handler: http.HandlerFunc(CategoryCreate),
		Methods: []string{http.MethodPost},
	}

	CategoryReadView = View{
		Route:   fmt.Sprintf("%s/read/", CategoryRouteGroup),
		Handler: http.HandlerFunc(CategoryRead),
		Methods: []string{http.MethodGet},
	}

	CategoryListView = View{
		Route:   fmt.Sprintf("%s/list", CategoryRouteGroup),
		Handler: http.HandlerFunc(CategoryList),
		Methods: []string{http.MethodGet},
	}

	CategoryUpdateView = View{
		Route:   fmt.Sprintf("%s/update/", CategoryRouteGroup),
		Handler: http.HandlerFunc(CategoryUpdate),
		Methods: []string{http.MethodPut},
	}

	CategoryDeleteView = View{
		Route:   fmt.Sprintf("%s/delete/", CategoryRouteGroup),
		Handler: http.HandlerFunc(CategoryDelete),
		Methods: []string{http.MethodDelete},
	}
)

func CategoryCreate(w http.ResponseWriter, r *http.Request) {
	// Entities To be Created; Category
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

	category, err := queries.CategoryCreate(ctx, models.CategoryCreateParams{
		UserID:    sql.NullInt64{Int64: input.userid, Valid: true},
		Name:      data["name"].(string),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Category models.Category `json:"category"`
	}

	output.Category = category

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CategoryRead(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/read/", CategoryRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Entities To Read; Category
	queries := models.New(database.DB)
	ctx := context.Background()

	category, err := queries.CategoryRead(ctx, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Category models.CategoryReadRow `json:"category"`
	}

	output.Category = category

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CategoryList(w http.ResponseWriter, r *http.Request) {
	// Entities To List; Category
	queries := models.New(database.DB)
	ctx := context.Background()

	categories, err := queries.CategoryBlogList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Categories []models.CategoryBlogListRow `json:"categories"`
	}

	output.Categories = categories

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/update/", CategoryRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Entities To Update; Category
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	category, err := queries.CategoryUpdate(ctx, models.CategoryUpdateParams{
		ID:        int64(id),
		Name:      data["name"],
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Category models.CategoryUpdateRow `json:"category"`
	}

	output.Category = category

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func CategoryDelete(w http.ResponseWriter, r *http.Request) {
	// Entities To Delete; Category
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/delete/", CategoryRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	err = queries.CategoryDelete(ctx, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// get all categories
	categories, err := queries.CategoryBlogList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// delete many to many relations
	for _, category := range categories {
		if category.CategoryID == int64(id) {
			err = queries.CategoryBlogDelete(ctx, models.CategoryBlogDeleteParams{
				BlogID:     category.BlogID,
				CategoryID: int64(id),
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
