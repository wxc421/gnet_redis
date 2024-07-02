package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// data := []byte("+OK\r\n")
	data := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	conn, err := net.DialTimeout("tcp", "localhost:6399", time.Second*30)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	defer conn.Close()
	_, err = conn.Write(data)
	time.Sleep(time.Second * 60)
	if err != nil {
		fmt.Printf("write failed , err : %v\n", err)
	}

}
