package dbapp

import (
	"apisvr/app/am"
	"context"
	"database/sql"
)

type CancelContext func()
type DoneContext func()

type DBHandler interface {
	Open() (*sql.DB, error)
	OpenCtx() (*sql.DB, context.Context, CancelContext, DoneContext, error)
	OpenCtx2() (*sql.DB, context.Context, context.CancelFunc, error)
	Open2(dbNm string) (*sql.DB, error)
	Close(*sql.DB)
	Ping() (bool, error)

	// R
	ReadTest() ([]am.TestVal, error)

	// C
	CreateTest(info am.TestVal) error

	// U
	UpdateTest(val int) error

	// D
	DelTest(id int) error
}
