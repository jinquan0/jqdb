package main

import (
    "flag"
    _ "github.com/denisenkom/go-mssqldb"
    jq "github.com/jinquan0/jqdb/jqmssql"
)

var (
    debug         = flag.Bool("debug", false, "enable debugging")
    server        = flag.String("server", "172.24.22.4", "SQL Server Host IP")
    port     *int = flag.Int("port", 1433, "SQL Server connect port")
    user          = flag.String("user", "testuser0", "connect user")
    password      = flag.String("password", "FDNvQ8#g", "connect password")
    //user          = flag.String("user", "cn\\infra01testuser", "connect user")
    database      = flag.String("database", "jq01test", "database name")
)

var DBC_PARA jq.DBconn

func init() {
	flag.Parse()
	DBC_PARA.Server = *server
	DBC_PARA.Port = *port 
	DBC_PARA.User = *user 
	DBC_PARA.Pass = *password
	DBC_PARA.Database = *database
}

func main() {
//	jq.CreateTable(DBC_PARA, 
//	`create table table1 (id int identity(101,1) primary key not null, 
//     fld1 int not null,
// 	 fld2 nvarchar(50) not null
// 	)`)

// 	jq.BulkInsertTable(DBC_PARA, "table1", 10)
	
	for i:=101; i < 120; i++ {
		jq.RandomUpdateTable(DBC_PARA, "table1", i)
	}
}
