package main

import (
	"fmt"
	"os"
	"time"

	"aldrico.com/go-sql-import/cmd/database"
)

func main() {
	start := time.Now()

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Required arguments: <dsn> <sql_file_path>")
		os.Exit(1)
	}

	cfg := database.ParseCfg(args[0])
	filePath := args[1]

	db := database.CreateDB(*cfg)
	defer db.Close()

	database.ImportSQL(db, filePath)

	elapsed := time.Since(start)
	fmt.Printf("Imported in %s\n", elapsed)
}
