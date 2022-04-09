package jqmssql

import (
    "database/sql"
    "log"
    "reflect"
    "unsafe"
    //"strconv"
    "errors"
    _ "github.com/denisenkom/go-mssqldb"
)

//
// https://vimsky.com/examples/usage/reflect-pointer-function-in-golang-with-examples.html
// 针对string 类型的处理
func GetUnexportedField(field reflect.Value) interface{} { 
    return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface() 
} 
   
func SetUnexportedField(field reflect.Value, value interface{}) { 
    reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())). 
        Elem(). 
        Set(reflect.ValueOf(value)) 
} 

//
// https://blog.csdn.net/weixin_34310785/article/details/93156955
// 针对integer 类型的处理
func SetIntegerField(field reflect.Value, value interface{}) { 
    //log.Println(reflect.ValueOf(value))
    //log.Printf("%s\n", reflect.ValueOf(value) )
    //log.Println(reflect.ValueOf(value).Interface())
    //log.Printf("%d\n",reflect.ValueOf(value).Interface().(int64) )
    // SQL Server 直接返回整型数据，而不是整数 字符串
    field.SetInt( reflect.ValueOf(value).Interface().(int64) )
} 

// 针对double 类型的处理
func SetDoubleField(field reflect.Value, value interface{}) { 
    // SQL Server 直接返回浮点数据，而不是浮点数 字符串
    field.SetFloat( reflect.ValueOf(value).Interface().(float64) )
} 

func AnyElemType(Ptr interface{}) (reflect.Value, reflect.Type, error) {
    vPtr := reflect.ValueOf(Ptr)
    val := vPtr.Elem()
    typ := vPtr.Elem().Type()
    if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
        return val,typ,errors.New("dest should be a struct's pointer")
    }
    return val,typ,nil
}

func ST_GetValueByKey(stPtr interface{}, field_name string) interface{} {
   return reflect.ValueOf(stPtr).Elem().FieldByName(field_name)
}


func Resultset2Anyarray(rows *sql.Rows, structPtrs []interface{}) (error,int) {    
    columns, err := rows.Columns() 
    if err != nil {
        return err,0
    }
    count := len(columns)
    values := make([]interface{}, count)
    valuePtrs := make([]interface{}, count)

    var idx int = 0 
    for rows.Next() {
        v,t,_ := AnyElemType(structPtrs[idx])   

        for i := range columns {
            valuePtrs[i] = &values[i]
        }

        rows.Scan(valuePtrs...)

        for i, col := range columns {
            val := values[i]

            b, ok := val.([]byte)
            var colValue interface{}  ///// 反射类型
            if (ok) {
                colValue = string(b)
            } else {
                colValue = val
            }
            
            //log.Println(col, colValue);v=v;t=t
            
            if t.Field(i).Tag.Get("json") == col {
                field:= v.Field(i)
                switch field.Kind() {
                case reflect.Invalid:
                    return nil,0
                case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
                    SetIntegerField(field, colValue)
                case reflect.String:
                    SetUnexportedField(field, colValue)
                case reflect.Float32,reflect.Float64:
                    SetDoubleField(field, colValue)
                } 
                
            }
            
        }

        idx++
    }

    return nil,idx

}

type ST_MsCommonQuery struct{
    AnyArray []interface{}
    Sql string
}
func MsCommonSelect2Maps(msQuery ST_MsCommonQuery, MsDb *sql.DB) (error, []map[string]interface{}, int) {
    rows,err := MsDb.Query( msQuery.Sql )
    if err != nil {
        log.Println("MSSQL/> ", err )
        return err,nil,0
    }else{
    }

    err1,numrow := Resultset2Anyarray(rows, msQuery.AnyArray)
    rows.Close()
    
    maps := make([]map[string]interface{}, 0)
    if err1 != nil {
        log.Println("MSSQL/> DescAutoMatch error: ", err1 )
        return err1,nil,0
    }else{
        for i:=0; i < numrow; i ++ {
            m := make(map[string]interface{})
            elem := reflect.ValueOf(msQuery.AnyArray[i]).Elem()
            relType := elem.Type()
            for j:=0; j < relType.NumField(); j++ {
                m[relType.Field(j).Name] = elem.Field(j).Interface()
            }
            maps = append(maps, m)
        }
        //log.Println(maps)
    }
    
    return nil,maps,numrow
}
