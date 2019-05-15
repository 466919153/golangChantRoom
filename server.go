package main

import (
	"fmt"

	"log"
	"net"
)
type toMessage struct {
	from string  //谁来发送
	to string  //给谁发
	privMessage string
}
var(
	chLogin chan string
	chLogout chan string
	chMessage chan string
	chToMessage chan toMessage
)

var mapAllClient map[string]net.Conn

func sendMsg(conn net.Conn,msg string){
	num,err:=conn.Write([]byte(msg))
	if err!=nil||num<=0 {
		fmt.Println("client err:",err)

		delete(mapAllClient,conn.RemoteAddr().String())

		chLogout<-conn.RemoteAddr().String()+"不在线了"

	}

}


func chHandle(){
	var msg string
	var to toMessage
	for {
		select {
		case msg =<-chLogin:
			for _,v:=range mapAllClient{
				go sendMsg(v,msg)
			}
		case msg =<-chLogout:
			for _,v:=range mapAllClient{
				go sendMsg(v,msg)
			}
		case msg =<-chMessage:
			for _,v:=range mapAllClient{
				go sendMsg(v,msg)
			}
		case to=<-chToMessage:
			if v,ok:=mapAllClient[to.to];ok{
				fmt.Println(v,ok)
				 sendMsg(v,to.privMessage)
			}else {
				 from:=mapAllClient[to.from]
				 sendMsg(from,"您地址输入错误或者对方不在线")
			}




		}

	}

}



func withClient(conn net.Conn){

	defer conn.Close()
	buf:=make([]byte,256)
	for {
		num,err:=conn.Read(buf)
		if err!=nil||num<=0{
			fmt.Println("Read(buf):",err)
			break
		}
		//写入消息

		if buf[0]=='@'{

			for i:=1; i<num;i++{
				if buf[i]=='#' {
					//@对方地址#消息
					chStruct :=toMessage{conn.RemoteAddr().String(),string(buf[1:i]),conn.RemoteAddr().String()+"对你说"+string(buf[i+1:num])}

					chToMessage<-chStruct

					break
				}


			}

		}else {
			chMessage<-conn.RemoteAddr().String()+"说"+string(buf[:num])
		}


	}
	delete(mapAllClient,conn.RemoteAddr().String())
	chLogout<-conn.RemoteAddr().String()+"下线了"
}



func main(){
	mapAllClient=make(map[string]net.Conn)
	chLogin=make(chan string)
	chLogout=make(chan string)
	chMessage=make(chan string)
	chToMessage=make(chan toMessage)


	lis,err:=net.Listen("tcp","localhost:9999")

	if err!=nil{
		log.Fatal("listen fault:",err)
	}
	defer lis.Close()

	go chHandle() //消息选择
	for {
		conn,err:=lis.Accept()
		if err!=nil{
			log.Fatal("accpet fault:",err)
			continue
		}
		connstr:=conn.RemoteAddr().String()

		mapAllClient[connstr]=conn

		chLogin<-connstr+"上线了"

		//通信
		go withClient(conn)


	}


}