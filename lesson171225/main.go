package main

import (
	"testCounterPackage/counter"
)

func main() {
	counter.Intervall = 5

	counter.Count()
	counter.PrintCount()

	counter.Count()
	counter.PrintCount()
}
