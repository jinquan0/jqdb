package jqmysql

import (
    "database/sql"
    "fmt"
)

func MySelect(MyDb *sql.DB, sql string, args... interface{}) bool {

    row:=MyDb.QueryRow(sql)

    err:=row.Scan(args...)
    if err != nil {
        fmt.Println("MySQL/> Query err: %v", err)
        return false
    }

    ParaLock_my_read_total.Lock(); MyReadTotal++; ParaLock_my_read_total.Unlock()

    return true
}