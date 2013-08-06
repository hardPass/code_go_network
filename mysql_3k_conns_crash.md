mysql 3000并发直接整个服务器挂掉
限制连接次数类似这样：

  n:=make(chan bool,100)
	func{
	n<-true
	do sql
	<-true
	}

限制他最大的并发数
2000以下
昨天头疼了我一个晚上
测试的时候2000机器人很好
3000开始直接整个服务器挂掉
一开始一直在找他驱动给的解决办法

	package main

	import (
		"database/sql"
		"fmt"
		_ "github.com/go-sql-driver/mysql"
		"math/rand"
		"runtime"
		"strconv"
		"time"
	)

	var  N = 2000
	var c = make(chan int, N)
	var db *sql.DB
	var limit = make(chan bool, 2 * runtime.NumCPU())
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))

	func main() {
		runtime.GOMAXPROCS(runtime.NumCPU())

		db, err := sql.Open("mysql", "root:kingweb@tcp(localhost:3306)/test?charset=utf8")
		if err != nil {
			fmt.Println("connect db fail %s", err)
		}
		db.SetMaxIdleConns(2 * runtime.NumCPU())

		t := time.Now()
		for i := 0; i < N; i++ {
			go updatemysql(db, i)
		}
		for i := 0; i < N; i++ {
			<-c
		}
		fmt.Println(time.Now().Sub(t))
	}

	func updatemysql(db *sql.DB, i int) {
		uid := r.Intn(1000)
		fmt.Println(uid)
		time.Sleep(time.Duration(uid) * time.Millisecond)
		stmt, err := db.Prepare("INSERT into test SET name='test" + strconv.Itoa(i) + "',num=?")
		if err != nil {
			fmt.Println(err)
		}
		res, err := stmt.Exec(uid)
		if err != nil {
			fmt.Println(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("insert test ", i, uid, id)
		c <- 1
	}
