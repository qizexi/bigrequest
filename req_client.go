package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

//处理请求的队列
func doClient(client net.Conn) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("链接8080端口失败：" + err.Error())
		return
	}
	fmt.Fprintf(conn, "test\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("读取数据失败：" + err.Error())
		return
	}
	fmt.Println(status)

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
	conn, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("监听8081端口失败：" + err.Error())
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
