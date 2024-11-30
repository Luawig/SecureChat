package client

import (
	"SecureChat/common"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func Start(username, password string) {
	// 加载客户端证书和私钥
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	// 加载服务器证书的CA
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	conn, err := tls.Dial("tcp", "localhost:8443", config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	encryptionKey := []byte(password) // 使用密码作为加密密钥

	go common.ReceiveMessages(conn, encryptionKey, username)
	common.SendMessages(conn, encryptionKey, username)
}
