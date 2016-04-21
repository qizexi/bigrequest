package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//处理请求的队列
func doClient(client net.Conn) {
	//打开数据库连接
	db, err := sql.Open("mysql", "root:@/bigdata")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	defer db.Close()

	//单条数据的写入
	sql := "INSERT INTO test_reg(r_nick, r_name, r_sex, r_phone, r_addr, r_recflag) VALUES "
	sql += " ('qizexi', 'abc', '1', '13607765481', '广东省广州市白云区', MD5('123456'))"

	//写入数据库
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		stmt.Exec()
	}

	var cflag string
	cflag = fmt.Sprintf("%v", client)
	//读取客户端发来的数据
	creader := bufio.NewReader(client)
	info, _ := creader.ReadString('\n')
	fmt.Println(info)

	//返回数据给客户端
	tt := int(time.Now().Unix())
	fmt.Fprintf(client, "%s\n", "hello :"+cflag+"#"+strconv.Itoa(tt))
	client.Close()

}

func main() {
	//监听8081端口
	conn, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("监听8082端口失败：" + err.Error())
		return
	}

	//接受来自客户端的请求
	for {
		client, err := conn.Accept()
		if err != nil {
			fmt.Println("接受客户端失败：" + err.Error())
			continue
		}

		//开启一个线程来保存客户端的请求
		go doClient(client)
	}
}
