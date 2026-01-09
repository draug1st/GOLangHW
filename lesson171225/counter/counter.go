package counter

import "fmt"

var iterator int = 0
var Intervall int = 0

func Count() int {
	iterator += Intervall
	return iterator
}

func PrintCount() {
	fmt.Println("Current Count:", iterator)
}
