package mdb

import (
	"apisvr/app/am"
	"fmt"
)

// ------------------------------------------------------------------------------
// CreateTest
// ------------------------------------------------------------------------------
func (m *MariadbHandler) CreateTest(info am.TestVal) error {
	db, dbErr := m.Open()
	var err error

	defer func() {
		m.Close(db)
	}()
	if dbErr != nil {
		am.Applog.Error("[CreateTest DB open error] : %s " + dbErr.Error())
		return dbErr
	}

	var strQry string = fmt.Sprintf(`
		INSERT INTO TestValT_TB (TEST_DT, VAL)
		VALUES (now(), %d)`, info.Val)

	_, err = db.Exec(strQry)

	if err != nil {
		am.Applog.Error("[CreateTest Query error] : %s [%s] (%v)", err.Error(), strQry, info)
		return err
	}

	am.Applog.Print(5, "[CreateTest ok]")
	return nil
}
