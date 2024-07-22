package handler

import (
	"context"
	"encoding/json"
	"gocommerce/model"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek apakah request memiliki header Authorization
		tokenString := r.Header.Get("Authorization")
		claims, IsValid := VerifyJWT(tokenString)
		if !IsValid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// adding username to the context if needed
		ctx := context.WithValue(r.Context(), ContextKeyUsername, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Contoh validasi sederhana: pastikan username dan password tidak kosong
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	err = model.RegisterUser(user)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "username": user.Username, "token": token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Contoh validasi sederhana: pastikan username dan password tidak kosong
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Cek apakah user dengan username tersebut ada di database
	hashedPassword, err := model.GetUserPassword(user.Username)
	if err != nil {
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Bandingkan password yang diinputkan dengan password di database
	if !model.CheckPasswordHash(user.Password, hashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "username": user.Username, "token": token})
}

func Profile(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(ContextKeyUsername).(string)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username})
}
