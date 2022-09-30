package main

import "os"

func test() {
	// тут не должно быть сработки так как функция не main
	os.Exit(0)
}

func main() {
	// анализатор должен находить ошибку,
	os.Exit(0) // want "os.Exit was being detected!"
}
