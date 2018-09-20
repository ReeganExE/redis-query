package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-redis/redis"
)

type param struct {
	addr    string
	pattern string
	key     string
	port    uint
	verbose bool
}

type HttpResponse struct {
	Value string
}

var (
	p      param
	client *redis.Client
)

func main() {
	flag.StringVar(&p.addr, "addr", "127.0.0.1:6379", "Redis address")
	flag.StringVar(&p.pattern, "pattern", "", "A pattern to extract the Value")
	flag.StringVar(&p.key, "key", "", "Key name")
	flag.UintVar(&p.port, "port", 0, "HTTP port")
	flag.BoolVar(&p.verbose, "verbose", false, "Print raw Value")

	flag.Parse()

	if p.pattern == "" {
		flag.Usage()
	}

	if p.key == "" {
		flag.Usage()
	}

	client = newClient(p.addr)

	if p.port != 0 {
		http.HandleFunc("/", httpHandler)

		fmt.Printf("HTTP server listens on port %d\n", p.port)

		if e := http.ListenAndServe(fmt.Sprintf(":%d", p.port), nil); e != nil {
			log.Fatal(e)
		}

		return
	}

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

func httpHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	key := req.FormValue("key")

	if key == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	value, e := getValue(key)

	if p.verbose {
		fmt.Println(value, e)
	}

	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, e.Error())))
		return
	}

	json.NewEncoder(w).Encode(&HttpResponse{value})
}

func getValue(key string) (string, error) {
	log.Printf("Getting key %s", key)

	hash, e := client.HGetAll(key).Result()

	if e != nil {
		return "", e
	}

	pattern := regexp.MustCompile(p.pattern)
	for _, val := range hash {
		if pattern.MatchString(val) {
			return pattern.FindStringSubmatch(val)[1], nil
		}
	}

	return "", nil
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
