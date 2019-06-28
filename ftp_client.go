package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const SERVER_IP string = "192.168.43.49:7788"

//发送文件
func SendFile(path string, conn net.Conn){
	//以只读打开文件
	fmt.Println(path)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	buf := make([]byte, 1024*4)
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

//接受文件
func RecvFile(fileName string, conn net.Conn, size string){
	//新建文件
	f, err := os.OpenFile(fileName,os.O_RDWR | os.O_APPEND | os. O_CREATE,066)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	i, err := strconv.ParseInt(size, 10, 64)
	buf := make([]byte, 1024*4)
	f.Seek(i,0)

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

func communication(choice string){
	//提示输入文件
	fmt.Println("请输入文件名：")
	var path string
	fmt.Scan(&path)

	//获取文件名 info.name()
	//info, err := os.Stat(path)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(1,info)
	//主动连接服务器
	conn, err := net.Dial("tcp",SERVER_IP)
	if err != nil {
		fmt.Println(err)
		return
	}
	//延迟关闭
	defer conn.Close()

	//给服务方发送文件名
	_, err = conn.Write([]byte(choice+path))
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	if "1" == string(buf[:n]){
		SendFile(path,conn)
	}else if "2" == string(buf[:n]){
		//判断本地是否存在文件
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println(info.Size())
			//文件存在,发送准备好了信息
			_, err = conn.Write([]byte("ok"+"0"))
			if err != nil {
				fmt.Println(err)
				return
			}
			RecvFile(info.Name(),conn,string(info.Size()))
		}
		fmt.Println(info.Size())
		//文件存在,发送准备好了信息
		_, err = conn.Write([]byte("ok"+string(info.Size())))
		if err != nil {
			fmt.Println(err)
			return
		}
		RecvFile(info.Name(),conn,string(info.Size()))
	}else if "3" == string(buf[:n]){
		//文件不存在
		//fmt.Println(3,string(buf[:n]))
		fmt.Println("文件不存在")
	} else{
		//文件存在
		fmt.Println(4,string(buf[:n]))
		fmt.Println("文件已经存在！换个名字！")
	}
}

func main(){
	LABEL:
	//提示输入恩建
	fmt.Println("请选择要进行的操作：\n1.上传文件\n2.下载文件\n")
	choice := make([]byte, 2)
	_, err := os.Stdin.Read(choice)
	if err != nil {
		fmt.Println(err)
		return
	}
	if choice[0] == '1' {
		communication("1")
	}else if choice[0] == '2'{
		communication("2")
	}else{
		fmt.Println("输入有误!")
		goto LABEL //回到开始
	}
}
