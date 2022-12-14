package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/segmentio/kafka-go"
    "context"
    "time"
)

// Function that return a string
func read(conn *websocket.Conn){
    topic := "rapport-topic"
    partition := 0


    fmt.Println("Starting consumer")

    for {
    con, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
    if err != nil {
        log.Fatal("failed to dial leader:", err)
    }

    con.SetReadDeadline(time.Now().Add(10*time.Second))
        batch := con.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

        b := make([]byte, 10e3) // 10KB max per message
        for {
            n, err := batch.Read(b)
            if err != nil {
                continue
            }
            message := string(b[:n])
            fmt.Println(message)
            conn.WriteMessage(1, []byte(message)) 
        }
    }
}


func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool { return true },
    }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }


    read(conn)
}


func setupRoutes() {
    http.HandleFunc("/", homePage)
    //    http.HandleFunc("/ws", wsEndpoint)

    http.HandleFunc("/ws", wsEndpoint)
}

func main() {
    fmt.Println("Hello World")
    setupRoutes()
    log.Fatal(http.ListenAndServe(":8080", nil))
}
