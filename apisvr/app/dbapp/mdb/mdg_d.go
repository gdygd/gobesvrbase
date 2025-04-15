package mdb

import (
	"apisvr/app/am"
	"context"
	"fmt"
)

// ------------------------------------------------------------------------------
// DelTest
// ------------------------------------------------------------------------------
func (m *MariadbHandler) DelTest(ctx context.Context, id int) error {
	var qeury string = fmt.Sprintf(`delete from TestValT_TB where VAL = %d`, id)

	db, dbErr := m.Open()
	var err error

	defer func() {
		m.Close(db)
	}()
	if dbErr != nil {
		am.Applog.Error("[DelTest DB open error] : %s " + dbErr.Error())
		return dbErr
	}

	am.Applog.Print(2, "[DelTest] (%d)", id)

	// _, err = db.Exec(qeury)
	_, err = db.ExecContext(ctx, qeury)

	if err != nil {
		am.Applog.Error("[DelTest Query error] : %s [%s]", err.Error(), qeury)

		return err
	}

	am.Applog.Print(5, "[DelTest ok]")

	return nil
}
