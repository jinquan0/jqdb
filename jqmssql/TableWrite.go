package jqmssql

import (
	"context"
	"log"
	"reflect"
    "crypto/rand"  //生成真随机数
    "github.com/satori/go.uuid" //生成UID
    _ "github.com/denisenkom/go-mssqldb"
)

// sql example:
/* create table tb1 (id int identity(101,1) primary key not null, 
    fld1 int not null,
 	fld2 nvarchar(50) not null
 	)
 */
func CreateTable(conn DBconn, tabname string, sql string) {
    db := MssqlConn(conn)
    //conn := MssqlConnWithMSA(*server, *port, *user, *database)
	_, err := db.Exec(sql)
	defer db.Exec("drop table " + tabname)
	if err != nil {
		log.Fatal("SQL Server/> create table failed with error", err)
	}
	MssqlDisconn(db)
}

func random_data() (int64, string) {
	num, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	uid := uuid.NewV4().String()
	return num.Int64(),uid
}

// sql example:
// insert into tb1 (fld1, fld2) values (123, 'abcdefg')
func BulkInsertTable(conn DBconn, tabname string, count int) {
	db := MssqlConn(conn)
	for i := 0; i < count; i++ {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal("SQL Server/> Begin tran failed", err)
		}

		num,uid:=random_data(); sql:=fmt.Sprintf("insert into "+tabname+" (fld1, fld2) values(%d, '%s')", num, uid)
		
		_, err = tx.Exec(sql)
		if err != nil {
			log.Fatal("SQL Server/> Insert failed", err)
		}
	}
	MssqlDisconn(db)
}

// sql example:
// sql_tx1: "insert into tb1 (fld1, fld2) values (1,'aaa')", 
// sql_tx2: "insert into tb1 (fld1, fld2) values (2,'bbb')", 
// sql_post_tx: "select fld1,fld2 from tb1"
func Transaction(conn DBconn, sql_tx1 string, sql_tx2 string, sql_post_tx string) {
    //db := MssqlConn(*server, *port, *user, *password, *database)
    db := MssqlConn(conn)

	tx1, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatal("BeginTx failed with error", err)
	}
	tx2, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatal("BeginTx failed with error", err)
	}
	_, err = tx1.Exec(sql_tx1)
	if err != nil {
		log.Fatal("tx1 failed with error", err)
	}
	_, err = tx2.Exec(sql_tx2)
	if err != nil {
		log.Fatal("tx2 failed with error", err)
	}
	tx1.Rollback()
	tx2.Commit()

	if sql_post_tx != "" {
		rows, err := db.Query(sql_post_tx)
		if err != nil {
			log.Fatal("post_tx read failed with error", err)
		}
		defer rows.Close()
		values := []int64{}
		for rows.Next() {
			var val int64
			err = rows.Scan(&val)
			if err != nil {
				log.Fatal("post_tx scan failed with error", err)
			}
			values = append(values, val)
		}
		if !reflect.DeepEqual(values, []int64{2}) {
			log.Fatal("Values is expected to be [1] but it is %v", values)
		}
	}

	MssqlDisconn(db)
}


