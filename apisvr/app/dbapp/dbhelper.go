package dbapp

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// ------------------------------------------------------------------------------
// ChekcNullString
// ------------------------------------------------------------------------------
func ChekcNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s, Valid: true,
	}
}

// ------------------------------------------------------------------------------
// zeroNullInt
// ------------------------------------------------------------------------------
func zeroNullInt(n int) sql.NullInt32 {
	num := int32(n)
	if num == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: num, Valid: true,
	}
}

// ------------------------------------------------------------------------------
// ZeroNullIntStr
// ------------------------------------------------------------------------------
func ZeroNullIntStr(n int) string {
	num := int32(n)
	if num == 0 {
		return "NULL"
	}
	return strconv.Itoa(n)
}

// ------------------------------------------------------------------------------
// emptyNullStr
// ------------------------------------------------------------------------------
func EmptyNullStr(s string) string {

	if len(s) == 0 {
		return "NULL"
	}
	return "'" + s + "'"
}

// ------------------------------------------------------------------------------
// ZeroNullFloatStr
// ------------------------------------------------------------------------------
func ZeroNullFloatStr(n float64) string {
	num := float64(n)
	if num == 0 {
		return "NULL"
	}
	//return strconv.Itoa(n)
	return fmt.Sprintf("%f", n)

}

func ServSysdate() string {
	t := time.Now()
	strServTm := fmt.Sprintf("TO_DATE('%04d%02d%02d%02d%02d%02d','YYYYMMDDHH24MISS')\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return strServTm
}

// ------------------------------------------------------------------------------
//
//	ConvertIntArrToStrArr
//
// ------------------------------------------------------------------------------
func ConvertIntArrToStrArr(data []int) []string {
	var strSlice []string
	for _, num := range data {
		strSlice = append(strSlice, fmt.Sprintf("%d", num))
	}
	return strSlice
}

// ------------------------------------------------------------------------------
//
//	MakeQuery
//
// ------------------------------------------------------------------------------
func MakeQuery(tmplname, query string, data any) (string, error) {
	t, err := template.New(tmplname).Parse(query)
	if err != nil {
		return "", err
	}

	var queryBuilder strings.Builder
	err = t.Execute(&queryBuilder, data)
	if err != nil {
		return "", err
	}

	return queryBuilder.String(), nil
}
