package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)


func main() {
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		// How long to wait for command to complete, in milliseconds
		Timeout: 50000,

		// MaxConcurrent is how many commands of the same type
		// can run at the same time
		MaxConcurrentRequests: 300,

		// VolumeThreshold is the minimum number of requests
		// needed before a circuit can be tripped due to health
		RequestVolumeThreshold: 10,

		// SleepWindow is how long, in milliseconds,
		// to wait after a circuit opens before testing for recovery
		SleepWindow: 1000,

		// ErrorPercentThreshold causes circuits to open once
		// the rolling measure of errors exceeds this percent of requests
		ErrorPercentThreshold: 50,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		resultChan := make(chan string, 1)
		errChan := hystrix.Go("my_command", func() error {
			resp, err := http.Get("http://localhost:6061")
			if err != nil {
				return err
			}
			defer r.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			resultChan <- string(b)

			return nil
		}, nil)


		// Block until we have a result or an error.
		select {
		case result := <-resultChan:
			log.Println("success:", result)
			w.WriteHeader(http.StatusOK)
		case err := <-errChan:
			log.Println("failure:", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	http.ListenAndServe(":6060", nil)
}
