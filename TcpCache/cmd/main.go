package main

/*
@Time    : 2021/3/4 20:27
@Author  : austsxk
@Email   : austsxk@163.com
@File    : main.go
@Software: GoLand
*/
import . "go_cache/TcpCache/internal"

func main() {
	c := Make("cache")
	go NewTcpClient(&c).Listen()
	New(&c).Listen()
}
