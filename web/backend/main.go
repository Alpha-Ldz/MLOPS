package main

import (
	"fmt"
	"net/http"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"encoding/json"
)

var dict = make(map[string]string)

func setupConsumer(topic string) {
	
		c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{topic}, nil)

	for {
		m, err := c.ReadMessage(-1)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("message at offset %s = %s", string(m.Key), string(m.Value))
		dict[string(m.Key)] = string(m.Value)
	}
}

func main() {
	
	go setupConsumer("rapport-topic2")

	http.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		json, err := json.Marshal(dict)
		if err != nil {
			fmt.Println(err)
		}
		message := fmt.Sprintf("%s", json)
		fmt.Fprintf(w, message)
	})

	http.ListenAndServe(":3333", nil)
}

