package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"github.com/valyala/fasthttp"
)

var logger = log.New(os.Stdout, "[DEBUG]", log.Ltime)

func init() {
	logfile, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logger.SetOutput(logfile)
}

// Genrate the QR code of file
func gen_qrcode(dataString string) []byte {
	png, err := qrcode.Encode(dataString, qrcode.Medium, 512)
	if err != nil {
		logger.Println("QR code generating error:", err)
	}
	return png
}

// Get the peer's conf file, return that file's QR code
func get_peer_token(machine_index int, index int) []byte {
	machine := "vpn" + strconv.Itoa(machine_index)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://" + machine + ":8080/peer" + strconv.Itoa(index))

	resp := fasthttp.AcquireResponse()
	cli := &fasthttp.Client{}
	cli.Do(req, resp)

	return gen_qrcode(string(resp.Body()))
}

// Select a machine and peer to be tunneled, returns a QR code of the peer.conf
func select_entry(account string) []byte {
	// The number of machines in total
	regnum, err := strconv.Atoi(os.Getenv("REGNUM"))
	if err != nil {
		logger.Println("Wrong VPN machine index", err)
	}
	num := uint32(regnum)

	// Select which machine
	sha_bytes := sha256.Sum256([]byte(account))
	val := binary.LittleEndian.Uint32(sha_bytes[:])

	// The selected machine's allowed peers
	_peer, err := strconv.Atoi(os.Getenv("VPN" + strconv.Itoa(int(val%num)) + "_PEER"))
	if err != nil {
		logger.Println("Wrong VPN total peers on", num, err)
	}
	peer := uint32(_peer)

	machine, index := int(val%num), int(val%peer+1)
	return get_peer_token(machine, index)
}

// Try login to SSO, return it could login or not
func auth(account string, passwd string) bool {

	args := struct {
		StudentID string `json:"studentID"`
		Password  string `json:"password"`
	}{
		StudentID: account,
		Password:  passwd,
	}

	url := os.Getenv("SSO_URL")
	if url == "" {
		log.Panicln("The SSO_URL is not setted.")
	}
	req := &fasthttp.Request{}
	req.SetRequestURI(url)
	requestBody, _ := json.Marshal(args)
	req.SetBody(requestBody)
	req.Header.SetMethod("POST")

	resp := &fasthttp.Response{}
	cli := &fasthttp.Client{}

	if err := cli.Do(req, resp); err != nil || len(resp.Body()) == 0 {
		return false
	}
	return true
}

func main() {
	app := fiber.New()

	app.Post("/Login", func(c *fiber.Ctx) error {
		account, raw_pw := c.FormValue("user"), c.FormValue("passwd")

		if auth(account, raw_pw) {
			tok := select_entry(account)
			return c.Send(tok)
		} else {
			return c.SendString("Permission denied!")
		}
	})

	app.Listen(":8080")
}
