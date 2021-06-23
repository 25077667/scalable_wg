package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Database need to implemented
var Database = make(map[string]string)

func init_db() {
	for i := 0; i < 5000; i++ {
		// user1:amazingvpn
		sha_bytes := sha256.Sum256([]byte("amazingvpn"))
		Database["user"+strconv.Itoa(i)] = hex.EncodeToString(sha_bytes[:])
	}
}

func get_acc_pw(data []string) (string, string) {
	account, pw := strings.Split(data[0], "="), strings.Split(data[1], "=")
	return account[1], pw[1]
}

func select_server() (string, string) {
	num, _ := strconv.Atoi(os.Getenv("REGNUM"))
	var result string
	for i := 0; i < num; i++ {
		result = os.Getenv("REMOTE" + strconv.Itoa(i))
	}
	// ! FIXME: Token
	return result, "TOKEN"
}

func main() {
	init_db()

	app := fiber.New()

	app.Post("/Login", func(c *fiber.Ctx) error {
		user_data := strings.Split(string(c.Body()), "&")
		if len(user_data) != 2 {
			return c.SendString("Invalid input!")
		}
		account, raw_pw := get_acc_pw(user_data)
		sha_bytes := sha256.Sum256([]byte(raw_pw))
		hashed_pw := hex.EncodeToString(sha_bytes[:])

		if Database[account] == hashed_pw {
			ip, tok := select_server()
			return c.SendString(ip + "\n" + tok)
		} else {
			return c.SendString("Permission denied!")
		}
	})

	app.Listen(":3000")
}
