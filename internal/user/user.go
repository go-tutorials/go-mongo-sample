package user

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	v "github.com/core-go/core/v10"

	"go-service/internal/user/handler"
	"go-service/internal/user/repository/adapter"
	"go-service/internal/user/repository/query"
	"go-service/internal/user/service"
)

type UserTransport interface {
	All(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(db *mongo.Database, logError func(context.Context, string, ...map[string]interface{})) (UserTransport, error) {
	validator, err := v.NewValidator()
	if err != nil {
		return nil, err
	}

	userRepository := adapter.NewUserAdapter(db, query.BuildQuery)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, logError, validator.Validate)
	return userHandler, nil
}
