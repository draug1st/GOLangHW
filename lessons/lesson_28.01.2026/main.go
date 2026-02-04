package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

// Константы для определения типа операции с балансом.
// Используются для избежания "магических строк" в коде и ошибок при написании.
const OPERATION_IN = "in"   // Операция пополнения
const OPERATION_OUT = "out" // Операция списания

// Balance - структура, описывающая кошелек (баланс).
// Содержит мьютекс для безопасного доступа из нескольких горутин одновременно.
type Balance struct {
	mu    sync.RWMutex // RWMutex позволяет эффективно управлять блокировками: много читателей или один писатель.
	value int          // Текущее значение баланса
	Limit int          // Лимит (максимальная сумма), которую может хранить кошелек
}

// Pay - метод для списания средств с баланса.
// Принимает сумму (sum), которую нужно списать.
// Возвращает ошибку, если средств недостаточно.
func (w *Balance) Pay(sum int) error {
	// Блокируем мьютекс на запись. Это значит, что пока эта функция выполняется,
	// никто другой не может ни читать, ни писать в w.value.
	// Это критически важно для предотвращения гонки данных (data race).
	w.mu.Lock()
	// defer гарантирует, что мьютекс разблокируется при выходе из функции,
	// даже если произойдет паника или return.
	defer w.mu.Unlock()

	// Проверяем, достаточно ли средств
	if w.value >= sum {
		w.value -= sum // Списываем средства
		return nil     // Возвращаем nil, так как ошибки нет
	}
	// Если средств не хватает, возвращаем ошибку
	return InsufficientFundsError
}

// Deposit - метод для пополнения баланса.
// Принимает сумму (sum), которую нужно добавить.
func (w *Balance) Deposit(sum int) error {
	// Снова блокируем на запись, так как будем менять состояние (w.value).
	w.mu.Lock()
	defer w.mu.Unlock()

	// Проверяем, не превысит ли баланс установленный лимит
	if w.value+sum <= w.Limit {
		w.value += sum // Пополняем
		return nil
	}
	// Если лимит превышен, возвращаем ошибку
	return OperationLimitError
}

// GetValue - метод для получения текущего значения баланса.
// Использует только чтение, поэтому используем RLock.
func (w *Balance) GetValue() int {
	// RLock (Read Lock) позволяет множеству горутин читать данные одновременно,
	// но запрещает кому-либо писать, пока есть хотя бы один читатель.
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.value
}

// GetLimit - метод для получения лимита.
// В данном случае поле Limit считается неизменяемым после инициализации (или не требует защиты,
// если мы уверены, что оно не меняется конкурентно), поэтому мьютекс не используется.
// Но для полной безопасности при возможности изменения Limit лучше тоже использовать RLock.
func (w *Balance) GetLimit() int {
	return w.Limit
}

// Payment - структура, описывающая платеж.
// Передается по каналам между воркерами.
type Payment struct {
	Id        int    // Уникальный идентификатор платежа
	Operation string // Тип операции (in или out)
	Amount    int    // Сумма операции
}

// Определение переменных ошибок, которые могут возникнуть.
// Это позволяет сравнивать ошибки с помощью errors.Is
var (
	InsufficientFundsError = errors.New("insufficient funds")       // Ошибка: недостаточно средств
	OperationLimitError    = errors.New("Operation limit exceeded") // Ошибка: превышен лимит
)

// processPayments - вспомогательная функция, которая решает, какой метод баланса вызвать
// в зависимости от типа операции в платеже.
func processPayments(wallet *Balance, payment Payment) error {
	switch payment.Operation {
	case OPERATION_IN:
		return wallet.Deposit(payment.Amount) // Вызываем пополнение
	case OPERATION_OUT:
		return wallet.Pay(payment.Amount) // Вызываем списание
	}
	return nil // Если операция неизвестна, ничего не делаем (можно было бы вернуть ошибку)
}

// successWorker - горутина, которая обрабатывает успешные платежи.
// Читает из канала successes и выводит информацию в консоль.
func successWorker(wg *sync.WaitGroup, successes <-chan Payment) {
	// Сообщаем WaitGroup, что эта горутина завершила работу, когда выйдем из функции.
	defer wg.Done()

	// range по каналу читает сообщения, пока канал не будет закрыт.
	for payment := range successes {
		fmt.Println(fmt.Sprintf("info: success payment %d ; amount %d", payment.Id, payment.Amount))
	}
	// Этот код выполнится только после того, как канал successes будет закрыт
	fmt.Println("success channel closed")
}

