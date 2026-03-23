package  transaction

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // Allow all origins for development
    CheckOrigin: func(r *http.Request) bool { return true },
}





func HandleWebSocket(c *gin.Context, rdb *redis.Client) {
    defer func() {
        if r := recover(); r != nil {
            log.Println("PANIC RECOVERED:", r)
        }
    }()

    log.Println("WS handler started")

    userIDAny, exists := c.Get("user_id")
    if !exists {
        log.Println("user_id not found")
        return
    }

    userID, ok := userIDAny.(string)
    if !ok {
        log.Println("user_id is not string")
        return
    }

    log.Println("UserID:", userID)

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    ctx := context.Background()
    key := fmt.Sprintf("user:txs:%s", userID)

    vals, err := rdb.LRange(ctx, key, 0, -1).Result()
    if err != nil {
        log.Println("Redis LRange error:", err)
    }

    payload, _ := json.Marshal(vals)
    if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
        log.Println("Write error:", err)
        return
    }

    channel := fmt.Sprintf("user_updates:%s", userID)
    pubsub := rdb.Subscribe(ctx, channel)
    defer pubsub.Close()

    ch := pubsub.Channel()

    for {
        msg, ok := <-ch
        if !ok {
            log.Println("Channel closed")
            return
        }

        log.Println("Received:", msg.Payload)

        vals, err := rdb.LRange(ctx, key, 0, -1).Result()
        if err != nil {
            log.Println("Redis error:", err)
            continue
        }

        payload, _ := json.Marshal(vals)

        if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
            log.Println("Write error:", err)
            return
        }
    }
}