package counter

import (
	"fmt"
	"testCounterPackage/config"
)

func init() {
	SetIntervall(config.Api.Counter.Interval)
}

func main() {
	fmt.Println("Hello from counter main!")
}
