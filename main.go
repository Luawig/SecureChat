package main

import (
	"SecureChat/client"
	"SecureChat/common"
	"SecureChat/server"
	"flag"
	"fmt"
	"os"
)

func main() {
	mode := flag.String("mode", "client", "启动模式: client 或 server")
	action := flag.String("action", "online", "操作: online 或 view")
	username := flag.String("username", "", "用户名")
	password := flag.String("password", "", "密码")
	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Println("用户名和密码不能为空")
		os.Exit(1)
	}

	switch *mode {
	case "client":
		if *action == "online" {
			client.Start(*username, *password)
		} else if *action == "view" {
			common.ViewChatLog(*username, *password)
		} else {
			fmt.Println("未知的操作")
			os.Exit(1)
		}
	case "server":
		if *action == "online" {
			server.Start(*username, *password)
		} else if *action == "view" {
			common.ViewChatLog(*username, *password)
		} else {
			fmt.Println("未知的操作")
			os.Exit(1)
		}
	default:
		fmt.Println("未知的启动模式")
		os.Exit(1)
	}
}
