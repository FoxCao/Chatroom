//client

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		//打开日志文件
		logfile, logerr := os.OpenFile("charLog_c.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

		if logerr != nil {
			panic(logerr)
		}

		//设置日志内容
		log.SetOutput(logfile)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		log.Printf("Fail to do this!")
		fmt.Println("Fail to do this!")
	}

}

//发送信息
func MessageSend(conn net.Conn) {

	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}

		_, err := conn.Write([]byte(input))

		if err != nil {
			conn.Close()

			//打开日志文件
			logfile, logerr := os.OpenFile("charLog_c.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

			if logerr != nil {
				panic(logerr)
			}

			//设置日志内容
			log.SetOutput(logfile)
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

			log.Printf("Fail to do this!:%s\n", err.Error())
			// fmt.Println("Fail to do this!")

			break

		}
	}
}
func main() {

	//拨号连接获取套接字
	conn, connerr := net.Dial("tcp", "127.0.0.1:8080")

	//关闭套接字
	defer conn.Close()

	//错误检查
	checkErr(connerr)

	go MessageSend(conn)
	// conn.Write([]byte("你好server，我是client"))

	buf := make([]byte, 1024)

	for {
		numOfBytes, err := conn.Read(buf)
		// checkErr(err)
		if err != nil {
			continue
		}
		//判断是否获取到数据
		if numOfBytes != 0 {
			//获取远程主机IP地址
			remoteAddr := conn.RemoteAddr()

			fmt.Printf("\treceived message from %s :%s\n", remoteAddr, string(buf[:numOfBytes]))
		}
	}
	fmt.Println("client program end!")
}
