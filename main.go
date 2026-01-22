package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Struktur Data Produk
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// Data dummy disimpan dalam slice (In-Memory)
var produkbaru = []Produk{
	{ID: 1, Nama: "Indomie Goreng", Harga: 3500, Stok: 50},
	{ID: 2, Nama: "Vit (air mineral lokal)", Harga: 3000, Stok: 100},
	{ID: 3, Nama: "Kecap (buat masak)", Harga: 12000, Stok: 20},
}

func main() {
	// 1. Health Check Endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API Running"})
	})

	// 2. Endpoint Get All & Create Produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			json.NewEncoder(w).Encode(produkbaru)
		} else if r.Method == "POST" {
			var p Produk
			err := json.NewDecoder(r.Body).Decode(&p)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}
			p.ID = len(produkbaru) + 1
			produkbaru = append(produkbaru, p)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(p)
		}
	})

	// 3. Endpoint Get, Update, & Delete by ID
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
			return
		}

		if r.Method == "GET" {
			for _, p := range produkbaru {
				if p.ID == id {
					json.NewEncoder(w).Encode(p)
					return
				}
			}
		} else if r.Method == "PUT" {
			var updateP Produk
			json.NewDecoder(r.Body).Decode(&updateP)
			for i, p := range produkbaru {
				if p.ID == id {
					updateP.ID = id
					produkbaru[i] = updateP
					json.NewEncoder(w).Encode(updateP)
					return
				}
			}
		} else if r.Method == "DELETE" {
			for i, p := range produkbaru {
				if p.ID == id {
					produkbaru = append(produkbaru[:i], produkbaru[i+1:]...)
					json.NewEncoder(w).Encode(map[string]string{"message": "sukses delete"})
					return
				}
			}
		}

		http.Error(w, "Produk belum ada", http.StatusNotFound)
	})

	fmt.Println("Server running di localhost:8080")
	http.ListenAndServe(":8080", nil)
}