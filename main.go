package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	handleError(err)
	
	for {
		conn, err := ln.Accept()
		handleError(err)
		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, conn, size)
		handleError(err)
		file := buf.Bytes()
		fmt.Println(file)
		fmt.Printf("\n\nreceived %d bytes from network\n\n", n)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
	return
}

func sendLargeFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", ":3000")
	handleError(err)
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes over network\n\n", n)
	return nil
}

func send(t time.Duration, size int) {
	time.Sleep(t * time.Second)
	sendLargeFile(size)
}

func main() {
	go send(2, 10000)
	fs := &FileServer {}
	fs.start()

}