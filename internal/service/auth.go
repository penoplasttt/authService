package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/penoplasttt/authService/internal/domain/models"
	"github.com/penoplasttt/authService/internal/lib/jwt"
	"github.com/penoplasttt/authService/internal/storage"
	"golang.org/x/crypto/bcrypt"
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

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
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", err)
		return "", fmt.Errorf("%s:%w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		a.log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)

	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user logged in succussfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", err)
		return "", fmt.Errorf("%s,%w", op, err)
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	uid, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", err)
			return 0, fmt.Errorf("%s:%w", op, ErrUserExists)
		}
		log.Error("failed to save user")
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user registered")

	return uid, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found ")
			return false, fmt.Errorf("%s:%w", op, ErrInvalidAppID)
		}
		return false, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("checking if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
