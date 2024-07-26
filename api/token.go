package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req renewAccessTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusUnauthorized)
		return
	}

	session, err := server.db.GetSession(context.Background(), refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusUnauthorized)

		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusUnauthorized)
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusUnauthorized)
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusUnauthorized)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		refreshPayload.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusUnauthorized)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rsp)
}
