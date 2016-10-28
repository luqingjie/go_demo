package main

import (
    "database/sql"
    "fmt"
    //需要在本地配置gobin，并且在gitbub上搞到驱动，并且本地编译通过，只要配置好
    //path,cmd下执行命令：go get github.com/go-sql-driver/mysql 
    //就可以再你配置的gobin下看到打包好的可以使用的代码 
    //项目主页 https://github.com/Go-SQL-Driver/MySQL ，里面的文档讲解的非常详细
    _ "github.com/go-sql-driver/mysql"
)
type userinfo struct {
    username    string
    email  string
    pwd   string
    id     int32
}
func main() {
    db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/db_go?charset=utf8")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    // user_info，需要替换你自己的表名
    rows, err := db.Query("SELECT * FROM user_info")
    if err != nil {
        panic(err.Error())
    }

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error())
    }

    // Make a slice for the values
    values := make([]sql.RawBytes, len(columns))

    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    // Fetch rows
    for rows.Next() {
        // get RawBytes from data
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }

        // Now do something with the data.
        // Here we just print each column as a string.
        var value string
        for i, col := range values {
            // Here we can check if the value is nil (NULL value)
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            fmt.Println(columns[i], ": ", value)
        }
        fmt.Println("-----------------------------------")

    }
            //插入数据
     stmt, err := db.Prepare("INSERT user_info SET username=?,email=?,pwd=?")
     checkErr(err)
     res, err := stmt.Exec("zhja", "研发","123")
     t_id,err:=res.LastInsertId();
     fmt.Println("insert lastid=",t_id)
    //result, err := db.Exec("INSERT INTO user_info (username, email, pwd) VALUES (?, ?, ?)","lily","销售","123")
     checkErr(err)
     //删除
    db.Exec("DELETE FROM user_info WHERE username=?", "zhja")
    checkErr(err)

    stmt1, err := db.Prepare("DELETE FROM user_info WHERE username=?")
    checkErr(err)
    stmt1.Exec("zhja")
    //查询
    var username, pwd, email string
    err = db.QueryRow("SELECT username,email,pwd FROM user_info WHERE id=?",1).Scan(&username, &email, &pwd)
    fmt.Println(username)
    fmt.Println(email)
    fmt.Println(pwd)

    rows1, err := db.Query("SELECT username,email,pwd FROM user_info WHERE username=?", "luqingjie")
    checkErr(err)
    for rows1.Next() {
        var username, email, pwd string
        if err := rows1.Scan(&username, &email, &pwd); err == nil {
            fmt.Println(err)
        }
        fmt.Println(username)
        fmt.Println(email)
        fmt.Println(pwd)
    }
    //事务,mysql表的类型需支持事务例：InnoDB
    tx, err := db.Begin()
    checkErr(err)
    stmt, err1 := tx.Prepare("INSERT INTO user_info (username, email, pwd) VALUES (?, ?, ?)")
    checkErr(err1)
    _, err2 := stmt.Exec("test", "测试", "123")
    checkErr(err2)
    //err3 := tx.Commit()
    err3 := tx.Rollback()
    checkErr(err3)
}
func checkErr(err error){
    if err != nil {
        panic(err)
    }
}