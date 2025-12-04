package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Render cloud database URL
	dsn := "postgresql://finance_service_14r5_user:6QJcvvToLvJ5EUGBGlrYWq15JIV4PrOC@dpg-d4jdlj7gi27c739kmt3g-a.singapore-postgres.render.com/finance_service_14r5?sslmode=require"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to cloud database: %v", err)
	}

	fmt.Println("\n✅ Connected to Render Cloud Database!")
	fmt.Println("═══════════════════════════════════════════════════")

	// Get all tables
	var tables []string
	db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`).Scan(&tables)

	fmt.Printf("\n📊 Tables Found: %d\n\n", len(tables))
	for i, table := range tables {
		// Count rows in each table
		var count int64
		db.Table(table).Count(&count)
		fmt.Printf("%d. %s (%d rows)\n", i+1, table, count)
	}

	// Show table sizes
	fmt.Println("\n📈 Table Sizes:")
	var sizes []map[string]interface{}
	db.Raw(`
		SELECT 
			schemaname,
			tablename,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
		FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
	`).Scan(&sizes)

	for _, size := range sizes {
		fmt.Printf("  %s: %v\n", size["tablename"], size["size"])
	}

	// Show column info for each table
	fmt.Println("\n🔍 Table Structures:")
	for _, table := range tables {
		fmt.Printf("\n📋 %s:\n", table)

		var columns []map[string]interface{}
		db.Raw(`
			SELECT 
				column_name,
				data_type,
				is_nullable
			FROM information_schema.columns
			WHERE table_name = ?
			ORDER BY ordinal_position
		`, table).Scan(&columns)

		for _, col := range columns {
			nullable := "NOT NULL"
			if col["is_nullable"] == "YES" {
				nullable = "NULLABLE"
			}
			fmt.Printf("   ├─ %s: %s (%s)\n", col["column_name"], col["data_type"], nullable)
		}
	}

	fmt.Println("\n✅ Database inspection complete!")
}
