package models

// Struktur Data Produk
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock  int    `json:"stock"`
}

// Data dummy disimpan dalam slice (In-Memory)
// var produkbaru = []Produk{
// 	{ID: 1, Nama: "Indomie Goreng", Harga: 3500, Stok: 50},
// 	{ID: 2, Nama: "Vit (air mineral lokal)", Harga: 3000, Stok: 100},
// 	{ID: 3, Nama: "Kecap (buat masak)", Harga: 12000, Stok: 20},
// }
