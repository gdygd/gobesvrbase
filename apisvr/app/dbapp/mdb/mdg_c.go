package mdb

import (
	"apisvr/app/am"
	"context"
	"fmt"
)

// ------------------------------------------------------------------------------
// CreateTest
// ------------------------------------------------------------------------------
func (m *MariadbHandler) CreateTest(ctx context.Context, info am.TestVal) error {
	var query string = fmt.Sprintf(`
		INSERT INTO TestValT_TB (TEST_DT, VAL)
		VALUES (now(), %d)`, info.Val)

	db, dbErr := m.Open()
	var err error

	defer func() {
		m.Close(db)
	}()
	if dbErr != nil {
		am.Applog.Error("[CreateTest DB open error] : %s " + dbErr.Error())
		return dbErr
	}

	_, err = db.ExecContext(ctx, query)

	if err != nil {
		am.Applog.Error("[CreateTest Query error] : %s [%s] (%v)", err.Error(), query, info)
		return err
	}

	am.Applog.Print(5, "[CreateTest ok]")
	return nil
}
