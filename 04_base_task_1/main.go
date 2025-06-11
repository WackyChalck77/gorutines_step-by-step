package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// функция берет контекст, имя. выдает строку, ошибку
func GetFile(ctx context.Context, name string) ([]byte, error) {
	if name == "" { //если пустое имя
		//возвращаем ноль, печатаем ошибку
		return nil, fmt.Errorf("invalid name %q", name)
	}

	ticker := time.NewTicker(time.Second) //раз в секунду идет тик
	//ждем какое-то время между попытками получить файл
	defer ticker.Stop() //останавливаем его в конце функции

	select { //ждем события. Блокируемся до одного из кейсов
	case <-ctx.Done(): //если контекст отменен
		return nil, ctx.Err() //возврат нуля и ошибки по контексту
	case <-ticker.C: //блокируемся до сигнала завершения
		//тикера, потом продолжаем дальше!
	}
	//проверка имени файла
	if strings.HasPrefix(name, "invalid") { //если name содержит
		//префикс "инвалид" - возвращаем ноль +
		//ошибку с текстовым описанием
		return nil, fmt.Errorf("invalid name %q", name)
	}
	//генерируем случайные байты
	b := make([]byte, 10)  //массив байт (8 бит), len=10
	n, err := rand.Read(b) //считываем случайные данные и записываем их
	//в байтовый слайс b. n - количество записаных байтов
	if err != nil {
		//при ошибке возвращаем ноль, формируем ошибку с именем
		return nil, fmt.Errorf("getting file %q: %w", name, err)
	}
	//возвращаем слайс с содержанием от нуля до n значений
	return b[:n], nil //и ноль в качестве ошибки
}

// GetFilesOld пример функции, которую нужно оптимизировать.
// Менять эту функцию не нужно, она нужна чтоб
// сравнить поведение двух функций после оптимизации
func GetFilesOld(ctx context.Context, names ...string) (result map[string][]byte, err error) {
	//принимаем контекст и неопределенное количество имен
	//возвращаем мапу ключ: строки- значение: байты, ошибку
	if len(names) == 0 {
		//если длина имени равна нулю
		//возвращаем: ноль-ноль
		return nil, nil
	}

	result = make(map[string][]byte, len(names))
	//результат равен мапе строка-байт, длины строки names
	for _, name := range names {
		//итерируем по диапазону names
		//по ключу имени записываем значение, которое
		//формирует функция GetFile в соответствии с контекстом и именем
		result[name], err = GetFile(ctx, name)
		if err != nil { //обрабатываем ошибку
			return nil, err
		}
	}

	return result, nil //возвращаем мапу с именами и в значении
	//сгенерированные темы
}

// GetFilesNew эту функцию можно менять, за исключением её
// сигнатуры
func GetFilesNew(ctx context.Context, names ...string) (result map[string][]byte, err error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	//здесь та же функция
	//принимаем контекст и неопределенное количество имен
	//возвращаем мапу ключ: строки- значение: байты, ошибку
	if len(names) == 0 {
		return nil, nil
	}
	result = make(map[string][]byte, len(names))
	wg.Add(len(names))
	for _, name := range names {
		//вкладыаем здесь горутины, которые исполняются параллельно
		go func(name string) {

			fileData, err := GetFile(ctx, name)
			if err != nil {
				fmt.Println("ошибка в горутине", name, err)
				return
			}
			mu.Lock() //чтобы не писать в мапу параллельно -делаем
			//mutual exception
			result[name] = fileData
			mu.Unlock()
			defer wg.Done()
		}(name)
	}

	wg.Wait()
	return result, nil
}

func main() {

	start := time.Now() //старт - это текущее время
	//файлы и ошибку получаем из новой функции, давай на вход
	//введем контекст
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(1015)*time.Millisecond)
	defer cancel()
	files, err := GetFilesNew(ctx,
		"one", "2", "three", "four", "five")
	if err != nil {
		//возвращаем лог с временем после старта и ошибку
		log.Fatalln(time.Since(start), err)
	}

	fmt.Println(time.Since(start)) //печатаем время после старта

	for value, name := range files {
		fmt.Println(name, value) //выводим все имена и присвоенные значения
	}

}
