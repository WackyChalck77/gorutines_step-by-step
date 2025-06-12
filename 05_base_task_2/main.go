package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	//печатаем результат функции Do, которая берет контекст
	//по бэкграунду и структуру юзер с перечисленными именами
	// fmt.Println(Do(context.Background(), []User{
	// 	{"aaa"}, {"bbb"}, {"ccc"},
	// 	{"ddd"}, {"eee"}}))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, _ := Do(ctx, []User{
		{"aaa"}, {"bbb"}, {"ccc"},
		{"ddd"}, {"eee"}})
	//	time.Sleep(time.Second * 5)
	fmt.Println("\nИтоговая мапа", result)
}

type User struct { //структура юзера с именем
	Name string
}

// ф-я получения по имени. вход: контекст, имя юзера
// выход: целое число, ошибка
func fetchByName(ctx context.Context, userName string) (int, error) {
	// Тут происходит сетевой поход, который по userName возвращает userID
	//добавим рандомное число от нуля до 1000
	netDuration := rand.Intn(6000)
	time.Sleep(time.Duration(netDuration) * time.Millisecond) // Имитация сетевого похода
	//создадим проверку завершения сетевого похода
	select {
	case <-ctx.Done():
		fmt.Println("Контекст завершен, не даждались сети,\nвремя работы одного из сетевых походов", time.Duration(netDuration), "\n")
		return 0, ctx.Err()
	default:
	}
	return rand.Int() % 100000, nil //возвращаем рандомное число
	//в пределах 99 999
}

// функция принимает контекст, структуру User, где имя-string
// выход: мапа строка-целое число, ошибка
func Do(ctx context.Context, users []User) (map[string]int, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	//мапа строка-целое
	collected := make(map[string]int)
	//итерируем пропуская порядковый номер в диапазоне значений
	//поступающей на вход в функцию структуры
	for n, u := range users {
		//создадим горутину, добавляем счетчик в вэйтгруп
		wg.Add(1)
		go func(u User) {
			defer wg.Done()
			fmt.Println("Работает горутина", n)
			//передаем в функцию контекст и каждое имя
			userID, err := fetchByName(ctx, u.Name)
			if err != nil { //обрабатываем ошибку
				fmt.Printf("Ошибка %v при получении имени %s!", err, u.Name)
				return
			}
			fmt.Println("Горутина", n, "пишет", u.Name)
			mu.Lock()
			collected[u.Name] = userID //наполняем мапу именем стуктуры
			mu.Unlock()
			//а ключ выдаем как userID из ф-ции fetchByName
		}(u)
	}
	wg.Wait()
	return collected, nil //возвращаем мапу и ноль, если нет ошибки
}
