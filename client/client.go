package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

var adress = "127.0.0.1:8000"

func getPUTInfos(str string) (string, string, error) {
	var fn = strings.TrimPrefix(str, "PUT ")
	var datas, err = readFile(fn)
	return fn, datas, err
}

func getGetInfos(str string) (string, string) {
	var strArgs = strings.TrimPrefix(str, "FN:")
	var args = strings.Split(strArgs, "DATA:")
	return args[0], args[1]
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

	reader := bufio.NewReader(os.Stdin)

	// connect to server
	conn, _ := net.Dial("tcp", adress)

	for {
		//wait for request

		fmt.Print("Command : ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(string(text), "\n")

		if strings.HasPrefix(text, "GET ") {
			fmt.Println("GET")

			var req = text + "\n"
			fmt.Println("GET sent")
			fmt.Fprintf(conn, req)
			message, _ := bufio.NewReader(conn).ReadString('\n')

			fmt.Println("message")

			fn, datas := getGetInfos(message)
			writeFile(fn, []byte(datas))
		}
		if strings.HasPrefix(text, "PUT ") {
			fmt.Println("PUT")

			fn, datas, err := getPUTInfos(text)
			if err != nil {
				fmt.Println(err)
			}
			var res = "PUT FN:" + fn + " DATA:" + datas + "\n"
			fmt.Println("PUT sent")
			fmt.Fprintf(conn, res)
		}
	}
}
