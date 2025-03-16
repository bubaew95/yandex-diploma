package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Response struct {
	Order   int64  `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

const RateLimit = 5

var requestCounts = 0

var statuses = []string{
	"REGISTERED", "INVALID", "PROCESSING", "PROCESSED",
}

func main() {
	route := chi.NewRouter()

	route.Get("/api/orders/{number}", func(w http.ResponseWriter, r *http.Request) {
		number := chi.URLParam(r, "number")

		reg := regexp.MustCompile("[^0-9]+")
		num := reg.ReplaceAllString(number, "")

		orderNum, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			log.Println(err)
		}

		if requestCounts > RateLimit {
			retry := RateLimit * time.Second
			once := sync.Once{}

			once.Do(func() {
				go func(retry time.Duration) {
					time.Sleep(retry)
					requestCounts = 0
					fmt.Println("requestCounts empty")
				}(retry)
			})

			w.Header().Set("Retry-After", string(retry))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(fmt.Sprintf("No more than %d requests per minute allowed", RateLimit)))
			return
		}

		time.Sleep(1 * time.Second)

		response := Response{
			Order:   orderNum,
			Status:  statuses[rand.Intn(len(statuses))],
			Accrual: rand.Intn(1000),
		}

		requestCounts++
		if response.Accrual < 40 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
		}
	})

	fmt.Println("Running on port 8082")
	if err := http.ListenAndServe(":8082", route); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
