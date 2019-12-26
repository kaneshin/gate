package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/k0kubun/pp"
	"github.com/kaneshin/gate/cmd/internal"
	"github.com/pkg/errors"
)

type Data struct {
	Target string
	Text   string `json:"text"`
}

func Post(data Data) error {
	pp.Println(data)
	target, ok := internal.Config.Targets[data.Target]
	if !ok {
		return errors.New(fmt.Sprintf("%s is not defined", data.Target))
	}

	if data.Text == "" {
		return nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var buf = bytes.NewBuffer(b)
	_, err = http.Post(target, "application/json", buf)
	if err != nil {
		return err
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "405 Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	text := r.FormValue("text")
	for _, v := range r.Form["target"] {
		if _, ok := internal.Config.Targets[v]; !ok {
			msg := v + " target is not found in your config"
			if err == nil {
				err = errors.New(msg)
			} else {
				err = errors.Wrap(err, msg)
			}
			continue
		}
		d := Data{
			Target: v,
			Text:   text,
		}
		if err1 := Post(d); err1 != nil {
			if err == nil {
				err = err1
			} else {
				err = errors.Wrap(err, err1.Error())
			}
		}
	}

	fmt.Fprintln(w, "success")
	if err != nil {
		fmt.Fprintf(w, "\n%v\n", err)
	}
}

func main() {
	internal.Load()

	http.HandleFunc("/", handler)

	port := fmt.Sprintf(":%d", internal.Config.Env.Port)
	log.Fatal(http.ListenAndServe(port, http.DefaultServeMux))
}
