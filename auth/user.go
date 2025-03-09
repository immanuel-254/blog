package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/immanuel-254/blog/auth/models"
	"github.com/immanuel-254/blog/database"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	// get data
	var data map[string]string
	GetData(data, w, r)

	// check if password match
	if data["password"] != data["confirm-password"] {
		http.Error(w, "invalid password or email", http.StatusBadRequest)
		return
	}

	// hash password
	hash, err := HashPassword(data["password"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create user
	user, err := queries.UserCreate(ctx, models.UserCreateParams{
		Email:     data["email"],
		Password:  hash,
		Isactive:  sql.NullBool{Bool: false, Valid: true},
		Isstaff:   sql.NullBool{Bool: false, Valid: true},
		Isadmin:   sql.NullBool{Bool: false, Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "create", user.ID, 0, w, r)

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(user.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(user.Email, "Activate Your Email", fmt.Sprintf("%s/activate/?token=%s", os.Getenv("DOMAIN"), one_time), EmailVerificationTemplate, w, r)

	resp := map[string]interface{}{"message": "signup successful"}
	SendData(resp, w, r)

}

func ActivateEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(database.DB)
	ctx := context.Background()

	// activate user
	user, err := queries.UserUpdateIsActive(ctx, models.UserUpdateIsActiveParams{
		ID:        int64(user_id),
		Isactive:  sql.NullBool{Bool: true, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "update", user.ID, int64(user_id), w, r)

	resp := map[string]interface{}{"message": "email has been verified"}
	SendData(resp, w, r)
}

func UserRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	user_id, err := strconv.ParseInt(queryParams.Get("user"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queries := models.New(database.DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	user, err := queries.UserRead(ctx, user_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "read", user.ID, authUser.ID, w, r)

	SendData(map[string]interface{}{"user": user}, w, r)
}

func UserList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queries := models.New(database.DB)

	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	users, err := queries.UserList(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "list", 0, authUser.ID, w, r)

	SendData(map[string]interface{}{"users": users}, w, r)
}

// Require auth
func ChangeEmailRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// get data
	var data map[string]string
	GetData(data, w, r)

	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	if data["email"] != authUser.Email {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(data["email"], "Change Your Email", fmt.Sprintf("%s/change-email/?token=%s", os.Getenv("DOMAIN"), one_time), ChangeEmailVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)

}

func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(database.DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	// get data
	var data map[string]string
	GetData(data, w, r)

	user, err := queries.UserUpdateEmail(ctx, models.UserUpdateEmailParams{
		ID:        int64(user_id),
		Email:     data["email"],
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "update", user.ID, int64(user_id), w, r)

	SendData(map[string]interface{}{"message": "email updated successfully"}, w, r)
}

func ChangePasswordRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(authUser.Email, "Change Your Password", fmt.Sprintf("%s/change-password/?token=%s", os.Getenv("DOMAIN"), one_time), ChangePasswordVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(database.DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	user, err := queries.UserLoginRead(ctx, authUser.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get data
	var data map[string]string
	GetData(data, w, r)

	check := CheckPasswordHash(data["old_password"], user.Password)

	if !check {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	if data["new_password"] != data["confirm_password"] {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	hash, err := HashPassword(data["new_password"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = queries.UserUpdatePassword(ctx, models.UserUpdatePasswordParams{
		ID:        int64(user_id),
		Password:  hash,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "update", user.ID, int64(user_id), w, r)

	SendData(map[string]interface{}{"message": "password updated successfully"}, w, r)
}

func ResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(authUser.Email, "Reset Your Password", fmt.Sprintf("%s/reset-password/?token=%s", os.Getenv("DOMAIN"), one_time), ResetPasswordVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(database.DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	user, err := queries.UserLoginRead(ctx, authUser.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get data
	var data map[string]string
	GetData(data, w, r)

	if data["new_password"] != data["confirm_password"] {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	hash, err := HashPassword(data["new_password"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = queries.UserUpdatePassword(ctx, models.UserUpdatePasswordParams{
		ID:        int64(user_id),
		Password:  hash,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "update", user.ID, int64(user_id), w, r)

	SendData(map[string]interface{}{"message": "password updated successfully"}, w, r)
}

func DeleteUserRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	// send email
	one_time, err := GenerateOneTimeToken(32, uint(authUser.ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendEmail(authUser.Email, "Delete User Account", fmt.Sprintf("%s/delete-user/?token=%s", os.Getenv("DOMAIN"), one_time), DeleteUserVerificationTemplate, w, r)

	SendData(map[string]interface{}{"message": "email sent"}, w, r)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	token := queryParams.Get("token")

	// verify token
	user_id, err := VerifyToken(token)

	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusBadRequest)
		return
	}

	queries := models.New(database.DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	if authUser.ID != int64(user_id) {
		http.Error(w, "Forbidden User", http.StatusForbidden)
		return
	}

	err = queries.UserDelete(ctx, int64(user_id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "delete", 0, int64(user_id), w, r)

	SendData(map[string]interface{}{"message": "user account deleted"}, w, r)
}

// require admin
func IsActiveChange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	user_id, err := strconv.ParseInt(queryParams.Get("user"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queries := models.New(database.DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	// get data
	var data map[string]string
	GetData(data, w, r)

	status, err := strconv.ParseBool(data["active"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := queries.UserUpdateIsActive(ctx, models.UserUpdateIsActiveParams{
		ID:        user_id,
		Isactive:  sql.NullBool{Bool: status, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "update", user.ID, authUser.ID, w, r)

	SendData(map[string]interface{}{"message": "user active status updated successfully"}, w, r)
}

func IsStaffChange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	user_id, err := strconv.ParseInt(queryParams.Get("user"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queries := models.New(database.DB)
	ctx := r.Context()

	auth := ctx.Value(current_user)

	if auth == nil {
		http.Error(w, "there is no current user", http.StatusInternalServerError)
		return
	}

	authUser := auth.(models.AuthUserReadRow)

	// get data
	var data map[string]string
	GetData(data, w, r)

	status, err := strconv.ParseBool(data["staff"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := queries.UserUpdateIsStaff(ctx, models.UserUpdateIsStaffParams{
		ID:        user_id,
		Isstaff:   sql.NullBool{Bool: status, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logging(queries, ctx, "user", "update", user.ID, authUser.ID, w, r)

	SendData(map[string]interface{}{"message": "user staff status updated successfully"}, w, r)
}
