package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/stock_simulation?parseTime=true")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Insert AMZN historical data
	err = insertAMZNData(db)
	if err != nil {
		log.Fatal("Failed to insert AMZN data:", err)
	}

	fmt.Println("✅ AMZN historical data inserted successfully!")
}

func insertAMZNData(db *sql.DB) error {
	// Sample data for AMZN (last 30 days)
	amznData := [][]interface{}{
		{"AMZN", time.Now().AddDate(0, 0, -30), 3180.50, 3225.30, 3165.80, 3200.20, 32345600},
		{"AMZN", time.Now().AddDate(0, 0, -29), 3200.20, 3238.40, 3185.90, 3215.80, 28234500},
		{"AMZN", time.Now().AddDate(0, 0, -28), 3215.80, 3251.20, 3198.30, 3235.10, 31765400},
		{"AMZN", time.Now().AddDate(0, 0, -27), 3235.10, 3268.60, 3212.40, 3248.90, 29187300},
		{"AMZN", time.Now().AddDate(0, 0, -26), 3248.90, 3275.30, 3225.20, 3261.70, 33654200},
		{"AMZN", time.Now().AddDate(0, 0, -25), 3261.70, 3289.90, 3245.10, 3274.40, 30432100},
		{"AMZN", time.Now().AddDate(0, 0, -24), 3274.40, 3298.80, 3256.30, 3285.20, 35976500},
		{"AMZN", time.Now().AddDate(0, 0, -23), 3285.20, 3315.50, 3268.10, 3302.60, 27234700},
		{"AMZN", time.Now().AddDate(0, 0, -22), 3302.60, 3328.30, 3285.40, 3318.90, 34543200},
		{"AMZN", time.Now().AddDate(0, 0, -21), 3318.90, 3345.70, 3302.80, 3332.50, 29876300},
		{"AMZN", time.Now().AddDate(0, 0, -20), 3332.50, 3356.20, 3315.20, 3348.80, 32654300},
		{"AMZN", time.Now().AddDate(0, 0, -19), 3348.80, 3372.50, 3331.90, 3365.40, 28187400},
		{"AMZN", time.Now().AddDate(0, 0, -18), 3365.40, 3389.10, 3348.30, 3378.70, 31543600},
		{"AMZN", time.Now().AddDate(0, 0, -17), 3378.70, 3402.20, 3361.80, 3391.30, 29876500},
		{"AMZN", time.Now().AddDate(0, 0, -16), 3391.30, 3418.40, 3374.50, 3405.80, 33765400},
		{"AMZN", time.Now().AddDate(0, 0, -15), 3405.80, 3431.90, 3388.60, 3419.50, 30321700},
		{"AMZN", time.Now().AddDate(0, 0, -14), 3419.50, 3445.30, 3402.20, 3433.40, 32654800},
		{"AMZN", time.Now().AddDate(0, 0, -13), 3433.40, 3459.60, 3416.10, 3447.20, 28456300},
		{"AMZN", time.Now().AddDate(0, 0, -12), 3447.20, 3473.80, 3430.30, 3461.50, 31876500},
		{"AMZN", time.Now().AddDate(0, 0, -11), 3461.50, 3488.40, 3444.60, 3475.30, 30234800},
		{"AMZN", time.Now().AddDate(0, 0, -10), 3475.30, 3502.70, 3458.40, 3489.90, 33543600},
		{"AMZN", time.Now().AddDate(0, 0, -9), 3489.90, 3516.50, 3472.80, 3504.60, 29567200},
		{"AMZN", time.Now().AddDate(0, 0, -8), 3504.60, 3531.30, 3487.40, 3519.80, 32234500},
		{"AMZN", time.Now().AddDate(0, 0, -7), 3519.80, 3546.90, 3502.90, 3535.40, 28876300},
		{"AMZN", time.Now().AddDate(0, 0, -6), 3535.40, 3562.20, 3518.50, 3551.10, 31654700},
		{"AMZN", time.Now().AddDate(0, 0, -5), 3551.10, 3578.60, 3534.30, 3567.80, 30432800},
		{"AMZN", time.Now().AddDate(0, 0, -4), 3567.80, 3595.40, 3550.90, 3584.50, 33765400},
		{"AMZN", time.Now().AddDate(0, 0, -3), 3584.50, 3612.20, 3567.60, 3601.30, 29321600},
		{"AMZN", time.Now().AddDate(0, 0, -2), 3601.30, 3629.80, 3584.40, 3618.90, 32456700},
		{"AMZN", time.Now().AddDate(0, 0, -1), 3618.90, 3647.50, 3601.80, 3636.60, 28456700},
	}

	// Insert AMZN data
	query := `INSERT IGNORE INTO historical_prices (symbol, date, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	fmt.Println("Inserting AMZN historical data...")
	for _, data := range amznData {
		_, err := db.Exec(query, data...)
		if err != nil {
			return fmt.Errorf("failed to insert AMZN data: %v", err)
		}
	}

	return nil
} 