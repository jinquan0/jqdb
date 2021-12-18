package jqmysql

import (
    "database/sql"
    mysql "github.com/go-sql-driver/mysql"
    
    "fmt"
    "time"
    "strings"
    
    "io/ioutil"
    "log"
    "crypto/tls";"crypto/x509"
)

type ST_MySQL_Endpoint struct {
  MyEndpoint string
  Host  string  
  Port  string  
  User  string  
  Pass  string  
  Sslca   string  
  Db    string  
}

func MyEndpointConstruct(conn *ST_MySQL_Endpoint) {
    if conn.Sslca != "null" {
        conn.MyEndpoint = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true&tls=custom&charset=utf8", conn.User , conn.Pass, conn.Host, conn.Port, conn.Db)
    }else{
        conn.MyEndpoint = strings.Join([]string{conn.User, ":", conn.Pass, "@tcp(", conn.Host, ":", conn.Port, ")/", conn.Db, "?charset=utf8"}, "")
    }
    fmt.Println(conn.MyEndpoint)
}

func MyConn(conn *ST_MySQL_Endpoint) *sql.DB {
    var err error

    ParaLock_my_conn_total.Lock(); MyConnErr++; ParaLock_my_conn_total.Unlock()
    
    // 构建MySQL连接串
    MyEndpointConstruct(conn)

    if conn.Sslca != "null" {
        rootCertPool := x509.NewCertPool()
        pem, _ := ioutil.ReadFile(conn.Sslca)
        if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
            log.Fatal("Failed to append PEM.")
        }
        mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})
    }else{
    }


    var MyDb *sql.DB
    MyDb,err = sql.Open("mysql", conn.MyEndpoint)
    if err != nil{
        r:=fmt.Sprintf("MySQL/> Open mysql failed, err:%v\n", err);
        fmt.Println(r)
        ParaLock_my_conn_err.Lock(); MyConnErr++; ParaLock_my_conn_err.Unlock()
    } else {
        //最大连接周期，超过100秒 的连接就close
        MyDb.SetConnMaxLifetime(100*time.Second)
        //设置最大连接数
        //MyDb.SetMaxOpenConns(1000) 
        //设置数据库最大闲置连接数
        //MyDb.SetMaxIdleConns(100)
    }

    //验证连接
    if err = MyDb.Ping(); err != nil {
        fmt.Println("MySQL/> open database fail!")
        return nil
    }
    if FlagDbg == true {
        fmt.Println("MySQL/> Connnect success")
    }
    return MyDb
}

func MyDisconn(MyDb *sql.DB){
    MyDb.Close()
    if FlagDbg == true {
        fmt.Println("MySQL/> database disconnected.")
    }
}