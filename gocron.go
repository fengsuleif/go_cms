
package main

import (
  "fmt"
  "sync"
  "time"
	"net/http"
	"io/ioutil"
  "github.com/roylee0704/gron"
)

func main() {
  var wg sync.WaitGroup
  wg.Add(1)

  c := gron.New()
  c.AddFunc(gron.Every(5*time.Second), func() {
    fmt.Println("runs every 5 seconds.")
	httpget("http://www.sanxiatansuo.com/jsdata/key_redis.php?index=goviews7")
  })
  c.Start()

  wg.Wait()
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