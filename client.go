package main


import (
	"fmt"
	"net"
	"os"
)


var ch chan int



func reader(conn *net.TCPConn) {
	buff := make([]byte, 256)
	for {
		j, err := conn.Read(buff)
		if err != nil {
			ch <- 1
			break
		}
		fmt.Printf("%s\n", buff[0:j])
	}
}


func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:%s  请输入host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1] //命令行 输入IP：port
	TcpAdd, _ := net.ResolveTCPAddr("tcp", service)

	conn, err := net.DialTCP("tcp", nil, TcpAdd)
	if err != nil {
		fmt.Println("服务没打开或者服务器故障")
		os.Exit(1)
	}
	defer conn.Close()

	go reader(conn)



	for {
		var msg string
		select {

		case <-ch:
			fmt.Println("server发生错误，请重新连接")
			os.Exit(2)
		default:
			fmt.Scan(&msg)
			bMsg := []byte( msg)
			conn.Write(bMsg)
		}
	}
}