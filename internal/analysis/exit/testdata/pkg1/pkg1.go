package pkg1

import "os"

func main() {
	// Тут мы на должны получить сработку анализатора, так как пакет pkg1 а не main
	os.Exit(0)
}
