package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories" // <-- 1. Tambahkan import ini
	"kasir-api/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// Catatan: Struct Produk di main.go ini sebenarnya sudah bisa dihapus 
// jika kamu sudah menggunakan models.Product dari package models agar tidak double.

func main() {
	// --- 1. Load Configuration ---
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// --- 2. Setup Database ---
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// --- 3. Inisialisasi Layer (Penting!) ---

	// A. Buat Repository dulu (Butuh db)
	productRepo := repositories.NewProductRepository(db)

	// B. Buat Service (Butuh repo, BUKAN db)
	productService := services.NewProductService(productRepo)

	// C. Buat Handler (Butuh service)
	productHandler := handlers.NewProductHandler(productService)

	// --- 4. Routes/Endpoints ---
	
	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API Running"})
	})

	// Setup routes menggunakan productHandler yang sudah dibuat di atas
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// --- 5. Run Server ---
	addr := ":" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("gagal running server:", err)
	}
}