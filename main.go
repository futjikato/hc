package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

func main() {
	no_load := flag.Bool("L", false, "Do not collect load information.")
	no_io := flag.Bool("S", false, "Do not collect disk IO information.")
	no_net := flag.Bool("N", false, "Do not collect network information.")

	db_path := flag.String("db", "./", "Path to directory in witch to save config.db sqlite3 file")

	flag.Parse()

	c := &StatConfig{}
	c.net = !*no_net
	c.load = !*no_load
	c.io = !*no_io
	c.lastIoWriteCount = make(map[string]int64)
	c.lastIoReadCount = make(map[string]int64)
	c.lastNetSentBytes = make(map[string]int64)
	c.lastNetRecvBytes = make(map[string]int64)

	db, db_err := sql.Open("sqlite3", filepath.Join(*db_path, "config.db"))
	if db_err != nil {
		panic(db_err)
	}
	defer db.Close()
	initDatabase(db)
	initLastValues(db, c)

	s := getStats(c)

	saveConfig(db, c)
	saveStats(db, s)
	fmt.Print(*s)
}

func initDatabase(db *sql.DB) {
	sql_stmt := `
		CREATE TABLE IF NOT EXISTS config_io(name TEXT PRIMARY KEY, read_count integer, write_count integer);
		CREATE TABLE IF NOT EXISTS config_net(name TEXT PRIMARY KEY, bytes_sent integer, bytes_recv integer);
		CREATE TABLE IF NOT EXISTS stats(id integer PRIMARY KEY, payload TEXT);
		DELETE FROM stats;
	`
	_, err := db.Exec(sql_stmt)
	if err != nil {
		panic(err)
	}
}

func initLastValues(db *sql.DB, c *StatConfig) {
	initLastIoValues(db, c)
	initLastNetValues(db, c)
}

func initLastIoValues(db *sql.DB, c *StatConfig) {
	r, err := db.Query("SELECT name, read_count, write_count FROM config_io")
	if err != nil {
		panic(err)
	}

	defer r.Close()
	for r.Next() {
		var rc int64
		var wc int64
		var name string
		err = r.Scan(&name, &rc, &wc)
		if err != nil {
			// todo maybe don´t panic here and just ignore row?
			panic(err)
		}

		c.lastIoReadCount[name] = rc
		c.lastIoWriteCount[name] = wc
	}
}

func initLastNetValues(db *sql.DB, c *StatConfig) {
	r, err := db.Query("SELECT name, bytes_sent, bytes_recv FROM config_net")
	if err != nil {
		panic(err)
	}

	defer r.Close()
	for r.Next() {
		var bs int64
		var br int64
		var name string
		err = r.Scan(&name, &bs, &br)
		if err != nil {
			// todo maybe don´t panic here and just ignore row?
			panic(err)
		}

		c.lastNetSentBytes[name] = bs
		c.lastNetRecvBytes[name] = br
	}
}

func saveConfig(db *sql.DB, c *StatConfig) {
	ioStmt, err := db.Prepare("INSERT OR REPLACE INTO config_io (name, read_count, write_count) VALUES (?,?,?)")
	if err != nil {
		panic(err)
	}
	defer ioStmt.Close()

	for ioName, rv := range c.lastIoReadCount {
		ioStmt.Exec(ioName, rv, c.lastIoWriteCount[ioName])
	}

	netStmt, netErr := db.Prepare("INSERT OR REPLACE INTO config_net (name, bytes_sent, bytes_recv) VALUES (?,?,?)")
	if netErr != nil {
		panic(netErr)
	}
	defer netStmt.Close()

	for netName, sv := range c.lastNetSentBytes {
		netStmt.Exec(netName, sv, c.lastNetRecvBytes[netName])
	}
}

func saveStats(db *sql.DB, s *StatSet) {
	stmt, err := db.Prepare("INSERT INTO stats (id, payload) VALUES (?,?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	b, jsonErr := json.Marshal(*s)
	if jsonErr != nil {
		panic(jsonErr)
	}
	stmt.Exec(nil, string(b))
}
