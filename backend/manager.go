package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"github.com/valyala/fasthttp"
)

// Database need to implemented
var Database = make(map[string]string)
var logger = log.New(os.Stdout, "[DEBUG]", log.Ltime)

func init() {
	logfile, _ := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logger.SetOutput(logfile)
	for i := 0; i < 5000; i++ {
		// user1:amazingvpn
		sha_bytes := sha256.Sum256([]byte("amazingvpn"))
		Database["user"+strconv.Itoa(i)] = hex.EncodeToString(sha_bytes[:])
	}
}

func gen_qrcode(dataString string) []byte {
	png, err := qrcode.Encode(dataString, qrcode.Medium, 512)
	if err != nil {
		logger.Println(err)
	}
	return png
}

func get_peer_token(machine_index int, index int) []byte {
	machine := "vpn" + strconv.Itoa(machine_index)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://" + machine + ":8080/peer" + strconv.Itoa(index))

	resp := fasthttp.AcquireResponse()
	cli := &fasthttp.Client{}
	cli.Do(req, resp)

	return gen_qrcode(string(resp.Body()))
}

func select_entry(account string) []byte {
	regnum, err := strconv.Atoi(os.Getenv("REGNUM"))
	if err != nil {
		logger.Println(err)
	}
	num := uint32(regnum)

	sha_bytes := sha256.Sum256([]byte(account))
	val := binary.LittleEndian.Uint32(sha_bytes[:])

	machine, index := int(val%num), int(val%105+1)

	return get_peer_token(machine, index)
}

func main() {
	app := fiber.New()

	app.Post("/Login", func(c *fiber.Ctx) error {
		account, raw_pw := c.FormValue("user"), c.FormValue("passwd")
		sha_bytes := sha256.Sum256([]byte(raw_pw))
		hashed_pw := hex.EncodeToString(sha_bytes[:])

		if Database[account] == hashed_pw {
			tok := select_entry(account)
			return c.Send(tok)
		} else {
			logger.Println("Wrong password: ", hashed_pw, " , Correct hash: ", Database[account])
			return c.SendString("Permission denied!")
		}
	})

	app.Listen(":8080")
}
