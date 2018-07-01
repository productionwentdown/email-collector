package main // import "github.com/productionwentdown/email-collector"

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	// modified from @badoux's checkmail
	"github.com/productionwentdown/email-collector/checkmail"
)

func main() {

	csvMutex := &sync.Mutex{}
	csvFile, err := os.OpenFile("list.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		w.Header().Add("Location", "/subscribed")
		w.WriteHeader(303)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
