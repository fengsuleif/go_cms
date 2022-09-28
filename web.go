/*
go run web.go
go build web.go
go build -ldflags "-s -w -H=windowsgui" web.go

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

go env -w GO111MODULE=auto
解除ssl验证
git config --global http.sslVerify "false"

go get  url
*/

package main

import (
	"fmt"
	"net/http"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])
}

func main() {
/*
	array:=[5] int{1,2,3,4,5}
	for i:=0;i<len(array);i++{
	fmt.Println(array[i])
	}
	
	for i,v:=range array{
	fmt.Println(i,v)
	}
	//直接创建切片
	//方法1 创建元素个数为9的切片，初始值是0
	otherslice1:=make([] int ,9)
	//方法2  创建元素个数为9的切片，预留10个元素的空间
	otherslice2:=make([] int ,9,10)
	 
	//根据数组创建切片
	var myslice [] int =array[:3]
	
	for _,v:=range myslice{
	fmt.Println(v)
	}
	s1:=[] int{9,8,7,6,5,4,3,2,1}
	s2:=[] int{0,1,2,3}
	//s3:=append(s1,s2...)
	//复制s2中的4个元素到s1的前四个元素位置进行替换
	copy(s1,s2)
	fmt.Println(len(otherslice1))
	fmt.Println(cap(otherslice2))
	//fmt.Println(s3)
	fmt.Println(s1)
	
	*/
	//map的使用方法
	
	var testMap map[string]int
	testMap = map[string]int{
	  "one": 1,
	  "two": 2,
	  "three": 3,
	}

	k := "two"
	v, ok := testMap[k]
	if ok {
	  fmt.Printf("The element of key %q: %d\n", k, v)
	} else {
	  fmt.Println("Not found!")
	}
	testMap["fort"]=8
	delete(testMap,"three")
	fmt.Println(testMap)
	
	cities1:=make(map[string] string,10)
    cities1["no1"]="feng"
	cities1["no2"]="su"
	cities1["no3"]="lei"
	fmt.Println("map 使用方法1")
	fmt.Println(cities1)
	
	cities2:=make(map[string] string)
    cities2["no1"]="feng"
	cities2["no2"]="su"
	cities2["no3"]="lei"
	fmt.Println("map 使用方法2")
	fmt.Println(cities2)
	
	cities3:=map[string]string{
	"no1":"feng",
	"no2":"su",
	"no3":"lei",
	}
	fmt.Println("map 使用方法3")
	fmt.Println(cities3)
	
	
	student:=make(map[string]map[string]string)
	student["stu1"]=make(map[string]string,3)
	student["stu1"]["name"]="feng"
	student["stu1"]["age"]="18"
	student["stu1"]["address"]="address"
	fmt.Println("map嵌套用例")
	fmt.Println(student)
	
	
	
	for k,v:=range cities1{
	fmt.Println("k:"+k+",v:"+v)
	}
	
	
	
	
	
	
	
	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8080", nil)
	
}

