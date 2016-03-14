// To run:
// go get github.com/githubnemo/CompileDaemon
// CompileDaemon -command="./gervice"

package main

import "os"
import "fmt"
import "time"
import "strconv"
import redis "gopkg.in/redis.v3"

type Driver interface {
  Start()
  Set()
  Get()
}


type ServiceCacheEntry struct {
  Name  string
  URL   string
}

type Updater struct {
  TTL           time.Time
  State         string
  Driver        Driver
  ServiceCache  []ServiceCacheEntry
}

var RedisClient *redis.Client

func main() {
  fmt.Println("Starting GOSD client")
  name := "service-2"
  url := "http://localhost:8881"

  // start
  redisDB,_ := strconv.Atoi(os.Getenv("RedisDB"))
  RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("RedisAddr"),
        Password: os.Getenv("RedisPassword"), // no password set
        DB:       int64(redisDB),  // use default DB
    })
  _, err := RedisClient.Ping().Result()
  if err != nil {
    panic(err.Error())
  }

  // set
  if err = RedisClient.HSet("gosd", name, url).Err(); err != nil {
    panic(err.Error())
  }

  // get
  val := tryRefreshForNTimes(3)
  for key,value := range val {
    fmt.Println(key,value)
  }

}

func tryRefreshForNTimes(n int) map[string]string {
  for n > 0 {
    val, err := RedisClient.HGetAllMap("gosd").Result()
    if err != nil {
      n--
    } else {
      return val
    }
  }
  return map[string]string{}
}
