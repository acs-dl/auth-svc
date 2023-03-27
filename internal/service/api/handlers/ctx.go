package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/config"
	"gitlab.com/distributed_lab/acs/auth/internal/data"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	usersCtxKey
	modulesCtxKey
	permissionsCtxKey
	refreshTokensCtxKey
	jwtParamsCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func JwtParams(r *http.Request) config.JwtCfg {
	return r.Context().Value(jwtParamsCtxKey).(config.JwtCfg)
}

func CtxJwtParams(jwtParams config.JwtCfg) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, jwtParamsCtxKey, jwtParams)
	}
}

func UsersQ(r *http.Request) data.Users {
	return r.Context().Value(usersCtxKey).(data.Users).New()
}

func CtxUsersQ(entry data.Users) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, usersCtxKey, entry)
	}
}

func ModulesQ(r *http.Request) data.Modules {
	return r.Context().Value(modulesCtxKey).(data.Modules).New()
}

func CtxModulesQ(entry data.Modules) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, modulesCtxKey, entry)
	}
}

func PermissionsQ(r *http.Request) data.Permissions {
	return r.Context().Value(permissionsCtxKey).(data.Permissions).New()
}

func CtxPermissionsQ(entry data.Permissions) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, permissionsCtxKey, entry)
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
