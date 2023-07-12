package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	var host string
	var port string

	fmt.Scanln(&host)
	fmt.Scanln(&port)

	conn, err := net.Dial("tcp", host+":"+port)
	if nil != err {
		log.Println(err)
	}

	go func() {
		data := make([]byte, 4096)

		for {
			n, err := conn.Read(data)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println("Server send : " + string(data[:n]))
			//time.Sleep(time.Duration(3) * time.Second)
		}
	}()

	for {
		var s string
		fmt.Scanln(&s)
		conn.Write([]byte(s))
		//time.Sleep(time.Duration(3) * time.Second)
	}
}
