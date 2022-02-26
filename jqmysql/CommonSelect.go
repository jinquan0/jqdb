package jqmysql

import (
    "database/sql"
    "fmt"
    "strconv"
    "reflect";"errors";"unsafe"
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
    //fmt.Println(reflect.ValueOf(value))
    //fmt.Printf("%s\n", reflect.ValueOf(value) )
    //fmt.Println(reflect.ValueOf(value).Interface())
    //fmt.Printf("%s\n",reflect.ValueOf(value).Interface().(string) )
    int64, err := strconv.ParseInt(reflect.ValueOf(value).Interface().(string), 10, 64); err=err
    //fmt.Printf("%d\n",int64)
    field.SetInt( int64 )
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

/*
 MySQL select 结果集 转换为 任意类型(结构体)数组
 返回值: select 查询结果数量
*/
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
            
            //fmt.Println(col, colValue)

            if t.Field(i).Tag.Get("json") == col {
                field:= v.Field(i)
                switch field.Kind() {
                case reflect.Invalid:
                    return nil,0
                case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
                    SetIntegerField(field, colValue)
                case reflect.String:
                    SetUnexportedField(field, colValue)
                } 
            }
        }

        idx++
    }

    return nil,idx

}

/*
 @@注意： 
   结构体中的字段--> 首字母必须大写, 不然报错: 
   panic: reflect: reflect.flag.mustBeAssignable using value obtained using unexported field

   json标签--> 必须和 MySql表结构的field名称保持一致, 后续处理会比对json标签、table filed名称是否匹配

例如:
type ST_MyFields_1 struct{
    Id             int     `json:"id"`
    Product_key    string  `json:"product_key"`
    Verbose_name   string  `json:"verbose_name"`
}

表结构Field name 与结构体json 标签一致(字符大小写)
+--------------+--------------+------+-----+---------+-------+
| Field        | Type         | Null | Key | Default | Extra |
+--------------+--------------+------+-----+---------+-------+
| id           | int(11)      | NO   | PRI | NULL    |       |
| product_key  | varchar(255) | YES  |     | NULL    |       |
| verbose_name | varchar(255) | YES  |     | NULL    |       |
+--------------+--------------+------+-----+---------+-------+

*/

type ST_MyCommonQuery struct{
    AnyArray []interface{}
    Sql string
}

func MyCommonSelect(myQuery ST_MyCommonQuery, MyDb *sql.DB) (error,int) {  
    rows,err := MyDb.Query( myQuery.Sql )
    if err != nil {
        return err,0
    }else{
    }
    
    err1,numrow := Resultset2Anyarray(rows, myQuery.AnyArray)
    rows.Close()
    if err1 != nil {
        fmt.Println("MySQL/> DescAutoMatch error: %v", err1 )
        return err1,0
    }else{ 
    }

    ParaLock_my_read_total.Lock(); MyReadTotal++; ParaLock_my_read_total.Unlock()
    
    return nil,numrow
}

func MyCommonSelectV2(myQuery ST_MyCommonQuery, MyDb *sql.DB) (error, []map[string]interface{}, int) {  
    rows,err := MyDb.Query( myQuery.Sql )
    if err != nil {
        return err,0
    }else{
    }
    
    err1,numrow := Resultset2Anyarray(rows, myQuery.AnyArray)
    rows.Close()
    
    maps := make([]map[string]interface{}, 0)
    if err1 != nil {
        fmt.Println("MySQL/> DescAutoMatch error: %v", err1 )
        return err1,nil,0
    }else{
        for i:=0; i < numrow; i ++ {
            m := make(map[string]interface{})
            elem := reflect.ValueOf(myQuery.AnyArray[i]).Elem()
            relType := elem.Type()
            for j:=0; j < relType.NumField(); j++ {
                m[relType.Field(j).Name] = elem.Field(j).Interface()
            }
            maps = append(maps, m)
        }
        //fmt.Println(maps)
    }
    
    ParaLock_my_read_total.Lock(); MyReadTotal++; ParaLock_my_read_total.Unlock()
    
    return nil,maps,numrow
}
