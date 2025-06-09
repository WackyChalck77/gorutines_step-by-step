package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var counter int32 //счетчик

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	ch := make(chan int) //канал

	for i := 1; i <= 15; i++ { //запускаем 3 горутины
		wg.Add(1)
		go func(i int) { //на вход поступает тип int
			defer wg.Done()
			mu.Lock() //ставим мьютексы
			ch <- i   //закидываем в канал номер от 0 до 2
			mu.Unlock()
			atomic.AddInt32(&counter, 1) //выполняем атомарно
			//counter++ //инкрементируем счётчик
			//		time.Sleep(time.Millisecond * 50)

		}(i) //даём переменную на вход
	}

	go func() { //эта горутина ждет завершения всех рабочих
		//горутин и закрывает канал
		wg.Wait()
		close(ch)
	}()

	//это выполняет главная горутина
	for val := range ch { //итерация по всему диапазону канала
		fmt.Println("Got:", val) //показываем что было в канале
		//канал автоматически закрывается?
	}
	fmt.Println("Final counter:", counter) //результат по счётчику
}
