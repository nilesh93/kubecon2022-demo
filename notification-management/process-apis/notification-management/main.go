package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nuwand/kubecon2022-demo/notification-api/helpers"
	"github.com/nuwand/kubecon2022-demo/notification-api/services"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type NotificationRequest struct {
	UserEmail string `json:"userEmail,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
	SendEmail bool   `json:"sendEmail"`
}

func main() {

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Route("/healthz", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			helpers.RespondwithJSON(w, 200, "healthy")
		})
	})

	r.Route("/api/v1/", func(r chi.Router) {
		// @Summary Create Notification
		// @Tags Notification
		// @Accept json
		// @Produce json
		// @Param data body NotificationRequest	true	"data"
		// @Success 200 {object} services.NotificationCreateResponse	"Okay"
		// @Failure 400 {string} string
		// @Failure 500 {string} string
		// @Router /notification [post]
		r.Post("/notification", func(w http.ResponseWriter, r *http.Request) {

			body := NotificationRequest{}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.RespondWithError(w, 400, err.Error())
				return
			}

			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			obj, err := NotificationCreate(ctx, body)
			if err != nil {
				helpers.RespondWithError(w, 500, err.Error())
				return
			}
			helpers.RespondwithJSON(w, 200, obj)
		})
	})

	logrus.Info("http server started")
	http.ListenAndServe(":4000", r)
}

func NotificationCreate(ctx context.Context, req NotificationRequest) (*services.NotificationCreateResponse, error) {

	if v, err := helpers.IsValid(req); !v {
		return nil, err
	}

	res, err := services.NECreateNotification(services.NotificationCreateRequest{
		UserEmail: req.UserEmail,
		Message:   req.Message,
		Type:      req.Type,
	})
	if err != nil {
		return nil, err
	}

	if req.SendEmail {
		go services.SendEmail(ctx, req.UserEmail)
	}

	return res, nil
}
