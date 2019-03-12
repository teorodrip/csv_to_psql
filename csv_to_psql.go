package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sqweek/dialog"
	"log"
	"os"
	"regexp"
)

const N_COLS = 24

const DB_USER = "uncharblog"
const DB_NAME = "pam_test"
const DB_PASS = "K5N3gwww5U8Yxfcv"
const DB_HOST = "192.168.27.122"
const CONN_STR = "dbname=" + DB_NAME +
	" user=" + DB_USER +
	" password=" + DB_PASS +
	" host=" + DB_HOST

type pgDB struct {
	Db           *sql.DB
	SqlInsertRow *sql.Stmt
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
		if err = p.PrepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func (p *pgDB) PrepareSqlStatements() error {
	var err error

	if p.SqlInsertRow, err = p.Db.Prepare("INSERT INTO uncharblog.grid_table (op, vis, news_heat, quantity, secur, side, status, lmt_pr, tif, fill_qty, avg_pr, filled, working_qty, idle, data_export_restricted, last_update, create_time, vwap, data_export_restricted_2, col_last, bid, ask, volume, d_adv) SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24 WHERE NOT EXISTS (SELECT 1 FROM uncharblog.grid_table WHERE op IS NOT DISTINCT FROM $1 AND vis IS NOT DISTINCT FROM $2 AND news_heat IS NOT DISTINCT FROM $3 AND quantity IS NOT DISTINCT FROM $4 AND secur IS NOT DISTINCT FROM $5 AND side IS NOT DISTINCT FROM $6 AND status IS NOT DISTINCT FROM $7 AND lmt_pr IS NOT DISTINCT FROM $8 AND tif IS NOT DISTINCT FROM $9 AND fill_qty IS NOT DISTINCT FROM $10 AND avg_pr IS NOT DISTINCT FROM $11 AND filled IS NOT DISTINCT FROM $12 AND working_qty IS NOT DISTINCT FROM $13 AND idle IS NOT DISTINCT FROM $14 AND data_export_restricted IS NOT DISTINCT FROM $15 AND last_update IS NOT DISTINCT FROM $16 AND create_time IS NOT DISTINCT FROM $17 AND vwap IS NOT DISTINCT FROM $18 AND data_export_restricted_2 IS NOT DISTINCT FROM $19 AND col_last IS NOT DISTINCT FROM $20 AND bid IS NOT DISTINCT FROM $21 AND ask IS NOT DISTINCT FROM $22 AND volume IS NOT DISTINCT FROM $23 AND d_adv IS NOT DISTINCT FROM $24);"); err != nil {
		return err
	}
	return nil
}

func NewSqlNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func main() {
	date_format := regexp.MustCompile("[0-9]*/[0-9]*/[0-9]*")

	path, err := dialog.File().Filter("CSV Files", "csv").Load()
	if err != nil {
		dialog.Message("%s (%s)", "Can not get file.", err.Error()).Title("Error!").Info()
		fmt.Fprintf(os.Stderr, "Error: can not get file, (%s)\n", err.Error())
		return
	}
	db, err := ConnectDB(CONN_STR)
	if err != nil {
		dialog.Message("%s (%s)", "can not connect to data base", err.Error()).Title("Error!").Info()
		fmt.Fprintf(os.Stderr, "Error: can not connect to data base, (%s)\n", err.Error())
		return
	}
	defer db.Db.Close()
	f, err := os.Open(path)
	if err != nil {
		dialog.Message("%s (%s)", "can not open file", err.Error()).Title("Error!").Info()
		fmt.Fprintf(os.Stderr, "Error: can not open file, (%s)\n", err.Error())
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		dialog.Message("%s (%s)", "can not read from csv", err.Error()).Title("Error!").Info()
		fmt.Fprintf(os.Stderr, "Error: can not read from csv, (%s)\n", err.Error())
		return
	}
	int_row := make([]interface{}, N_COLS)
	data = data[1:]
	for _, row := range data {
		if len(row) != N_COLS {
			dialog.Message("%s", "Error in the format of the table (must have 24 columns).").Title("Error!").Info()
			fmt.Fprintf(os.Stderr, "Error: in the format of the table\n")
			return
		}
		for j, cell := range row {
			if j == 15 || j == 16 {
				cell = date_format.FindString(cell)
			}
			int_row[j] = NewSqlNullString(cell)
		}
		_, err := db.SqlInsertRow.Exec(int_row...)
		if err != nil {
			dialog.Message("%s (%s)", "Error executing the query", err.Error()).Title("Error!").Info()
			fmt.Fprintf(os.Stderr, "Error: executing the query, (%s)\n", err.Error())
			return
		}
	}
	dialog.Message("%s", "Table successfully uploaded.").Title("OK").Info()
	fmt.Printf("Table uploaded.\n")
}
