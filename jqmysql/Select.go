package jqmysql

import (
    "database/sql"
    "fmt"
)

func MySelect(MyDb *sql.DB, sql string, args... interface{}) bool,n_row {

    row:=MyDb.QueryRow(sql)

    err:=row.Scan(args...)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("MySQL/> There were no rows, but otherwise no error occurred")
            return false,0   
        }else{
            fmt.Println("MySQL/> Query err: %v", err)
            return false,0
        }
    }

    ParaLock_my_read_total.Lock(); MyReadTotal++; ParaLock_my_read_total.Unlock()

    return true,1
}
