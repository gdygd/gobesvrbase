package mdb

import (
	"apisvr/app/am"
	"fmt"
)

// ------------------------------------------------------------------------------
// DelTest
// ------------------------------------------------------------------------------
func (m *MariadbHandler) DelTest(id int) error {
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

	var strQry string = fmt.Sprintf(`delete from TestValT_TB where VAL = %d`, id)

	_, err = db.Exec(strQry)

	if err != nil {
		am.Applog.Error("[DelTest Query error] : %s [%s]", err.Error(), strQry)

		return err
	}

	am.Applog.Print(5, "[DelTest ok]")

	return nil
}
