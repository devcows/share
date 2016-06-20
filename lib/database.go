package lib

import (
	"database/sql"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/tylerb/graceful"
)

type Server struct {
	ID               int       `json:"id"`
	Path             string    `json:"path"`
	Port             int       `json:"port"`
	InsertedDatetime time.Time `json:"inserted_datetime"`
	ListIps          []string  `json:"list_ips"`
	Srv              *graceful.Server
}

var settings SettingsShare

func OpenDatabase() (*sql.DB, error) {
	destDb, err := sql.Open("sqlite3_share", settings.Daemon.DatabaseFilePath)
	if err != nil {
		return nil, err
	}
	destDb.Ping()

	return destDb, nil
}

func InitDB(settings_params SettingsShare) error {
	settings = settings_params

	sqlite3conn := []*sqlite3.SQLiteConn{}
	sql.Register("sqlite3_share",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				return nil
			},
		})

	err2 := CreateTable()
	if err2 != nil {
		return err2
	}

	return nil
}

func CreateTable() error {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS servers(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Path TEXT,
		Port int,
		ListIps TEXT,
		InsertedDatetime datetime
	);
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}

	_, err2 := destDb.Exec(sql_table)
	if err2 != nil {
		return err2
	}

	defer destDb.Close()
	return nil
}

func StoreServer(server Server) (int64, error) {
	sql_add := `
	INSERT OR REPLACE INTO servers(
		Path,
		Port,
		ListIps,
		InsertedDatetime
	) values(?, ?, ?, CURRENT_TIMESTAMP)
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return -1, err
	}

	stmt, err2 := destDb.Prepare(sql_add)
	if err2 != nil {
		return -1, err2
	}
	defer stmt.Close()

	res, err3 := stmt.Exec(server.Path, server.Port, strings.Join(server.ListIps, "||"))
	if err3 != nil {
		return -1, err3
	}

	id, err4 := res.LastInsertId()
	if err4 != nil {
		return -1, err4
	}

	defer destDb.Close()

	return id, nil
}

func RemoveServer(id int64) error {
	sql_remove := `
	delete from servers where
	Id = ?
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}

	stmt, err := destDb.Prepare(sql_remove)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(id)
	if err2 != nil {
		return err2
	}

	defer destDb.Close()

	return nil
}

func ListServers() ([]Server, error) {
	sql_select := `
	select Id, Path, Port, ListIps from servers
	order by InsertedDatetime DESC
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return []Server{}, err
	}

	rows, err2 := destDb.Query(sql_select)
	if err2 != nil {
		return []Server{}, err2
	}
	defer rows.Close()

	var results []Server
	for rows.Next() {
		item := Server{}

		listIps := ""
		err3 := rows.Scan(&item.ID, &item.Path, &item.Port, &listIps)
		if err3 == nil {
			item.ListIps = strings.Split(listIps, "||")
			results = append(results, item)
		}
	}

	defer destDb.Close()

	return results, nil
}
