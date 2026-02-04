package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NumHipsters   = 10000
	StoreCapacity = 300
	IPhonePrice   = 1800
)

var (
	soldCount      SafeCounter
	failedCount    SafeCounter
	totalProcessed SafeCounter
)

type SafeCounter struct {
	mu sync.Mutex
	v  int64
}

func (sc *SafeCounter) Inc() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.v++
}

func (sc *SafeCounter) Value() int64 {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.v
}

type Store struct {
	balance int64
	mu      sync.Mutex
}

func (s *Store) BuyIphone(cost int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.balance += cost
}

type Hipster struct {
	ID      int
	Balance int64
}

func main() {
	store := &Store{}

	hipsterQueue := make(chan Hipster, NumHipsters)

	fmt.Println("Генерация очереди хипстеров...")
	for i := 1; i <= NumHipsters; i++ {
		balance := int64(rand.Intn(4001))
		hipsterQueue <- Hipster{ID: i, Balance: balance}
	}
	close(hipsterQueue)

	fmt.Println("Магазин откроется через 3 секунды...")
	timer := time.NewTimer(3 * time.Second)
	<-timer.C
	fmt.Println(">>> Магазин открыт! <<<")

	var wg sync.WaitGroup
	wg.Add(StoreCapacity)

	fmt.Printf("Касс запущено: %d \n", StoreCapacity)

	for i := 0; i < StoreCapacity; i++ {
		go func(workerID int) {
			defer wg.Done()
			for hipster := range hipsterQueue {
				processHipster(store, hipster)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Магазин закрыт. Все посетители обслужены.")

	printStats(store)
}

func processHipster(store *Store, h Hipster) {
	if h.Balance >= IPhonePrice {
		store.BuyIphone(IPhonePrice)
		h.Balance -= IPhonePrice
		soldCount.Inc()
	} else {
		failedCount.Inc()
	}
	totalProcessed.Inc()
}

func printStats(store *Store) {
	fmt.Println("\n-------------------------------------------")
	fmt.Println("        Дневной отчет о продажах           ")
	fmt.Println("-------------------------------------------")
	fmt.Printf("Всего хипстеров:          %d\n", totalProcessed.Value())
	fmt.Printf("iPhones продано:          %d\n", soldCount.Value())
	fmt.Printf("Ушли плакать в Starbucks: %d\n", failedCount.Value())
	fmt.Printf("Магазин заработал:        $%d\n", store.balance)
	fmt.Println("-------------------------------------------")
}
