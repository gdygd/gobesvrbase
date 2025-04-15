package mdb

import (
	"apisvr/app/am"
	"context"
	"fmt"
)

// ------------------------------------------------------------------------------
// UpdateTest
// ------------------------------------------------------------------------------
func (m *MariadbHandler) UpdateTest(ctx context.Context, val int) error {
	var query string = fmt.Sprintf(`UPDATE TestValT_TB
							SET TEST_DT = now
							WHERE VAL = %d`, val)
	db, dbErr := m.Open()
	var err error

	defer func() {
		m.Close(db)
	}()
	if dbErr != nil {
		am.Applog.Error("[UpdateTest DB open error] : %s " + dbErr.Error())
		return dbErr
	}

	// _, err = db.Exec(query)
	_, err = db.ExecContext(ctx, query)

	if err != nil {
		am.Applog.Error("[UpdateTest Query error] : %s [%s]", err.Error(), query)
		return err
	}

	am.Applog.Print(5, "[UpdateTest ok]")

	return nil
}
