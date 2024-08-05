package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kinya-h/blog_engine/db"
	"github.com/kinya-h/blog_engine/util"
)

type createUserRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	PasswordHash string `json:"password_hash" binding:"required,min=6"`
	Email        string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	hashedPassword, err := util.HashPassword(req.PasswordHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusInternalServerError)
		return
	}

	arg := db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	result, err := server.db.CreateUser(context.Background(), arg)
	if err != nil {
		http.Error(w, fmt.Sprintf(" Error %s", err), http.StatusInternalServerError)
		return
	}

	userId, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("An Error Occured Creating user : %s", err.Error())
		http.Error(w, fmt.Sprintf(" Error %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User created successfully , id : %d ", userId)))
}

type loginUserRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	PasswordHash string `json:"password_hash" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	user, err := server.db.GetUser(context.Background(), req.Username)
	if err != nil {
		fmt.Printf("an error occured : %s", err.Error())
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInsufficientStorage)
		return

	}

	err = util.CheckPassword(req.PasswordHash, user.PasswordHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusUnauthorized)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		string(user.Role),
		time.Minute*15,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusUnauthorized)
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		string(user.Role),
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	session, err := server.db.CreateSession(context.Background(), db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    r.UserAgent(),
		ClientIp:     r.RemoteAddr,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	id, err := session.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	print(id)
	rsp := loginUserResponse{
		// SessionID:             sessionId.,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(rsp)

}
