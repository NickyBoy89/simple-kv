package main

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	ExternalPort = 3000
	KeyName      = "key"
	ValueName    = "value"
)

var (
	ErrNoKeyProvided   = errors.New("error: No key provided")
	ErrNoValueProvided = errors.New("error: No value provided")
	ErrMultipleKeys    = errors.New("error: Multiple keys provided")
	ErrMultipleValues  = errors.New("error: Multiple values provided")
	ErrNoSuchKey       = errors.New("error: No such key exists")
)

var MainDatabase Database = NewKvDatabase("data")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var rawkey []string
			var rawvalue []string
			var exists bool
			if rawkey, exists = r.Form[KeyName]; !exists {
				http.Error(w, ErrNoKeyProvided.Error(), http.StatusBadRequest)
				return
			}
			if rawvalue, exists = r.Form[ValueName]; !exists {
				http.Error(w, ErrNoValueProvided.Error(), http.StatusBadRequest)
				return
			}

			if len(rawkey) > 1 {
				http.Error(w, ErrMultipleKeys.Error(), http.StatusBadRequest)
				return
			}
			if len(rawvalue) > 1 {
				http.Error(w, ErrMultipleValues.Error(), http.StatusBadRequest)
				return
			}

			log.Infof("Inserting key: %s, value: %s", rawkey[0], rawvalue[0])

			written := MainDatabase.Put(rawkey[0], rawvalue[0])
			_ = written

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
			return
		case http.MethodDelete:
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var rawkey []string
			var exists bool
			if rawkey, exists = r.Form[KeyName]; !exists {
				http.Error(w, ErrNoKeyProvided.Error(), http.StatusBadRequest)
				return
			}

			if len(rawkey) > 1 {
				http.Error(w, ErrMultipleKeys.Error(), http.StatusBadRequest)
				return
			}

			log.Infof("Deleting key: %s", rawkey[0])

			deleted := MainDatabase.Delete(rawkey[0])
			if !deleted {
				http.Error(w, ErrNoSuchKey.Error(), http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
			return
		case http.MethodGet:
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var rawkey []string
			var exists bool
			if rawkey, exists = r.Form[KeyName]; !exists {
				http.Error(w, ErrNoKeyProvided.Error(), http.StatusBadRequest)
				return
			}

			if len(rawkey) > 1 {
				http.Error(w, ErrMultipleKeys.Error(), http.StatusBadRequest)
				return
			}

			log.Infof("Getting value at key: %s", rawkey[0])

			result, contains := MainDatabase.Get(rawkey[0])

			if !contains {
				http.Error(w, ErrNoSuchKey.Error(), http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, result)
			return
		}
	})

	log.Infof("Started listening for requests on port %d", ExternalPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", ExternalPort), nil); err != nil {
		log.Fatal(err)
	}
}
