package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
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
