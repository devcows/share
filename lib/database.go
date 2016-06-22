package lib

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/cznic/ql/driver"
	"github.com/tylerb/graceful"
)

type Server struct {
	UUID      string    `json:"uuid"`
	Path      string    `json:"path"`
	Port      int       `json:"port"`
	CreatedAt time.Time `json:"created_at"`
	ListIps   []string  `json:"list_ips"`
	Srv       *graceful.Server
}

var settings SettingsShare

func OpenDatabase() (*sql.DB, error) {
	destDb, err := sql.Open("ql", settings.Daemon.DatabaseFilePath)
	if err != nil {
		return nil, err
	}
	destDb.Ping()

	return destDb, nil
}

func InitDB(settings_params SettingsShare) error {
	settings = settings_params
	err2 := CreateTable()
	if err2 != nil {
		return err2
	}

	return nil
}

func CreateTable() error {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS Servers(
		UUID string,
		Path string,
		Port int,
		ListIps string,
		CreatedAt time
	);
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}

	defer destDb.Close()

	tx, err1 := destDb.Begin()
	if err1 != nil {
		return err1
	}

	_, err2 := tx.Exec(sql_table)
	if err2 != nil {
		return err2
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func StoreServer(server Server) error {
	sqlAdd := `
	INSERT INTO Servers(
		UUID,
		Path,
		Port,
		ListIps,
		CreatedAt
	) values($1, $2, $3, $4, $5)
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}
	defer destDb.Close()

	tx, err1 := destDb.Begin()
	if err1 != nil {
		return err1
	}

	stmt, err2 := tx.Prepare(sqlAdd)
	if err2 != nil {
		return err2
	}
	defer stmt.Close()

	_, err3 := stmt.Exec(server.UUID, server.Path, server.Port, strings.Join(server.ListIps, "||"), server.CreatedAt)
	if err3 != nil {
		return err3
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func RemoveServer(uuid string) error {
	sqlRemove := `
	delete from Servers where
	UUID == $1
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}
	defer destDb.Close()

	tx, err1 := destDb.Begin()
	if err1 != nil {
		return err1
	}

	stmt, err2 := tx.Prepare(sqlRemove)
	if err2 != nil {
		return err2
	}
	defer stmt.Close()

	if _, err = stmt.Exec(uuid); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ListServers() ([]Server, error) {
	var results []Server

	sql_select := `
	select UUID, Path, Port, ListIps, CreatedAt from Servers order by CreatedAt
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return results, err
	}
	defer destDb.Close()

	rows, err1 := destDb.Query(sql_select)
	if err1 != nil {
		return results, err1
	}

	for rows.Next() {
		item := Server{}

		listIps := ""
		err = rows.Scan(&item.UUID, &item.Path, &item.Port, &listIps, &item.CreatedAt)
		if err == nil {
			item.ListIps = strings.Split(listIps, "||")
			results = append(results, item)
		}
	}

	return results, nil
}
