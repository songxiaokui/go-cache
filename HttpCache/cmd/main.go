package main

/*
@Time    : 2021/3/3 22:30
@Author  : austsxk
@Email   : austsxk@163.com
@File    : main.go
@Software: GoLand
*/
import (
	in "go_cache/HttpCache/internal"
)

func main() {
	c := in.Make("cache")
	in.New(&c).Listen()
}
