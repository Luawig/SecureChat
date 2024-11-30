package server

import (
	"SecureChat/common"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

func Start(username, password string) {
	// 加载服务器证书和私钥
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal(err)
	}

	// 加载客户端证书的CA
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	listener, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("服务器启动，等待连接...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, password, username)
	}
}

func handleConnection(conn net.Conn, password, username string) {
	encryptionKey := []byte(password) // 使用密码作为加密密钥
	go common.ReceiveMessages(conn, encryptionKey, username)
	common.SendMessages(conn, encryptionKey, username)
}
