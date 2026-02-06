package main

import "fmt"

func test() string {
	return "testtttt go"
}

func main() {
	var test string = test()
	var test2 string
	test2 = "test2!!"
	println(test, test2)
	fmt.Println("hello!!fmt")
}
