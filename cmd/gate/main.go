package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/kaneshin/gate"
	"github.com/kaneshin/gate/cmd/internal"
)

type Data struct {
	Target string
	Text   string `json:"text"`
}

var config = gate.NewConfig().WithHTTPClient(http.DefaultClient)

func postToSlackIncoming(data Data) error {
	svc := gate.NewSlackIncomingService(config).WithBaseURL(data.Target)
	_, err := svc.PostTextData(gate.TextData{
		Text: data.Text,
	})
	if err != nil {
		return err
	}
	return nil
}

func postToLINENotify(data Data) error {
	return errors.New("no implementation")
}

func post(name, text string) string {
	v, ok := internal.Config.Targets.Load(name)
	if !ok {
		return fmt.Sprintf("✘ %s not found in your config\n", name)
	}
	data := Data{
		Target: v.(string),
		Text:   text,
	}
	var err error
	if internal.Config.Slack.IsIncoming(name) {
		err = postToSlackIncoming(data)
	}
	if internal.Config.LINE.IsNotify(name) {
		err = postToLINENotify(data)
	}
	if err != nil {
		return fmt.Sprintf("✘ %s %v\n", name, err)
	}
	return fmt.Sprintf("✔ %s\n", name)
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
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, name := range r.Form["target"] {
		wg.Add(1)
		name := name
		go func(mu *sync.Mutex) {
			defer wg.Done()
			mu.Lock()
			fmt.Fprint(w, post(name, text))
			mu.Unlock()
		}(&mu)
	}
	wg.Wait()
}

func main() {
	err := internal.Load()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handler)

	port := fmt.Sprintf(":%d", internal.Config.Env.Port)
	fmt.Printf("Listening for HTTP on %s%s\n", internal.Config.Env.Host, port)
	log.Fatal(http.ListenAndServe(port, http.DefaultServeMux))
}
