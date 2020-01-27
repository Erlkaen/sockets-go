package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
)

var staticCode = 0

var clients []Client

// Client struct
type Client struct {
	code       int
	connection net.Conn
}

const port = ":8000"

func getPUTInfos(str string) (string, string) {
	var strArgs = strings.TrimPrefix(str, "PUT FN:")
	var args = strings.Split(strArgs, "DATA:")
	return args[0], args[1]
}

func getGetInfos(str string) (string, string, error) {
	var fn = strings.TrimPrefix(str, "GET ")
	fn = strings.Trim(fn, "\n")
	fn = strings.Trim(fn, " ")
	fmt.Println("a.txt" == fn)
	var datas, err = readFile(fn)
	return fn, datas, err
}

func handleConnection(conn net.Conn) error {
	var client = Client{staticCode, conn}
	clients = append(clients, client)
	staticCode++
	fmt.Println("Connection with client ", client.code)
	for {
		message, error := bufio.NewReader(client.connection).ReadString('\n')
		if error == io.EOF {
			fmt.Println("Connection closed")
			return error
		}
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println("Client : ", message)
		message = strings.TrimSuffix(string(message), "\n")

		if strings.HasPrefix(message, "GET ") {
			fn, datas, err := getGetInfos(message)
			if err != nil {
				fmt.Println("error : No such file in directory")
				return err
			}
			fmt.Println("ok")
			var res = "FN:" + fn + " DATA:" + datas + "\n"
			fmt.Fprintf(client.connection, res)
		}
		if strings.HasPrefix(message, "PUT ") {
			fn, datas := getPUTInfos(message)
			var err = writeFile(fn, []byte(datas))
			if err != nil {
				fmt.Fprintf(client.connection, "File could not be writed")
			}
			fmt.Fprintf(client.connection, "File writed")
		}
	}
}

func writeFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, 0644)
	return err
}

func readFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	return string(data), err
}

func main() {
	fmt.Println("Start server...")
	// listen
	ln, _ := net.Listen("tcp", port)

	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}
