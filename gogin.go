/*
go mod init gin

go mod edit -require github.com/gin-gonic/gin@latest
*/

package main
import (
   // "fmt"
	//"bufio"
	"io/ioutil"
    "os"
	"log"
	"github.com/gin-gonic/gin"
)

 

 

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Geektutu")
	})
	
	r.GET("/1", func(c *gin.Context) {
		file, _ := os.Open("a1.txt") 
		defer file.Close() 
		// ReadAll接收一个io.Reader的参数 返回字节切片
		bytes, _ := ioutil.ReadAll(file)
		c.String(200,string(bytes) )
	})
	
	log.Println("Run...")
	r.Run(":8200") // Default listen and serve on 0.0.0.0:8080
}