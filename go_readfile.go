package main

import (
    "fmt"
    "io"
    "os"
	"log"
)

func rw(){
	fp,err:=os.OpenFile("a.txt",os.O_WRONLY,0644)
	if err!=nil{
	fmt.Println("文件打开失败")
	return
	}else{
	content:="222222222222\r\n" 
	n,_:=fp.Seek(0,2)
	_,err=fp.WriteAt([]byte(content),n)
	}
	defer fp.Close()
}

func w1(){
	//写文件1
	fp,err := os.Create("./a1.txt")
	if err!=nil{
		//文件创建失败
		/*
		1.路径不存在
		2.文件权限
		3.程序打开文件上限
		 */
		fmt.Println("文件创建失败")
		return
	}
 
	//\n不会换行  原因 在windows文本文件中换行\r\n  回车  在linux中换行\n
	fp.WriteString("hello world\r\n")
	fp.WriteString("性感荷官在线发牌")
	defer fp.Close()
}


func w2(){
	//写文件2
	fp,err := os.Create("./a2.txt")
	if err!=nil{
		//文件创建失败
		/*
		1.路径不存在
		2.文件权限
		3.程序打开文件上限
		 */
		fmt.Println("文件创建失败")
		return
	}

	//写操作
	//slice := []byte{'h','e','l','l','o'}
	//count,err1 := fp.Write(slice)
	count,err1 := fp.Write([]byte("性感老王在线授课"))

	if err1!=nil {
		fmt.Println("写入文件失败")
		return
	}else {
		fmt.Println(count)
	}

	defer fp.Close()
}


func w3(){
	//写文件3
	fp,err := os.Create("./a3.txt")
	if err!=nil{
		//文件创建失败
		/*
		1.路径不存在
		2.文件权限
		3.程序打开文件上限
		 */
		fmt.Println("文件创建失败")
		return
	}

	//写操作
	//获取光标流位置'
	//获取文件起始到结尾有多少个字符
	//count,_:=fp.Seek(0,os.SEEK_END)
	count,_:=fp.Seek(0,io.SeekEnd)
	fmt.Println(count)
	//指定位置写入
	fp.WriteAt([]byte("hello world"),count)
	fp.WriteAt([]byte("hahaha"),0)
	fp.WriteAt([]byte("秀儿"),19)

	defer fp.Close()
}


func main() { 
	log.Println("这是一条优雅的日志。")
    v := "优雅的"
    log.Printf("这是一个%s日志\n", v)
    //fatal系列函数会在写入日志信息后调用os.Exit(1)。Panic系列函数会在写入日志信息后panic。
    log.Fatalln("这是一天会触发fatal的日志") 
    log.Panicln("这是一个会触发panic的日志。") //执行后会自动触发一个异常




w1();
w2();
w3();
rw();
    //打开文件
    fp, err := os.Open("./a.txt")
    if err != nil {
        fmt.Println("err=", err)
        return
    }

    buf := make([]byte, 1024*2) //2k大小
    //n代表从文件读取内容的长度
    n, err1 := fp.Read(buf)
    if err1 != nil && err1 != io.EOF {
        fmt.Println("err1=", err1)
        return
    }
    fmt.Println("buf=", string(buf[:n]))

    //关闭文件
    defer fp.Close()
	}