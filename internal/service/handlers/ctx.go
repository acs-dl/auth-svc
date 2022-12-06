package handlers

import (
	"context"
	"github.com/mhrynenko/jwt_service/internal/data"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	usersCtxKey
	refreshTokensCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func UsersQ(r *http.Request) data.Users {
	return r.Context().Value(usersCtxKey).(data.Users).New()
}

func CtxUsersQ(entry data.Users) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, usersCtxKey, entry)
	}
}

func RefreshTokensQ(r *http.Request) data.RefreshTokens {
	return r.Context().Value(refreshTokensCtxKey).(data.RefreshTokens).New()
}

func CtxRefreshTokensQ(entry data.RefreshTokens) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, refreshTokensCtxKey, entry)
	}
}
