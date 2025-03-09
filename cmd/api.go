package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/immanuel-254/blog/auth"
)

var (
	Login = auth.View{
		Route:   "/login",
		Handler: http.HandlerFunc(auth.Login),
	}

	Logout = auth.View{
		Route:   "/logout",
		Handler: http.HandlerFunc(auth.Logout),
	}

	Signup = auth.View{
		Route:   "/signup",
		Handler: http.HandlerFunc(auth.Signup),
	}

	ActivateEmail = auth.View{
		Route:   "/activate",
		Handler: http.HandlerFunc(auth.ActivateEmail),
	}

	UserRead = auth.View{
		Route:       "/read",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAdmin},
		Handler:     http.HandlerFunc(auth.UserRead),
	}

	UserList = auth.View{
		Route:       "/list",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAdmin},
		Handler:     http.HandlerFunc(auth.UserList),
	}

	ChangeEmailRequest = auth.View{
		Route:       "/change-email-request",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.ChangeEmailRequest),
	}

	ChangeEmail = auth.View{
		Route:       "/change-email",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.ChangeEmail),
	}

	ChangePasswordRequest = auth.View{
		Route:       "/change-password-request",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.ChangePasswordRequest),
	}

	ChangePassword = auth.View{
		Route:       "/change-password",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.ChangePassword),
	}

	ResetPasswordRequest = auth.View{
		Route:       "/reset-password-request",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.ResetPasswordRequest),
	}

	ResetPassword = auth.View{
		Route:       "/reset-password",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.ResetPassword),
	}

	DeleteUserRequest = auth.View{
		Route:       "/delete-user-request",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.DeleteUserRequest),
	}

	DeleteUser = auth.View{
		Route:       "/delete-user",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAuth},
		Handler:     http.HandlerFunc(auth.DeleteUser),
	}

	IsActiveChange = auth.View{
		Route:       "/isactive",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAdmin},
		Handler:     http.HandlerFunc(auth.IsActiveChange),
	}

	IsStaffChange = auth.View{
		Route:       "/isstaff",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAdmin},
		Handler:     http.HandlerFunc(auth.IsStaffChange),
	}

	SessionList = auth.View{
		Route:       "/session/list",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAdmin},
		Handler:     http.HandlerFunc(auth.SessionList),
	}

	LogList = auth.View{
		Route:       "/log/list",
		Middlewares: []func(http.Handler) http.Handler{auth.RequireAdmin},
		Handler:     http.HandlerFunc(auth.LogList),
	}
)

func Api() {
	mux := http.NewServeMux()

	allviews := []auth.View{
		Login,
		Logout,
		Signup,
		ActivateEmail,
		UserRead,
		UserList,
		ChangeEmailRequest,
		ChangeEmail,
		ResetPasswordRequest,
		ResetPassword,
		DeleteUserRequest,
		DeleteUser,
		IsActiveChange,
		IsStaffChange,

		SessionList,

		LogList,
	}

	auth.Routes(mux, allviews)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", os.Getenv("PORT")), // Custom port
		//Handler:      internal.LoggingMiddleware(internal.Cors(internal.New(internal.ConfigDefault)(mux))), // Attach the mux as the handler
		Handler:      auth.LoggingMiddleware(mux),
		ReadTimeout:  10 * time.Second, // Set read timeout
		WriteTimeout: 10 * time.Second, // Set write timeout
		IdleTimeout:  30 * time.Second, // Set idle timeout
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting server:", err)
	}
}
