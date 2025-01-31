package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
)

func ParseCfg(dsn string) *mysql.Config {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return cfg
}

func CreateDB(cfg mysql.Config) *sql.DB {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Hour * 3)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to database")

	return db
}

func utf8clean(s *string) {
	newString := ""
	for _, char := range *s {
		rune, _ := utf8.DecodeRuneInString(string(char))
		if rune != 65279 {
			newString += string(char)
		}
	}
	*s = newString
}

func ImportSQL(db *sql.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	temp := ""
	count := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		utf8clean(&line)
		line = strings.TrimSpace(line)
		temp += line

		if strings.HasSuffix(line, ";") {
			_, err := db.Exec(temp)
			if err != nil {
				log.Fatal(err)
			}
			temp = ""
		}

		count++
		fmt.Println(count)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("SQL file imported")
}
