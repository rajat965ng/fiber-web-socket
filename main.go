package main

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))
		log.Println(c.Params("id"))
		log.Println(c.Query("v"))
		log.Println(c.Cookies("session"))

		var (
			messageType int
			msg         []byte
			err         error
		)

		for {
			if messageType, msg, err = c.ReadMessage(); err != nil {
				log.Println("read error :", err)
				break
			}
			log.Printf("receive: %s", msg)

			response := string(msg)

			if err = c.WriteMessage(messageType, []byte(fmt.Sprintf("Server: %s", response))); err != nil {
				log.Println("Write error: ", err)
				break
			}
		}

	}))

	app.Listen(":3000")
}
