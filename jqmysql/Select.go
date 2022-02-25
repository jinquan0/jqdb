package jqmysql

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)

func MySelect(MyDb *sql.DB, query string, args... interface{}) bool {

    row:=MyDb.QueryRow(query)

    err:=row.Scan(args...)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("MySQL/> There were no rows, but otherwise no error occurred")
            return false
        }else{
            fmt.Println("MySQL/> Query err: %v", err)
            return false
        }
    }

    ParaLock_my_read_total.Lock(); MyReadTotal++; ParaLock_my_read_total.Unlock()

    return true
}
