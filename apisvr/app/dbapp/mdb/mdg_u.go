package mdb

import (
	"apisvr/app/am"
	"fmt"
)

// ------------------------------------------------------------------------------
// UpdateTest
// ------------------------------------------------------------------------------
func (m *MariadbHandler) UpdateTest(val int) error {
	db, dbErr := m.Open()
	var err error

	defer func() {
		m.Close(db)
	}()
	if dbErr != nil {
		am.Applog.Error("[UpdateTest DB open error] : %s " + dbErr.Error())
		return dbErr
	}

	var strUpdQry string = fmt.Sprintf(`UPDATE TestValT_TB
							SET TEST_DT = now
							WHERE VAL = %d`, val)

	_, err = db.Exec(strUpdQry)

	if err != nil {
		am.Applog.Error("[UpdateTest Query error] : %s [%s]", err.Error(), strUpdQry)
		return err
	}

	am.Applog.Print(5, "[UpdateTest ok]")

	return nil
}
