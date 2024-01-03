package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/penoplasttt/authService/internal/domain/models"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	panic("not implemented")
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error){
	panic("not implemented")
}

func (a *Auth)IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}
