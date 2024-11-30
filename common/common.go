package common

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// Message 定义消息结构体
type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
}

// 接收消息并显示
func ReceiveMessages(conn net.Conn, encryptionKey []byte, username string) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("连接关闭:", err)
			return
		}
		encryptedMessage := string(buf[:n])

		// 解密消息
		messageJSON, err := DecryptMessage(encryptionKey, encryptedMessage)
		if err != nil {
			log.Println("解密消息失败:", err)
			continue
		}

		var message Message
		err = json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			log.Println("解析消息失败:", err)
			continue
		}

		fmt.Printf("\n[%s] %s: %s\n", message.Timestamp.Format("2006-01-02 15:04:05"), message.Username, message.Content)
		fmt.Print("输入消息: ")

		// 保存加密消息
		SaveMessageToFile(username+".log", encryptedMessage)
	}
}

// 监听命令行输入并发送消息
func SendMessages(conn net.Conn, encryptionKey []byte, username string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("输入消息: ")
		content, _ := reader.ReadString('\n')

		message := Message{
			Timestamp: time.Now(),
			Username:  username,
			Content:   content,
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Println("编码消息失败:", err)
			continue
		}

		// 加密消息
		encryptedMessage, err := EncryptMessage(encryptionKey, messageJSON)
		if err != nil {
			log.Println("加密消息失败:", err)
			continue
		}

		_, err = conn.Write([]byte(encryptedMessage))
		if err != nil {
			log.Println("发送消息失败:", err)
			return
		}

		// 保存加密消息
		SaveMessageToFile(username+".log", encryptedMessage)
	}
}

// 加密消息
func EncryptMessage(key, message []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(message))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], message)

	return hex.EncodeToString(ciphertext), nil
}

// 解密消息
func DecryptMessage(key []byte, encryptedMessage string) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// 保存消息到文件
func SaveMessageToFile(filename, message string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("无法打开文件:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(message + "\n"); err != nil {
		log.Println("无法写入文件:", err)
	}
}

// 查看本地聊天记录
func ViewChatLog(username, password string) {
	encryptionKey := []byte(password) // 使用密码作为加密密钥
	file, err := os.Open(username + ".log")
	if err != nil {
		log.Fatal("无法打开聊天记录文件:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		encryptedMessage := scanner.Text()
		messageJSON, err := DecryptMessage(encryptionKey, encryptedMessage)
		if err != nil {
			log.Println("解密消息失败:", err)
			continue
		}

		var message Message
		err = json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			log.Println("解析消息失败:", err)
			continue
		}

		fmt.Printf("[%s] %s: %s\n", message.Timestamp.Format(time.RFC3339), message.Username, message.Content)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("读取聊天记录文件时出错:", err)
	}
}
