//server
package main

import (
	// "bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	// "regexp"
	"strings"
)

var (
	logFileName = flag.String("log", "chatLog.log", "Log file name")
	//用来存储所有连接过服务器的conn,将连接用户和对应的conn做一个映射
	onlineConns = make(map[string]net.Conn)
	//用来记录当前链接用户IP
	// connNow = make(chan string)
	//初始化消息队列
	messageQueue = make(chan string, 1000)
	quitChan     = make(chan bool)
)

// //发送消息
// func SendInfo(conn net.Conn) {

// 	defer func() {
// 		if conn != nil {
// 			conn.Close()
// 		}
// 	}()

// 	//轮询发送消息
// 	var input string
// 	for {
// 		//获取用户输入,用户输入格式固定为，“127.0.0.1：50727#消息”
// 		reader := bufio.NewReader(os.Stdin)
// 		data, _, _ := reader.ReadLine()

// 		input = strings.ToUpper(string(data))

// 		if input == "EXIT" {
// 			conn.Close()
// 			os.Exit(0)
// 			break
// 		}

// 		_, senderr := conn.Write(data)

// 		checkErr(senderr)
// 	}
// }
func checkErr(err error) {
	if err != nil {
		//打开日志文件
		logfile, logerr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

		if logerr != nil {
			fmt.Println("Fail to find", *logfile)
			os.Exit(1)
		}
		log.SetOutput(logfile)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		//写入日志文件信息
		log.Printf("%s\n", err.Error())
		// fmt.Println("fail to do this!")
		// panic(err)
		// } else {
		// 	//打开日志文件
		// 	logfile, logerr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

		// 	if logerr != nil {
		// 		fmt.Println("Fail to find", *logfile)
		// 		os.Exit(1)
		// 	}
		// 	log.SetOutput(logfile)
		// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		// 	//写入日志文件信息
		// 	log.Printf("Success do it !\n")
	}
}

//处理接收的消息，存入消息队列
func ProcessInfo(conn net.Conn) {
	// fmt.Println("执行到函数ProcessInfo")
	writeToLog("执行到ProcessInfo")
	buf := make([]byte, 1024)
	defer func() {

		if conn != nil {
			conn.Close()
			writeToLog("执行完ProcessInfo")
		}
	}()
	//轮询
	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			continue
		}

		if numOfBytes != 0 {
			// //获取远程机器IP地址
			// remoteaddr := conn.RemoteAddr()
			// //如果不使用buf[:numOfBytes]直接使用buf输出的话可能会有多余字符
			// fmt.Printf("\t\tHas received this message from %s :%s\n", remoteaddr, string(buf[:numOfBytes]))

			message := string(buf[:numOfBytes])
			//将消息写入到消息队列中
			messageQueue <- message

			// currentIP := fmt.Sprintf("%s", conn.RemoteAddr())
			// connNow <- currentIP

			// cur := <-connNow
			// fmt.Printf("当前发送消息用户为：%s", cur)

		}
	}
}

//消费消息
func ConsumeMessage() {
	// fmt.Println("执行到函数ConsumeMessage")
	writeToLog("执行到ConsumeMessage")
	for {
		select {
		case message := <-messageQueue:
			//对消息进行解析
			doProcessMessage(message)
		case <-quitChan:
			break
		}
	}

	writeToLog("执行完ConsumeMessage")
}

//解析客户端消息
func doProcessMessage(message string) {
	// fmt.Println("执行到函数doProcessMessage")
	writeToLog("执行到函数doProcessMessage")
	contents := strings.Split(message, "#")
	// fmt.Println(len(contents))
	if len(contents) > 1 {
		//获取接受者主机地址
		addr := contents[0]

		//获取发送的消息
		//如果发送的消息中包含“#”,拼接消息
		sendMessage := strings.Join(contents[1:], "#")

		fmt.Println(sendMessage)
		//处理字符串中的空格
		addr = strings.Trim(addr, " ")

		//从链接map中查找接收消息的IP地址是否已经登录
		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMessage))

			if err != nil {
				fmt.Println("online conns send faliure")
			}
		} else {

		}
	}
	writeToLog("执行完成函数doProcessMessage")
	// fmt.Println("执行完成函数doProcessMessage")
}

//写入日志运行信息
func writeToLog(info string) {
	//打开日志文件
	logFile, logErr := os.OpenFile("chatLog.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if logErr != nil {
		fmt.Println(logErr.Error())
	}
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()
	//设置日志模板
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//写入日志
	log.Printf("%s\n", info)
}

func main() {
	//初始化map
	// onlineConns = make(map[string]net.Conn)

	//c/s编程首先开启监听
	listen_socket, err := net.Listen("tcp", "127.0.0.1:8080")
	writeToLog("开启监听")
	//不要忘记关闭socket
	defer func() {
		if listen_socket != nil {
			listen_socket.Close()
			writeToLog("套接字已关闭")
		}
	}()
	checkErr(err)

	fmt.Println("server is starting!")

	//开启协程消费messageQueue
	go ConsumeMessage()
	//轮询接收消息
	for {
		conn, connerr := listen_socket.Accept()
		checkErr(connerr)

		//将conn存储到映射表中
		remoteaddr := conn.RemoteAddr()
		//将地址转换成字符串
		addr := fmt.Sprintf("%s", remoteaddr)
		onlineConns[addr] = conn

		//查看map中存储了那些conn
		for item := range onlineConns {
			fmt.Println(item)
		}
		//并发接收请求
		go ProcessInfo(conn)

		//并发发送消息
		// go SendInfo(conn)
	}

}
