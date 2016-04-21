package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	//请求等待超时时间
	MAXREQTIME int = 1000
	//请求队列最小阀值
	MINREQNUM int = 1000
)

//客户端请求队列
var reqlist []net.Conn

//请求处理定时器
var reqtime int = MAXREQTIME

//保存请求的队列
func saveReqQueue(client net.Conn) {
	reqlist = append(reqlist, client)
}

//处理请求的队列
func doClient() {
	//打开数据库资源
	db, err := sql.Open("mysql", "root:@/bigdata")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	defer db.Close()

	//已经处理请求的下标
	var rindex int

	for {
		//休眠，等待更多客户端的请求
		reqtime--
		//客户端数量还不够多并且没有达到处理的时间
		if MINREQNUM > len(reqlist) && reqtime > 0 {
			continue
		}
		reqtime = MAXREQTIME
		rindex = 0

		//没有客户端发起请求
		if len(reqlist) <= 0 {
			continue
		}

		//构造批量写入数据的请求
		sql := "INSERT INTO test_reg(r_nick, r_name, r_sex, r_phone, r_addr, r_recflag) VALUES "

		//遍历处理客户端请求
		var cflag string
		for _, client := range reqlist {
			sql += " ('qizexi', 'abc', '1', '13607765481', '广东省广州市白云区', MD5('123456')),"

			cflag = fmt.Sprintf("%v", client)
			//读取客户端发来的数据
			creader := bufio.NewReader(client)
			info, _ := creader.ReadString('\n')
			fmt.Println(info)

			//返回数据给客户端
			tt := int(time.Now().Unix())
			fmt.Fprintf(client, "%s\n", "hello :"+cflag+"#"+strconv.Itoa(tt))
			client.Close()
			rindex++
		}

		//去掉后面那个多余的逗号(,)
		sql = strings.Trim(sql, ",")
		//写入数据库
		stmt, err := db.Prepare(sql)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			stmt.Exec()
		}

		//切除已经处理的切片数据
		reqlist = reqlist[rindex:]
	}
}

func main() {
	//监听8080端口
	conn, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("监听8080端口失败：" + err.Error())
		return
	}

	//开启一个线程专门消费客户端的请求
	go doClient()

	//接受来自客户端的请求
	for {
		client, err := conn.Accept()
		if err != nil {
			fmt.Println("接受客户端失败：" + err.Error())
			continue
		}

		//开启一个线程来保存客户端的请求
		go saveReqQueue(client)
	}
}
