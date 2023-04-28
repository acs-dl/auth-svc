package helpers

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
)

func SetTokensCookies(w http.ResponseWriter, access, refresh string) {
	refreshCookie := &http.Cookie{
		Name:     data.RefreshCookie,
		Value:    refresh,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, refreshCookie)

	accessCookie := &http.Cookie{
		Name:     data.AccessCookie,
		Value:    access,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, accessCookie)
}

func ClearTokensCookies(w http.ResponseWriter) {
	refreshCookie := &http.Cookie{
		Name:     data.RefreshCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, refreshCookie)

	accessCookie := &http.Cookie{
		Name:     data.AccessCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, accessCookie)
}
