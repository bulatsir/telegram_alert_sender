package main

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"log"
	"net/http"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	//    "fmt"
)

var Config Conf

type Conf struct {
	Key string `yaml:"key"`
}

type (

	// Timestamp is a helper for (un)marhalling time
	Timestamp time.Time

	// HookMessage is the message we receive from Alertmanager
	HookMessage struct {
		// Alert is a single alert.
		Alerts []struct {
			Labels      map[string]string `json:"labels"`
			Annotations struct {
				Description string `json:"description"`
				Summary     string `json:"summary"`
			}

			StartsAt string `json:"startsAt,omitempty"`
			EndsAt   string `json:"EndsAt,omitempty"`
		}
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
		Status            string            `json:"status"`
		Receiver          string            `json:"receiver"`
		GroupLabels       map[string]string `json:"groupLabels"`
		CommonLabels      map[string]string `json:"commonLabels"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
	}
)

func main() {

	//read configuration

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/alerts", alertsHandler)
	log.Fatal(http.ListenAndServe(":9270", nil))
}

//send message to telegram
func sendMessage(s string, reply string) {
	bot, err := tgbotapi.NewBotAPI(Config.Key)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("send message to %s", s)
	chatId, _ := strconv.ParseInt(s, 10, 64)
	msg := tgbotapi.NewMessage(chatId, reply)

	bot.Send(msg)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func alertsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	default:
		http.Error(w, "unsupported HTTP method", 400)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := enc.Encode("OK"); err != nil {
		log.Printf("error encoding messages: %v", err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var m HookMessage
	if err := dec.Decode(&m); err != nil {
		log.Printf("error decoding message: %v", err)
		http.Error(w, "invalid request body", 400)
		return
	}

	status := m.Status
	alarm_name := m.CommonLabels["alertname"]
	//    alarm_list := m.CommonLabels

	log.Printf("Alarm received: " + alarm_name + "[" + status + "]")

	var chatId string
	var cluster string
	var namespace string
	var pod string
	var phase string
	var details string

	//received alert list handling
	for _, value := range m.Alerts {
		summary := value.Annotations.Summary
		description := value.Annotations.Description
		for k, v := range value.Labels {
			phase = ""
			if k == "cluster" {
				cluster = v
			}
			if k == "namespace" {
				namespace = v
			}
			if k == "pod" {
				pod = v
			}
			if k == "phase" {
				phase = v
			}
			if k == "label_chat_id" {
				chatId = "-" + v
			}

		}

		if len(phase) > 0 {
			details = "cluster: " + cluster + "\n" + "namespace: " + namespace + "\n" + "pod: " + pod + "\n" + "phase: " + phase
		} else {
			details = "cluster: " + cluster + "\n" + "namespace: " + namespace + "\n" + "pod: " + pod
		}

		alertMessage := summary + "\n\n" + details + "\n" + "status: " + description

		if chatId != "-0" {
			log.Printf("received chatId: " + chatId)
			sendMessage(chatId, alertMessage)
		}

	}

	return
}
