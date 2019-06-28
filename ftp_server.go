package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const SERVER_IP string = "192.168.43.49:7788"
const FILE_ADDR string = "./"

func SendFile(path string, conn net.Conn, size string){
	defer conn.Close()
	//以只读打开文件
	fmt.Println(path)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	i, err := strconv.ParseInt(size, 10, 64)
	fmt.Println(i)
	buf := make([]byte, 1024*4)
	f.Seek(i,0)
	//读取文件内容
	for {
		n,err := f.Read(buf)
		fmt.Println(buf[:n])
		if err != nil{
			if err == io.EOF {
				fmt.Println("文件发送完毕")
			}else{
				fmt.Println(err)
			}
			return
		}
		conn.Write(buf[:n])
	}
}

func RecvFile(path string, conn net.Conn){
	defer conn.Close()
	//新建文件
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	buf := make([]byte, 1024*4)

	//接受文件内容
	for {
		n,err := conn.Read(buf)
		if err != nil{
			if err == io.EOF {
				fmt.Println("文件接受完毕")
			}else{
				fmt.Println(err)
			}
			return
		}
		//向文件写内容
		f.Write(buf[:n])
	}
}

func HandleChoice(conn net.Conn){
	defer conn.Close()
	buf := make([]byte,1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(buf[0]) == "1" {
		//path是服务器指定路径加文件名
		path := FILE_ADDR+string(buf[1:n])
		//确认文件是否存在,存在就提示，不存在就可以创建
		_, err = os.Open(path)
		if err != nil {
			//fmt.Println("难道是这里错了？？")
			_, err = conn.Write([]byte("1"))
			if err != nil {
				fmt.Println(err)
				return
			}
			RecvFile(path,conn)
		}else{
			fmt.Println("文件已经存在")
			_, err = conn.Write([]byte("4"))
			return
		}
	}else if string(buf[0]) == "2" {
		//path是服务器指定路径加文件名
		path := FILE_ADDR+string(buf[1:n])
		//确认文件是否存在
		_, err = os.Open(path)
		if err != nil {
			fmt.Println(err)
			_, err = conn.Write([]byte("3"))
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}else {
			//给请求方发送“ok"确认
			_, err = conn.Write([]byte("2"))
			if err != nil {
				fmt.Println(err)
				return
			}
			ok := make([]byte,1024)
			_, err := conn.Read(ok)
			//fmt.Println(string(ok[:n]))
			if err != nil {
				fmt.Println(err)
				return
			}
			if string(ok[:2]) == "ok" {
				fmt.Println(string(ok[2:n]))
				SendFile(path,conn,string(ok[2:n]))
			}else{
				fmt.Print("未收到接收方准备确认的信息，连接关闭")
			}
		}
	}else{
		fmt.Println("错误的请求！")
	}
}

func main(){
	//监听
	lisetener, err := net.Listen("tcp",SERVER_IP)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer lisetener.Close()

	for {
		//阻塞等待连接
		conn, err := lisetener.Accept()
		if err != nil {
			fmt.Print(err)
			return
		}
		defer conn.Close()
		go HandleChoice(conn)
	}
}