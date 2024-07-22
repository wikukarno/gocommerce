package main

import (
    "gocommerce/database"
    "gocommerce/handler"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    database.InitDB()
    defer database.DB.Close()

    r := mux.NewRouter()

    r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
    r.Handle("/profile", handler.AuthMiddleware(http.HandlerFunc(handler.Profile)))

    // Tambahkan logging untuk memastikan route terdaftar
    r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        tmpl, _ := route.GetPathTemplate()
        log.Printf("Registered Route: %s\n", tmpl)
        return nil
    })

    // Gunakan 'r' sebagai handler dalam ListenAndServe
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("ListenAndServe error: ", err)
    }
}
