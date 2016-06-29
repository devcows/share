package lib

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/cznic/ql/driver"
)

type Server struct {
	UUID      string    `json:"uuid"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	ListIps   []string  `json:"list_ips"`
	Flags     []string  `json:"flags"`
}

var settings *SettingsShare

func OpenDatabase() (*sql.DB, error) {
	destDb, err := sql.Open("ql", settings.ShareDaemon.DatabaseFilePath)
	if err != nil {
		return nil, err
	}
	destDb.Ping()

	return destDb, nil
}

func InitDB(settings_params SettingsShare) error {
	settings = &settings_params
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
		Flags string,
		CreatedAt time
	);
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}

	defer destDb.Close()

	tx, err := destDb.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(sql_table); err != nil {
		return err
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
		Flags,
		CreatedAt,
	) values($1, $2, $3, $4)
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}
	defer destDb.Close()

	tx, err := destDb.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sqlAdd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(server.UUID, server.Path, strings.Join(server.Flags, "||"), server.CreatedAt)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func FindServer(uuid string) (Server, error) {
	return Server{}, nil
}

func RemoveServer(uuid string) error {
	sqlRemove := `
	delete from Servers where
	UUID == $1
	`

	if _, err := FindServer(uuid); err != nil {
		return err
	}

	destDb, err := OpenDatabase()
	if err != nil {
		return err
	}
	defer destDb.Close()

	tx, err := destDb.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sqlRemove)
	if err != nil {
		return err
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

	sqlSelect := `
	select UUID, Path, Flags, CreatedAt from Servers order by CreatedAt
	`

	destDb, err := OpenDatabase()
	if err != nil {
		return results, err
	}
	defer destDb.Close()

	rows, err := destDb.Query(sqlSelect)
	if err != nil {
		return results, err
	}

	for rows.Next() {
		item := Server{}
		listFlags := ""
		if err = rows.Scan(&item.UUID, &item.Path, &listFlags, &item.CreatedAt); err == nil {
			item.Flags = strings.Split(listFlags, "||")
			results = append(results, item)
		}
	}

	return results, nil
}

func SearchServerByUUID(uuid string) (Server, error) {
	servers, err := ListServers()

	if err != nil {
		return Server{}, err
	}

	for _, server := range servers {
		if server.UUID == uuid {
			return server, nil
		}
	}

	return Server{}, fmt.Errorf("Server not found with UUID: %s", uuid)
}
