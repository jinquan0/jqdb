package jqmysql

import (
    "database/sql"
    "log"
    "reflect"
)

func MyInsert(MyDb *sql.DB, sql string, mydata interface{}, info_flag bool) bool {

    rv := reflect.ValueOf(mydata)
    var args []interface{}
    for i := 0; i < rv.NumField(); i++ {
        args = append(args, rv.Field(i).Interface())
    }
    //log.Println(args)
    //log.Println(sql)

    //开启事务
    tx, err := MyDb.Begin()
    if err != nil {
        log.Println("MySQL/> transaction fail")
        return false
    }

    //准备sql语句, IGNORE: no gaps
    stmt, err := tx.Prepare(sql )
    if err != nil {
        log.Println("MySQL/> Prepare fail")
        return false
    }
    //将参数传递到sql语句中并且执行
    // @@@ 在一个slice后加上 … ，就能把slice拆开成一个可变参数的形式传入函数
    res, err := stmt.Exec(args...)
    if err != nil {
        log.Println("MySQL/> Exec fail")
        return false
    }

    //将事务提交
    tx.Commit()
    
    //获得上一个插入自增的id
    i,_:=res.LastInsertId()
    if info_flag == true {
        log.Printf("MySQL/> LastInsertId[%d]: %v \n", i, mydata )
    }
    
    ParaLock_my_write_total.Lock(); MyWriteTotal++; ParaLock_my_write_total.Unlock()
    return true
}
