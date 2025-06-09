package main

import (
	"fmt"
)

func main() {
	//var mu sync.Mutex

	var max int // переменная int
	//приравнякм ее к нуля бля начала
	max = 0
	for i := 1000; i > 0; i-- { //итерируем
		//от тысячи до нуля
		go func(i int) { //получаем разный обход, нужно передать
			//в качестве аргумента i, чтобы избежать data race
			if i%2 == 0 && i > max {
				//	mu.Lock() //мьютексы хороши, но и без них работает
				max = i
				//	mu.Unlock()
			}
		}(i)
	}
	fmt.Printf("Maximum is %d", max)
}
