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

	"github.com/immanuel-254/blog/auth"
	"github.com/immanuel-254/blog/blog/models"
	"github.com/immanuel-254/blog/database"
)

const BlogRouteGroup = "/blog"

var (
	BlogCreateView = View{
		Route:   fmt.Sprintf("%s/create", BlogRouteGroup),
		Handler: http.HandlerFunc(BlogCreate),
		Methods: []string{http.MethodPost},
	}

	BlogReadView = View{
		Route:   fmt.Sprintf("%s/read/", BlogRouteGroup),
		Handler: http.HandlerFunc(BlogRead),
		Methods: []string{http.MethodGet},
	}

	BlogListView = View{
		Route:   fmt.Sprintf("%s/list", BlogRouteGroup),
		Handler: http.HandlerFunc(BlogList),
		Methods: []string{http.MethodGet},
	}

	BlogUpdateView = View{
		Route:   fmt.Sprintf("%s/update/", BlogRouteGroup),
		Handler: http.HandlerFunc(BlogUpdate),
		Methods: []string{http.MethodPut},
	}

	BlogDeleteView = View{
		Route:   fmt.Sprintf("%s/delete/", BlogRouteGroup),
		Handler: http.HandlerFunc(BlogDelete),
		Methods: []string{http.MethodDelete},
	}
)

func BlogCreate(w http.ResponseWriter, r *http.Request) {
	// Entities To be Created; Blog
	var data map[string]string
	auth.GetData(data, w, r)

	queries := models.New(database.DB)
	ctx := context.Background()

	var input struct {
		userid int64
	}

	if _, ok := data["userid"]; ok {
		id, err := strconv.ParseInt(data["userid"], 10, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input.userid = id
	}

	blog, err := queries.BlogCreate(ctx, models.BlogCreateParams{
		UserID:    sql.NullInt64{Int64: input.userid, Valid: true},
		Title:     data["title"],
		Body:      data["body"],
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var categories []int64

	if _, ok := data["categories"]; ok {

		rawCategories := data["categories"]
		parts := strings.Split(rawCategories, ",")

		for _, part := range parts {
			trimmedPart := strings.TrimSpace(part) // Remove any extra spaces
			if part == "" {
				continue // Skip empty parts (e.g., due to a trailing comma)
			}
			num, err := strconv.ParseInt(trimmedPart, 10, 64) // Convert to int64
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			categories = append(categories, num)
		}
	}

	rawCategories := []interface{}{}

	if _, ok := data["categories"]; ok {
		err := json.Unmarshal([]byte(data["categories"]), &rawCategories)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, v := range rawCategories {
			if num, ok := v.(float64); ok {
				categories[i] = int64(num)
			}
		}
	}

	for _, value := range categories {
		err = queries.AssignBlogToCategory(ctx, models.AssignBlogToCategoryParams{
			BlogID:     blog.ID,
			CategoryID: value,
			CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var output struct {
		Blog models.Blog `json:"blog"`
	}

	output.Blog = blog

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func BlogRead(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/read/", BlogRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Entities To Read; Blog
	queries := models.New(database.DB)
	ctx := context.Background()

	blog, err := queries.BlogRead(ctx, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Blog models.BlogReadRow `json:"blog"`
	}

	output.Blog = blog

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func BlogList(w http.ResponseWriter, r *http.Request) {
	// Entities To Read; Blog, Category
	queries := models.New(database.DB)
	ctx := context.Background()

	blogs, err := queries.BlogList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Blogs []models.BlogListRow `json:"blogs"`
	}

	output.Blogs = blogs

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func BlogUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/update/", BlogRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Entities To Update; Blog, Category
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	blog, err := queries.BlogUpdate(ctx, models.BlogUpdateParams{
		ID:        int64(id),
		Title:     data["title"].(string),
		Body:      data["body"].(string),
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output struct {
		Blog models.BlogUpdateRow `json:"blog"`
	}

	output.Blog = blog

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func BlogDelete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/delete/", BlogRouteGroup))
	idStr = strings.TrimLeft(idStr, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Entities To Delete; Blog
	queries := models.New(database.DB)
	ctx := context.Background()

	err = queries.BlogDelete(ctx, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get blog categories
	categories, err := queries.BlogCategoriesList(ctx, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get blog comments
	comments, err := queries.BlogCommentsList(ctx, sql.NullInt64{Int64: int64(id), Valid: true})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// delete many to many relations
	for _, category := range categories {
		err = queries.CategoryBlogDelete(ctx, models.CategoryBlogDeleteParams{
			BlogID:     int64(id),
			CategoryID: category.CategoryID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	for _, comment := range comments {
		err = queries.CommentDelete(ctx, int64(comment.ID))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
