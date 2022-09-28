package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
	"encoding/json"
	"strings"
	"sort"
	_  "github.com/mattn/go-sqlite3"
)

const (
	MAX_UPLOAD_SIZE = 1024 * 1024 * 50 //50MB
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	/*
			log.Println("message")
		    // Fatalln is Println() followed by a call to os.Exit(1).
		    //log.Fatalln("fatal message")
		    // Panicln is Println() followed by a call to panic().
		    log.Panicln("panic message")
	*/
}

func main() {
	fmt.Println("执行路径：")
	fmt.Println(filepath.Abs(filepath.Dir(os.Args[0])))
	fmt.Println("当前时间：")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	//使用https
	//server.ListtenAndServeTLS("cert.pem","key.pem");
	//httpget("http://www.sanxiatansuo.com/jsdata/key_redis.php?index=goviews")
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.Handle("/kindeditor/", http.StripPrefix("/kindeditor/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/page", get_singel_contorl)
	http.HandleFunc("/search", search_contorl)
	http.HandleFunc("/add", insert_art_contorl)
	http.HandleFunc("/form", save_form_contorl)
	http.HandleFunc("/view", view_art_contorl)
	http.HandleFunc("/edit", update_art_contorl)
	http.HandleFunc("/update", update_art_contorl)
	http.HandleFunc("/upload",uploadHandler)
	server.ListenAndServe()
}


func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Printf("File is too big")
		return
	}
	file, headers, err := r.FormFile("imgFile")
	if err != nil {
		log.Printf("Error when try to get file: %v", err)
		return
	}
	fileName:=headers.Filename
	var fileExt string = strings.ToLower(fileName[strings.Index(fileName, ".")+1:])//获取上传文件的类型
	fileExt_array:=[]string{"png", "jpg", "jgeg","gif","bmp"} 
	if !in_array_fast(fileExt,fileExt_array){
		log.Printf("图片类型错误") 
	}
	 
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		return
	}
	 
	err = ioutil.WriteFile("./public/upload/"+fileName, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		return
	}
	
	
	var hash map[string]interface{} = make(map[string]interface{})
	fileUrl:="/kindeditor/upload/"+fileName
	if err == nil {
		hash["error"] = 0
		hash["url"] = fileUrl
	} else {
		hash["error"] = 1
		hash["message"] = err.Error()
	}
	fmt.Println(fileUrl)
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	d, _ := json.Marshal(hash) 
	w.Write(d)
}

type User struct {
	Id      int
	Title   template.HTML
	Author  template.HTML
	Content template.HTML
	Time    string
}

type Art_list struct {
	Id    int
	Title string
}

type Post struct {
	Id     int
	Title  string
	Author string
}

type Post_all struct {
	Totalpage int
	Artdata   []Post
}

func Posts(limit, page int) (posts []Post, err error) {
	Db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)

	pagesize := 15
	if page < 1 {
		return
	}
	lnum := (page - 1) * pagesize

	rows, err := Db.Query("SELECT id,title,author FROM  info  order by id desc limit $1 ,$2", lnum, limit)

	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func get_count(sql_str string) (total int) {
	Db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)

	rows, err := Db.Query(sql_str)
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return
		}
	}
	rows.Close()
	return
}

func get_art(id int) (posts User, err error) {
	Db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)
	rows, err := Db.Query("SELECT id,title,author,content,create_at FROM  info where id =$1", id)
	for rows.Next() {
		err = rows.Scan(&posts.Id, &posts.Title, &posts.Author, &posts.Content, &posts.Time)
		if err != nil {
			fmt.Println(err)
		}
	}
	rows.Close()
	return
}

func Posts_search(limit, page int, word, search_type string) (posts []Post, err error) {
	sql_str := "SELECT id,title ,author FROM info where " + search_type + " like '%" + word + "%'"

	Db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)

	pagesize := 15
	if page < 1 {
		return
	}
	lnum := (page - 1) * pagesize
	rows, err := Db.Query(sql_str+" order by id desc limit $1 ,$2", lnum, limit)

	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Author)
		fmt.Println(post)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	data := Post_all{}
	sql_str := "SELECT count(*) FROM info"

	posts, _ := Posts(10, 1)

	data.Totalpage = get_count(sql_str)
	data.Artdata = posts
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, data)
}

func search_contorl(w http.ResponseWriter, r *http.Request) {
	data := Post_all{}
	r.ParseForm()
	word := r.PostFormValue("word")
	search_type := r.PostFormValue("search_type")

	fmt.Println(search_type)
	posts, _ := Posts_search(10, 1, word, search_type)
	sql_str := "SELECT count(*) FROM info where " + search_type + " like '%" + word + "%'"
	data.Totalpage = get_count(sql_str)
	data.Artdata = posts

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, data)
}

