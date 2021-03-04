package main

import (
	"flag"
	"fmt"
)

/*
@Time    : 2021/3/4 20:31
@Author  : austsxk
@Email   : austsxk@163.com
@File    : client.go
@Software: GoLand
*/

func main() {

	server := flag.String("h", "127.0.0.1", "cache server address")
	op := flag.String("c", "get", "command, could be get/set/del")
	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")
	flag.Parse()
	fmt.Println(*server, *op, *key, *value)

}