// failureWorker - горутина, которая обрабатывает ошибки.
// Читает из канала failures.
func failureWorker(wg *sync.WaitGroup, failures <-chan error) {
	defer wg.Done()

	// Читаем ошибки из канала, пока он открыт
	for err := range failures {
		fmt.Println(fmt.Sprintf("Log: error %s", err))
		// Проверяем тип ошибки и выводим соответствующее сообщение с разным уровнем логирования
		if errors.Is(err, InsufficientFundsError) {
			fmt.Println(fmt.Sprintf("critical : %s", err.Error()))
		}
		if errors.Is(err, OperationLimitError) {
			fmt.Println(fmt.Sprintf("warning: %s", err.Error()))
		}
	}
	fmt.Println("failure channel closed")
}

// worker - основная рабочая лошадка.
// Читает платежи из канала payments, обрабатывает их через кошелек
// и отправляет результат либо в successes, либо в failures.
func worker(n int, wg *sync.WaitGroup, wallet *Balance, payments <-chan Payment, successes chan<- Payment, failures chan<- error) {
	// Уменьшаем счетчик WaitGroup при завершении воркера
	defer wg.Done()

	// Читаем задачи (платежи) из канала payments.
	// Цикл завершится, когда канал payments будет закрыт и опустеет.
	for p := range payments {
		fmt.Println(fmt.Sprintf("worker %d processing payment %d", n, p.Id))

		// Пытаемся выполнить операцию с балансом
		err := processPayments(wallet, p)

		if err != nil {
			// Если произошла ошибка, логируем детали (текущий баланс или лимит)
			// и отправляем ошибку в канал failures
			if errors.Is(err, InsufficientFundsError) {
				b := wallet.GetValue()
				fmt.Println(fmt.Sprintf("insufficient funds %d", b))
			}
			if errors.Is(err, OperationLimitError) {
				fmt.Println(fmt.Sprintf("operation limit exceeded %d", wallet.GetLimit()))
			}
			failures <- err
		} else {
			// Если успешно, отправляем платеж в канал successes
			successes <- p
		}
	}
}

// generatePayments - функция-генератор, которая создает платежи.
// Работает в отдельной горутине.
func generatePayments(count int, payments chan<- Payment) {
	operation := ""

	// Очень важно закрыть канал payments после того, как мы закончили писать в него.
	// Это сигнал для всех воркеров (worker), что работы больше не будет,
	// и они могут завершить свои циклы range.
	defer close(payments)

	for i := 1; i <= count; i++ {
		// Генерируем случайное число для выбора типа операции
		n := rand.Intn(100)
		if n > 50 {
			operation = OPERATION_OUT
		} else {
			operation = OPERATION_IN
		}

		// Отправляем новый платеж в канал.
		// Если канал буферизированный и полон, или небуферизированный и никто не читает,
		// эта операция заблокируется.
		payments <- Payment{i, operation, rand.Intn(1000)}
	}
}

func main() {
	// Создаем каналы для обмена данными
	payments := make(chan Payment)  // Канал для задач (входящих платежей)
	successes := make(chan Payment) // Канал для успешных платежей
	failures := make(chan error)    // Канал для ошибок

	workerCount := 20 // Количество воркеров, которые будут обрабатывать платежи параллельно
	// Инициализируем кошелек с начальным значением и лимитом
	wallet := Balance{value: 100000, Limit: 10000000}

	// wg (WaitGroup) используется для ожидания завершения работы воркеров обработки (worker)
	var wg sync.WaitGroup
	// cg (Consumer WaitGroup) используется для ожидания завершения воркеров-потребителей результатов (successWorker и failureWorker)
	var cg sync.WaitGroup

	// Запускаем обработчики результатов (успех и ошибки)
	cg.Add(2) // Добавляем 2, так как запускаем две такие горутины
	go successWorker(&cg, successes)
	go failureWorker(&cg, failures)

	// Запускаем пул воркеров для обработки платежей
	for i := 0; i < workerCount; i++ {
		wg.Add(1) // Увеличиваем счетчик ожидаемых горутин
		go worker(i, &wg, &wallet, payments, successes, failures)
	}

	// Запускаем генерацию платежей в отдельной горутине.
	// Она будет наполнять канал payments.
	go generatePayments(10000, payments)

	// Основной поток блокируется здесь и ждет, пока ВСЕ воркеры (worker) завершат работу.
	// Воркеры завершатся, когда закроется канал payments (в generatePayments) и они дочитают всё из него.
	wg.Wait()

	// После того как все воркеры закончили работу, мы точно знаем,
	// что никто больше не будет писать в каналы failures и successes.
	// Поэтому их можно безопасно закрыть.
	close(failures)
	close(successes)

	// Теперь ждем, пока обработчики результатов (successWorker, failureWorker)
	// дочитают остатки из каналов и завершатся.
	cg.Wait()

	fmt.Println("done")
}