func update_art_contorl(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm() 如果使用/enctype="application/x-www-form-urlencoded"
	r.ParseMultipartForm(20480)//enctype="multipart/form-data"
	if r.Method == "POST" {
		id, err := strconv.Atoi(r.PostFormValue("id"))
		if err != nil { 
			fmt.Println(err) 
		} 
		
		content := User{
			id,
			template.HTML(r.PostFormValue("title")),
			template.HTML(r.PostFormValue("author")),
			template.HTML(r.PostFormValue("contentx")),
			r.PostFormValue("create_at"),
		}
		var rid int64=edit(content)
		fmt.Println(rid)
		
		file, _, _ := r.FormFile("file")
		if file!=nil{
		//文件上传start 
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		r.ParseMultipartForm(MAX_UPLOAD_SIZE); 
		
		fmt.Println("##########") 
		fmt.Println(r.FormFile("file"))
		fmt.Println("##########") 
		file, headers, err := r.FormFile("file")
		/**/
		if err != nil {
			log.Printf("Error when try to get file: %v", err)
		}
		
		fileName:=headers.Filename
		var fileExt string = strings.ToLower(fileName[strings.Index(fileName, ".")+1:])
		//获取上传文件的类型
		if !(fileExt == "png" ||fileExt == "jpg" ||fileExt == "jpeg"){
			log.Printf("图片类型错误")
			return
		}
		
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Read file error: %v", err)
		}
		 
		err = ioutil.WriteFile("./public/upload/"+fileName, data, 0666)
		if err != nil {
			log.Printf("Write file error: %v", err)
		}
		//文件上传end
	}

		if rid > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		pageno, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			fmt.Println(err)
		}
		data, _ := get_art(pageno)
		fmt.Println(r.Method)

		fmt.Println(data)
		t, _ := template.ParseFiles("editform.html")
		t.Execute(w, data)
	}
}

func process(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, "Hello World!")
}

func get_singel_contorl(w http.ResponseWriter, r *http.Request) {
	pageno, err := strconv.Atoi(r.URL.Query().Get("num"))
	if err != nil {
		fmt.Println("转换失败")
	}
	posts, _ := Posts(10, pageno)
	t, _ := template.ParseFiles("index_data.html")
	t.Execute(w, posts)
}

func headers(w http.ResponseWriter, r *http.Request) {
	h := r.Header
	r.ParseForm() //解析参数，默认是不会解析的
	fmt.Println("path", r.URL.Path)
	fmt.Println("参数", r.URL.Query().Get("feng"))
	fmt.Println("Method", r.Method)
	fmt.Println("RequestURI", r.RequestURI)
	fmt.Println("head", r.Header.Get("Authorization"))
	fmt.Println("RawQuery", r.URL.RawQuery)
	fmt.Fprintln(w, h)
}

func insert_art_contorl(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("theform.html")
	t.Execute(w, nil)
}

func save_form_contorl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("name", r.PostFormValue("post"))
	fmt.Println("name", r.Form["post"])

	db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)
	//插入数据
	stmt, err := db.Prepare("INSERT INTO info(title, author,content, create_at) values(?,?,?,?)")
	checkErr(err)
	res, err := stmt.Exec(r.PostFormValue("title"), r.PostFormValue("author"), r.PostFormValue("contentx"), time.Now().Format("2006-01-02"))
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func view_art_contorl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.URL.Query().Get("id")

	db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)
	//查询数据
	rows, err := db.Query("SELECT id,title,author,content,create_at FROM info where id=" + id)
	checkErr(err)

	for rows.Next() {
		var id int
		var title string
		var author string
		var content string
		var create_at string
		err = rows.Scan(&id, &title, &author, &content, &create_at)
		checkErr(err)

		user := User{
			id,
			template.HTML(title),
			template.HTML(author),
			template.HTML(content),
			create_at,
		}
		//fmt.Println(user.author)
		t, _ := template.ParseFiles("view.html")
		//w.Header().Set("Content-Type","text/html")
		t.Execute(w, user)
	}
}

//删除数据
func del(id int) (total int64) {
	db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)
	stmt, err := db.Prepare("delete from info where id=?")
	checkErr(err)
	res, err := stmt.Exec(id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	return affect
}

//修改数据
func edit(content User) (total int64) {
	db, err := sql.Open("sqlite3", "./sqlite3/u.db")
	checkErr(err)
	stmt, err := db.Prepare("update  info set title=?,  author=? ,content=?,create_at=? where id=?")
	checkErr(err)
	res, err := stmt.Exec(content.Title, content.Author, content.Content, content.Time, content.Id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	return affect
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//查找字符串是否在字符串数组中
func in_array(target string, str_array []string) bool { 
	for _, element := range str_array{ 
	if target == element{  
            return true  
        }  
    } 
    return false  
} 


//查找字符串是否在字符串数组中 大数据版,需要引入"sort"
func in_array_fast(target string, str_array []string) bool { 
    sort.Strings(str_array) 
    index := sort.SearchStrings(str_array, target) 
    if index < len(str_array) && str_array[index] == target {  
        return true  
    }  
    return false  
} 

func httpget(url string){
	resp, err := http.Get(url) 
	if err != nil { 
	fmt.Println(" handle error2") 
	} 
	defer resp.Body.Close() 
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil { 
	 fmt.Println(" handle error2") 
	} 
	fmt.Println(string(body))
}
