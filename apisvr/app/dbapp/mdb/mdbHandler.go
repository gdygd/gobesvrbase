package mdb

import (
	"context"
	"database/sql"
	"time"

	"apisvr/app/am"

	"apisvr/app/dbapp"

	_ "github.com/go-sql-driver/mysql"
)

var Temp dbapp.DBHandler = nil

// ------------------------------------------------------------------------------
// Struct
// ------------------------------------------------------------------------------
type MariadbHandler struct {
	user      string
	pw        string
	dbNm      string
	host      string
	Connected bool
}

// ------------------------------------------------------------------------------
// Open
// ------------------------------------------------------------------------------
func (m *MariadbHandler) Open() (*sql.DB, error) {
	dbSrc := m.user + ":" + m.pw + "@tcp(" + m.host + ")/" + m.dbNm
	database, err := sql.Open("mysql", dbSrc)

	if err != nil {
		am.Applog.Error("database open err : %v", err)
		m.Connected = false
		return nil, err
	}

	// var ctx context.Context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = database.PingContext(ctx); err != nil {
		am.Applog.Error("database open ping context err : %v", err)
		m.Connected = false
		return nil, err
	}

	m.Connected = true

	return database, nil
}

// ------------------------------------------------------------------------------
// Open
// ------------------------------------------------------------------------------
func (m *MariadbHandler) OpenCtx() (*sql.DB, context.Context, dbapp.CancelContext, dbapp.DoneContext, error) {

	return nil, nil, nil, nil, nil

}

func (m *MariadbHandler) OpenCtx2() (*sql.DB, context.Context, context.CancelFunc, error) {
	return nil, nil, nil, nil
}

// ------------------------------------------------------------------------------
// Open
// ------------------------------------------------------------------------------
func (m *MariadbHandler) Open2(dbNm string) (*sql.DB, error) {
	dbSrc := m.user + ":" + m.pw + "@tcp(" + m.host + ")/" + dbNm
	database, err := sql.Open("mysql", dbSrc)

	if err != nil {
		am.Applog.Error("database open err : %v", err)
		m.Connected = false
		return nil, err
	}

	// var ctx context.Context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = database.PingContext(ctx); err != nil {
		am.Applog.Error("database open ping context err : %v", err)
		m.Connected = false
		return nil, err
	}

	m.Connected = true

	return database, nil
}

// ------------------------------------------------------------------------------
// Close
// ------------------------------------------------------------------------------
func (m *MariadbHandler) Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// ------------------------------------------------------------------------------
// Ping
// ------------------------------------------------------------------------------
func (m *MariadbHandler) Ping() (bool, error) {

	var err error
	var rows *sql.Rows = new(sql.Rows)
	db, dbErr := m.Open()

	defer func() {
		am.Applog.Print(6, "[Ping close..]")
		m.Close(db)
		rows.Close()
	}()

	if dbErr != nil {
		am.Applog.Error("[Ping db open error] : %s " + dbErr.Error())
		return false, dbErr
	}

	am.Applog.Print(6, "[Ping qeury...]")
	var val int = 0
	rows, err = db.Query(`select 1 VAL from dual`)
	if err != nil {
		am.Applog.Print(6, "[Ping qeury...] err %v", err)
		return false, err
	}

	if rows.Next() {
		err := rows.Scan(&val)

		if err != nil {
			am.Applog.Error("[Ping Query error(1)] : %s" + err.Error())
			return false, err
		}
	}

	am.Applog.Print(6, "[Ping qeury...](%d)", val)

	return true, nil
}

// ------------------------------------------------------------------------------
// ChangeHostAddress
// ------------------------------------------------------------------------------
func (m *MariadbHandler) ChangeHostAddress(host string) {
	m.host = host
}

// ------------------------------------------------------------------------------
// GetConnected
// ------------------------------------------------------------------------------
func (m *MariadbHandler) GetConnected() bool {
	return m.Connected
}

// ------------------------------------------------------------------------------
// NewMariadbHandler
// ------------------------------------------------------------------------------
func NewMariadbHandler(user, pw, db, host string) *MariadbHandler {

	return &MariadbHandler{user: user, pw: pw, dbNm: db, host: host, Connected: false}
}
