package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	//"time"
)

func main() {
	db, err := sql.Open("mysql", "root:@/adolf?charset=utf8")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT article SET title=?,content=?,time=?")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	//更新数据
	stmt, err = db.Prepare("update article set content=? where id=?")
	checkErr(err)

	res, err = stmt.Exec("新研发部门", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM article")
	checkErr(err)

	for rows.Next() {
		var id int
		var title string
		var content string
		var time string
		err = rows.Scan(&id, &title, &content, &time)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(title)
		fmt.Println(content)
		fmt.Println(time)
	}

	//删除数据
	stmt, err = db.Prepare("delete from article where id=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
