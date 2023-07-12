package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var NickName []byte

func printError(err error) int {
	if nil != err {
		fmt.Println(err)
		return -1
	}

	return 0
}

func firstSendMessages(conn net.Conn) {
	message := "노찬영"
	_, err := conn.Write([]byte(message))
	printError(err)
}

func getRecvCurrentTime() string {
	currentTime := time.Now()
	hour := currentTime.Hour()
	min := currentTime.Minute()

	meridiem := ""

	if hour > 12 {
		meridiem = "오후"
		hour = hour - 12
	} else {
		meridiem = "오전"
	}

	return fmt.Sprintf("%s %02d:%02d", meridiem, hour, min)
}

func getCurrentTimeChatCreated() string {
	currentTime := time.Now()

	year := currentTime.Year()
	hour := currentTime.Hour()
	min := currentTime.Minute()
	weekday := currentTime.Weekday().String()

	// 영어로 받아지는 weekday를 한글말로 치환해주는 case문 추가 필요. ex) Wednesday -> 수요일
	// case
	return fmt.Sprintf("%d년 %d월 %d일 %s", year, hour, min, weekday)
}

func receiveMessages(conn net.Conn) {
	recvBuf := make([]byte, 4096)
	for {
		currentTime := getRecvCurrentTime()
		recvByte, err := conn.Read(recvBuf)
		ret := printError(err)
		if ret != 0 {
			//if io.EOF == err {
			//	fmt.Println("Read null")
			return
			//}
		}

		if recvByte > 0 {
			data := recvBuf[:recvByte]
			fmt.Println("(" + currentTime + ")" + " " + string(NickName) + " " + ":" + " " + string(data))
		}

	}
}

func sendMessages(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// 키보드 입력을 받아서 메시지 전송
		message := scanner.Text()
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Printf("Failed to send message to client %s: %v", conn.RemoteAddr(), err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from stdin: %v", err)
	}
}

func setClientNickName(conn net.Conn) {
	_, err := conn.Read(NickName)
	printError(err)
}

func connectHandler(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected: %s", clientAddr)

	firstSendMessages(conn)
	setClientNickName(conn)

	currentTime := getCurrentTimeChatCreated()

	fmt.Println(currentTime)
	fmt.Println(string(NickName) + "님과 1:1 채팅이 연결되었습니다.")

	go receiveMessages(conn)

	sendMessages(conn)
}

func main() {
	NickName = make([]byte, 200)

	socket, err := net.Listen("tcp", "182.172.222.42:9080")
	ret := printError(err)
	if ret != 0 {
		fmt.Println("server listen error.")
		return
	}

	defer socket.Close()
	for {
		conn, err := socket.Accept()
		printError(err)

		defer conn.Close()
		go connectHandler(conn)
	}

}
