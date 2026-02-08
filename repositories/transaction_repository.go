package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// inisialisasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	// inisialisasi modeling transactionDetails -> nanti kita insert ke db
	details := make([]models.TransactionDetail, 0)
	// loop setiap item
	for _, item := range items {
		var productName string
		var productID, price, stock int
		// get product dapet pricing
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		// hitung current total = quantity * pricing
		// ditambahin ke dalam subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// kurangi jumlah stok
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		// item nya dimasukkin ke transactionDetails
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// insert transaction details
	// Di dalam fungsi CreateTransaction atau sejenisnya
for i := range details {
    // Pastikan ID transaksi induk sudah di-assign ke struct detail
    details[i].TransactionID = transactionID
    
    // Gunakan tx.Exec untuk menjaga atomicity (kesatuan transaksi)
    _, err = tx.Exec(
        "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
        details[i].TransactionID, 
        details[i].ProductID, 
        details[i].Quantity, 
        details[i].Subtotal,
    )
    if err != nil {
        // Jika satu gagal, tx.Rollback() biasanya dipanggil di level service/handler
        return nil, err 
    }
}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}

func (repo *TransactionRepository) GetSalesReport(startDate, endDate string) (models.SalesSummaryResponse, error) {
    var report models.SalesSummaryResponse

    summaryQuery := `
        SELECT 
            COALESCE(SUM(total_amount), 0), 
            COUNT(id) 
        FROM transactions 
        WHERE created_at::date BETWEEN $1 AND $2`

    err := repo.db.QueryRow(summaryQuery, startDate, endDate).Scan(
        &report.TotalRevenue,
        &report.TotalTransaksi,
    )
    if err != nil {
        return report, err
    }

    topProductQuery := `
        SELECT 
            p.name, 
            SUM(td.quantity) as total_qty
        FROM transaction_details td
        JOIN products p ON td.product_id = p.id
        JOIN transactions t ON td.transaction_id = t.id
        WHERE t.created_at::date BETWEEN $1 AND $2
        GROUP BY p.name
        ORDER BY total_qty DESC
        LIMIT 1`

    err = repo.db.QueryRow(topProductQuery, startDate, endDate).Scan(
        &report.ProdukTerlaris.Nama,
        &report.ProdukTerlaris.QtyTerjual,
    )

    if err != nil {
        report.ProdukTerlaris.Nama = "Belum ada transaksi"
        report.ProdukTerlaris.QtyTerjual = 0
    }

    return report, nil
}