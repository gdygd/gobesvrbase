package mdb

import (
	"apisvr/app/am"
)

// ------------------------------------------------------------------------------
// ReadGroup
// ------------------------------------------------------------------------------
func (m *MariadbHandler) ReadTest() ([]am.TestVal, error) {
	db, dbErr := m.Open()

	var rdata am.TestVal = am.TestVal{}

	defer m.Close(db)
	if dbErr != nil {
		am.Applog.Error("[ReadGroup DB open error] : %s " + dbErr.Error())
		return nil, dbErr
	}

	rows, err := db.Query(`SELECT now() as dt, 1 as val from dual`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	datas := make([]am.TestVal, 0)

	for rows.Next() {

		err := rows.Scan(&rdata.Dt, &rdata.Val)

		if err != nil {
			am.Applog.Error("[ReadTest Query error(1)] : %s" + err.Error())
			return nil, err
		}
		datas = append(datas, rdata)

	}
	if err = rows.Err(); err != nil {
		am.Applog.Error("[ReadTest Query error(2)] : %s " + err.Error())
		return nil, err
	}
	am.Applog.Print(2, "[ReadTest ok]")

	return datas, nil
}
