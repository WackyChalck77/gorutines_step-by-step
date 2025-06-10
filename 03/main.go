package main

import (
	"fmt"
	"net/http"
	"sync"
)

func getURLStat(i int, url string, wg *sync.WaitGroup) string {
	r, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при подключении", err)
	}
	r.Body.Close()
	fmt.Println(i, r.Status, url)
	defer wg.Done()
	return r.Status
}

func main() {
	var wg sync.WaitGroup
	//сюда надо слайс с адресами
	urls := []string{"https://avito.ru", "https://google.com",
		"https://ya.ru", "https://music.yandex.ru",
		"https://wildberries.ru/"}
	urls_len := len(urls)

	wg.Add(urls_len)
	//итерируем по слайсу и выдаем горутинам работу
	for i, url := range urls {
		//	fmt.Printf("%d, %s \n", i, url)
		go getURLStat(i, url, &wg) //передаем указатель
	}

	// var url string
	// url = "https://avito.ru"
	// getURLStat(url)
	//	time.Sleep(time.Second * 5)
	// go func() {
	// 	wg.Wait()
	// }()
	wg.Wait()
}
