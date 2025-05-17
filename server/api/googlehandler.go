package api

import (
	"context"
	"encoding/json"
	"log"
	"my-meal-planner/models"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// handleGoogleLogin initiates the Google OAuth flow
func (h *Handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Debug info
	log.Println("Starting OAuth flow...")

	// Check if OAuth config is properly initialized
	oauthConfig := h.store.GetOAuthConfig()
	if oauthConfig == nil {
		http.Error(w, "OAuth configuration is not available", http.StatusInternalServerError)
		log.Println("ERROR: OAuth config is nil")
		return
	}

	log.Printf("OAuth Config: ClientID=%s, RedirectURL=%s",
		oauthConfig.ClientID,
		oauthConfig.RedirectURL)

	// Generate a random state for CSRF protection
	state := uuid.New().String()

	// Store the state in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Minute.Seconds() * 5), // 5 minutes
		HttpOnly: true,
	})

	// Get the authorization URL from the OAuth config
	authURL := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline, // Get a refresh token as well
		oauth2.ApprovalForce,     // Force approval to get a refresh token
	)

	log.Println("Redirecting to:", authURL)

	// Redirect to Google's consent page
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// handleGoogleCallback handles the callback from Google OAuth
func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Get the state from the callback
	stateParam := r.URL.Query().Get("state")
	if stateParam == "" {
		http.Error(w, "State parameter missing", http.StatusBadRequest)
		return
	}

	// Get the state cookie
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "State cookie missing: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verify state parameter matches state cookie
	if stateCookie.Value != stateParam {
		http.Error(w, "State mismatch: cookie="+stateCookie.Value+" param="+stateParam, http.StatusBadRequest)
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Delete the cookie
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	// Exchange code for token using the OAuth config
	token, err := h.store.GetOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the token to get user info
	client := h.store.GetOAuthConfig().Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	// Parse the user info
	var userInfo struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or update the user
	user := &models.User{
		ID:    userInfo.Sub, // Use the Google ID as the user ID
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}

	if err := h.store.CreateOrUpdateUser(user); err != nil {
		http.Error(w, "Failed to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	jwtToken, err := h.store.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to the frontend with the token
	clientURL := os.Getenv("FRONTEND_URL")
	if clientURL == "" {
		clientURL = "http://localhost:5173" // fallback for local dev
	}

	http.Redirect(w, r, clientURL+"?token="+jwtToken, http.StatusTemporaryRedirect)
}
