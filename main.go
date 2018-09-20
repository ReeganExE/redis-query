package main

import (
	"flag"
	"log"
	"regexp"

	"github.com/go-redis/redis"
)

type param struct {
	addr    string
	pattern string
	key     string
}

var (
	p param
)

func main() {
	flag.StringVar(&p.addr, "addr", "127.0.0.1:6379", "Redis address")
	flag.StringVar(&p.pattern, "pattern", "", "A pattern to extract the value")
	flag.StringVar(&p.key, "key", "", "Key name")

	flag.Parse()

	if p.pattern == "" {
		flag.Usage()
	}

	if p.key == "" {
		flag.Usage()
	}

	client := newClient(p.addr)
	hash, e := client.HGetAll(p.key).Result()

	if e != nil {
		log.Fatal(e)
	}

	pattern := regexp.MustCompile(p.pattern)
	for _, val := range hash {
		if pattern.MatchString(val) {
			println(pattern.FindStringSubmatch(val)[1])
		}
	}
}

func newClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	return client
}
