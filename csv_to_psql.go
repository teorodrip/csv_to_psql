package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

const DB_USER = "uncharblog"
const DB_NAME = "pam_test"
const DB_PASS = "K5N3gwww5U8Yxfcv"
const DB_HOST = "192.168.27.122"
const CONN_STR = "dbname=" + DB_NAME +
	" user=" + DB_USER +
	" password=" + DB_PASS +
	" host=" + DB_HOST

type pgDB struct {
	Db *sql.DB
}

func ConnectDB(ConnStr string) (*pgDB, error) {
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to %s\n", ConnStr)
		log.Print(err)
		return nil, err
	} else {
		p := &pgDB{Db: db}
		if err = p.Db.Ping(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func main() {
	var query strings.Builder

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: wrong number of arguments, expected 1 (path to csv).\n")
		return
	}
	if len(os.Args[1]) < 5 || os.Args[1][len(os.Args[1])-4:] != ".csv" {
		fmt.Fprintf(os.Stderr, "Error: wrong format of file, expected (.csv).\n")
		return
	}
	db, err := ConnectDB(CONN_STR)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: can not connect to data base, (%s)\n", err.Error())
		return
	}
	defer db.Db.Close()
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: can not open file, (%s)\n", err.Error())
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: can not read from csv, (%s)\n", err.Error())
		return
	}
	query.Write([]byte("INSERT INTO uncharblog.csv_test_1 (col_1, col_2) VALUES "))
	for i, row := range data {
		for j, cell := range row {
			if j == 0 {
				query.Write([]byte("("))
			}
			query.Write([]byte("'"))
			query.Write([]byte(cell))
			if j < len(row)-1 {
				query.Write([]byte("',"))
			} else {
				query.Write([]byte("')"))
			}
		}
		if i < len(data)-1 {
			query.Write([]byte(","))
		}
	}
	_, err = db.Db.Exec(query.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: can not perform the query, (%s):\n%s\n", err.Error(), query.String())
		return
	}
	fmt.Printf("Table uploaded.\n")
}
