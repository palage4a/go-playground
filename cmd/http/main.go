package main

import (
	"fmt"
	"log"
	"net/http"
	// "encoding/json"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// body := r.Body
		// defer body.Close()

		// b, err := json.Marshal(r)
		// if err != nil {
		// 	http.Error(w, "json unmarshal error", 400)
		// }

		log.Printf("%q\n", r)
		fmt.Fprintf(w, "Pong, %q", r)
	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}
