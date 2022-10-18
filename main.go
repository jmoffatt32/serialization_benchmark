package main

import (
	"bufio"
	"fmt"
	"math"
	jsonmethod "msg-test/pkg/json_method"
	"net"
	"os"
	"time"
)

type Method interface {
	enc() ([]byte, error)
	dec() error
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func check_detailed(err error, spec string) {
	if err != nil {
		fmt.Println(spec)
		fmt.Println(err)
		return
	}
}

var DELIM byte = 0
var TIME_DELIM byte = 255

func startServer() net.Conn {

	// Start server
	l, err := net.Listen("tcp", ":6969")
	check(err)

	// Listen for connection from 'Recieve'
	c, err := l.Accept()
	check(err)

	netData, err := bufio.NewReader(c).ReadString('\n')
	check(err)
	if netData == "READY\n" {
		os.Stdout.Write([]byte(jsonmethod.NAME + " Benchmark Test:\n"))
	}

	return c

}

func startClient() net.Conn {

	// Connect to 'Send's server
	c, err := net.Dial("tcp", "localhost:6969")
	check(err)

	c.Write([]byte("READY\n"))

	return c
}

func sendMessage(c net.Conn, p []int64) {

	// Encode struct
	b, err := jsonmethod.Encode(p)
	check_detailed(err, "Error in Encode:")

	// Add delimiter
	b = append(b, DELIM)

	// Write bytes to connection
	c.Write(b)
}

func receiveMessage(c net.Conn) []byte {

	// Read bytes from buffer
	b, err := bufio.NewReader(c).ReadBytes(DELIM)
	check_detailed(err, "Error Reading Bytes")

	if b[0] == DELIM {
		return nil
	}

	// Remove Delimiter
	b = b[:len(b)-1]

	// Decode
	var p []int64
	err = jsonmethod.Decode(b, &p)
	check_detailed(err, "Error Decoding")

	// Return timestamp
	end := time.Now()
	t, err := end.MarshalBinary()
	check(err)
	return t
}

func receiveTime(c net.Conn) time.Time {

	t := make([]byte, 15)
	_, err := bufio.NewReader(c).Read(t)
	check(err)

	var end time.Time
	err = end.UnmarshalBinary(t)
	check(err)

	return end
}

func makePayload(size int) []int64 {

	d := make([]int64, size)
	for i := 0; i < len(d); i++ {
		d[i] = int64(i)
	}
	return d
}

func SendProcess() {
	// Function to be ran by sending process

	c := startServer()
	// Once connection is established, take timestamp of `START_TIME`

	for i := 3; i < 9; i++ {
		var size = int(math.Pow(10, float64(i)))
		var p = makePayload(size)

		start := time.Now()

		sendMessage(c, p)

		end := receiveTime(c)

		var size_Mb = float32(size) / 1000000
		fmt.Printf("Message size-%.3f Megabytes:\n", size_Mb)
		fmt.Println(end.Sub(start))

	}

	STOP := make([]byte, 1)
	STOP[0] = DELIM
	c.Write(STOP)

}

func ReceiveProcess() {
	// Function to be ran by recieving process

	c := startClient()

	// Listen for incoming messages
	for {

		// Receives message, demarshals, and takes timestamp of finishing time
		t := receiveMessage(c)

		if t == nil {
			return
		}

		// Write timestamp back to server to compute overall time
		c.Write(t)

	}
	// Once all bytes are read, take timestamp of `END_TIME`

	// Send end-time back to server

}

func main() {

	go SendProcess()

	time.Sleep(1 * time.Second)

	ReceiveProcess()

}
