package main

import (
	"math/rand"
	"time"
	"os"
	"bufio"
	"strings"
	"math"
	"strconv"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"context"
	"encoding/json"
	"fmt"
)

var writer *kafka.Writer

var writer2 *kafka.Writer

func Configure(kafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}

func Configure2(kafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	w = kafka.NewWriter(config)
	writer2 = w
	return w, nil
}

func Push(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return writer.WriteMessages(parent, message)
}

func Push2(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return writer2.WriteMessages(parent, message)
}

const (
	Bike = iota
	Scooter
	Motorcycle
)

type Engine struct {
	id        int
	user	  int
	vType     int
	Longitude float64
	Latitude  float64
	Status    string
	DestLong  float64
	DestLat   float64
	Speed 	  float64
}

type Rapport struct {
	Id 	  int
	User int
	VType int
	Longitude float64
	Latitude float64
	Status string
}

type Rapport2 struct {
	Id 	  int
	Longitude float64
	Latitude float64
}

type Tuple struct {
	Longitude float64
	Latitude  float64
}

func readLines(path string) ([]Tuple, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rtn := make([]Tuple, 0)

	for _, line := range lines {
		s := strings.Split(line, ",")
		a, _ := strconv.ParseFloat(s[0], 64)
		b, _ := strconv.ParseFloat(s[1], 64)
		tmp := Tuple{Longitude: a, Latitude: b}
		rtn = append(rtn, tmp)
	}

	return rtn, scanner.Err()
}

func dir(lat1, lon1, lat2, lon2 float64) (float64, float64) {
	var pLon float64
	var pLat float64

	if lon2 > lon1 {
		pLon = 1
	} else {
		pLon = -1
	}

	if lat2 > lat1 {
		pLat = 1
	} else {
		pLat = -1
	}

	return pLon, pLat
}

func direction(lat1, lon1, lat2, lon2 float64) (float64, float64) {
	dLon := math.Abs(lon2 - lon1)
	dLat := math.Abs(lat2 - lat1)


	pLon := dLon / (dLon + dLat)
	pLat := dLat / (dLon + dLat)


	return pLon, pLat
}

func (e *Engine) Update(locations []Tuple) {
	if e.Status == "moving" {

		pLon, pLat := direction(e.Latitude, e.Longitude, e.DestLat, e.DestLong)

		dLon, dLat := dir(e.Latitude, e.Longitude, e.DestLat, e.DestLong)

		tmpLon := e.Longitude + (pLon * e.Speed * dLon)
		tmpLat := e.Latitude + (pLat * e.Speed * dLat)

		aLon, aLat := dir(tmpLat, tmpLon, e.DestLat, e.DestLong)

		if aLon != dLon || aLat != dLat {
			e.Longitude = e.DestLong
			e.Latitude = e.DestLat

			e.Status = "waiting"
		} else {
			e.Longitude = tmpLon
			e.Latitude = tmpLat
		}


	} else if e.Status == "waiting" {
		if rand.Intn(3) == 0 {
			e.Status = "moving"
			loc := locations[rand.Intn(len(locations))]

			for loc.Longitude == e.Longitude && loc.Latitude == e.Latitude {
				loc = locations[rand.Intn(len(locations))]
			}

			e.DestLong = loc.Longitude
			e.DestLat = loc.Latitude

		}
	}

	//push a rapport to kafka
	r := Rapport{
		Id: e.id,
		User: e.user,
		VType: e.vType,
		Longitude: e.Longitude,
		Latitude: e.Latitude,
		Status: e.Status,
	}

	b, _ := json.Marshal(r)

	fmt.Println(string(b))

	Push(context.Background(), []byte(strconv.Itoa(e.id)), []byte(b))

	a := Rapport2{
		Id: e.id,
		Longitude: e.Longitude,
		Latitude: e.Latitude,
	}

	b, _ = json.Marshal(a)

	if e.vType == Scooter {
		Configure2([]string{"localhost:9092"}, "simulation", "Scooter-topic")
	} else if e.vType == Bike {
		Configure2([]string{"localhost:9092"}, "simulation", "Bike-topic")
	} else if e.vType == Motorcycle {
		Configure2([]string{"localhost:9092"}, "simulation", "Motorcycle-topic")
	}
	Push2(context.Background(), []byte(strconv.Itoa(e.id)), []byte(b))
}

func NewBike(locations []Tuple) *Engine {
	var e Engine

	loc := locations[rand.Intn(len(locations))]


	e.id = rand.Intn(1000)
	e.vType = Bike
	e.Status = "waiting"
	e.Longitude = loc.Longitude
	e.Latitude = loc.Latitude
	e.Speed = 0.0001
	e.DestLong = 0
	e.DestLat = 0
	e.user = rand.Intn(100)

	return &e
}

func NewScooter(locations []Tuple) *Engine {
	var e Engine

	loc := locations[rand.Intn(len(locations))]

	e.id = rand.Intn(1000)
	e.vType = Scooter
	e.Status = "waiting"
	e.Longitude = loc.Longitude
	e.Latitude = loc.Latitude
	e.Speed = 0.0002
	e.DestLong = 0
	e.DestLat = 0
	e.user = rand.Intn(100)

	return &e
}

func NewMotorcycle(locations []Tuple) *Engine {
	var e Engine

	loc := locations[rand.Intn(len(locations))]

	e.id = rand.Intn(1000)
	e.vType = Motorcycle
	e.Status = "waiting"
	e.Longitude = loc.Longitude
	e.Latitude = loc.Latitude
	e.Speed = 0.0003
	e.DestLong = 0
	e.DestLat = 0
	e.user = rand.Intn(100)

	return &e
}

func engineLoop(locations []Tuple, e *Engine) {
	for {
		e.Update(locations)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	locations, _ := readLines("dataset/destinations")


	// config kafka
	Configure([]string{"localhost:9092"}, "simulation", "rapport-topic2")


	bike := NewBike(locations)
	scooter := NewScooter(locations)
	motorcycle := NewMotorcycle(locations)

	go engineLoop(locations, bike)
	go engineLoop(locations, scooter)
	go engineLoop(locations, motorcycle)

	// while the user do not write quit in the console the program will continue to run
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == "quit" {
			break
		}
	}
}
