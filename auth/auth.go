package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/immanuel-254/blog/auth/models"
)

func AuthLogin(queries *models.Queries, ctx context.Context, data map[string]string) (string, int, error) {
	user, err := queries.UserLoginRead(ctx, data["email"])

	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("%s, (%s)", data["email"], err.Error())
	}

	check := CheckPasswordHash(data["password"], user.Password)

	if !check {
		return "", http.StatusBadRequest, fmt.Errorf("invalid credentials")
	}

	// create key
	key := base64.StdEncoding.EncodeToString(GenerateAESKey())

	// create session

	session, err := queries.SessionCreate(ctx, models.SessionCreateParams{
		Key:       key,
		UserID:    user.ID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	err = queries.LogCreate(ctx, models.LogCreateParams{
		DbTable:   "session",
		Action:    "create",
		ObjectID:  session.ID,
		UserID:    session.UserID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return key, http.StatusOK, nil
}
