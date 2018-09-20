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
	Value string `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	args   param
	client *redis.Client
)

func main() {
	flag.StringVar(&args.addr, "addr", "127.0.0.1:6379", "Redis address")
	flag.StringVar(&args.pattern, "pattern", "", "A pattern to extract the Value")
	flag.StringVar(&args.key, "key", "", "Key name")
	flag.UintVar(&args.port, "port", 0, "HTTP port")
	flag.BoolVar(&args.verbose, "verbose", false, "Print raw Value")

	flag.Parse()

	client = newClient(args.addr)

	if args.port != 0 {
		// Start a HTTP server on the specified port
		http.HandleFunc("/query", httpHandler)

		fmt.Printf("HTTP server listens on port %d\n", args.port)

		if e := http.ListenAndServe(fmt.Sprintf(":%d", args.port), nil); e != nil {
			log.Fatal(e)
		}

		return
	}

	if args.pattern == "" {
		flag.Usage()
		return
	}

	if args.key == "" {
		flag.Usage()
		return
	}

	value, e := getValue(args.key)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Println(value)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	key := req.FormValue("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ErrorResponse{"Invalid parameter. Use ?key=..."})
		return
	}

	value, e := getValue(key)

	if args.verbose {
		fmt.Println(value, e)
	}

	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ErrorResponse{e.Error()})
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

	pattern := regexp.MustCompile(args.pattern)
	for _, val := range hash {
		if pattern.MatchString(val) {
			groups := pattern.FindStringSubmatch(val)
			if len(groups) == 2 {
				// return this first matched group
				return groups[1], nil
			}

			// we don't know which group to be return,
			// return the whole matched one
			return groups[0], nil
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
