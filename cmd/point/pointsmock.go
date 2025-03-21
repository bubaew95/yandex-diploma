package main

import (
	"encoding/json"
	"fmt"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/systemdto"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const RateLimit = 2

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

			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(fmt.Sprintf("No more than %d requests per minute allowed", RateLimit)))
			return
		}

		time.Sleep(1 * time.Second)

		response := systemdto.CalculationSystem{
			Order:   num,
			Status:  statuses[rand.Intn(len(statuses))],
			Accrual: 729.98,
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
