package main // import "github.com/productionwentdown/email-collector"

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	// modified from @badoux's checkmail
	"github.com/productionwentdown/email-collector/checkmail"
)

var filename string
var listen string
var redirect string
var slack string

func main() {
	flag.StringVar(&filename, "file", "list.csv", "file to append records to")
	flag.StringVar(&listen, "listen", ":8080", "address to listen to")
	flag.StringVar(&redirect, "redirect", "/subscribed", "path to redirect to upon success")
	flag.StringVar(&slack, "slack", "", "optional slack webhook url")
	flag.Parse()

	csvMutex := &sync.Mutex{}
	csvFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer csvFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	csvWriter := csv.NewWriter(csvFile)

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			w.Write([]byte("An error occurred. Please try again"))
			return
		}
		email := r.Form.Get("email")
		if err := checkmail.ValidateFormat(email); err != nil {
			log.Println(email, err)
			w.WriteHeader(400)
			w.Write([]byte("Email is not valid"))
			return
		}
		err = checkmail.ValidateHost(email)
		if _, ok := err.(checkmail.SmtpError); !ok && err != nil {
			log.Println(email, err)
			w.WriteHeader(400)
			w.Write([]byte("Email is not valid"))
			return
		}
		log.Println(email, "success")
		csvMutex.Lock()
		csvWriter.Write([]string{email, time.Now().String()})
		csvWriter.Flush()
		csvMutex.Unlock()
		if len(slack) > 4 {
			object := map[string]string{
				"text": "A human (" + email + ") subscribed to the list! Wohoo!",
			}
			payload, err := json.Marshal(object)
			if err != nil {
				log.Println("json.Marshal failed, not supposed to happen!")
				w.WriteHeader(500)
				return
			}
			resp, err := http.Post(slack, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				log.Println(body)
				w.WriteHeader(500)
				return
			}
		}
		w.Header().Add("Location", redirect)
		w.WriteHeader(303)
	})

	log.Fatal(http.ListenAndServe(listen, nil))

}
