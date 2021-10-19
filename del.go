package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"log"
	"net/http"
)

// Структура используемая при отображении view.html
type DeathNote struct {
	SignatureCount int // хранение количества записей
	Signatures []string //хранение самих записей
}

//обработка ошибок
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//getStrings возвращает сегмент строк, прочитанный из fileName
//по одной строке на каждую строку файла
func getStrings(fileName string) []string {
	var line []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil //если файла нет возвращаем nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = append(line, scanner.Text())
	}
	check(scanner.Err())
	return line //возвращаем сегмент строк
}

//читает записи книги смерти и выводит их вместе со счетчиком
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	signatures := getStrings("signatures.txt") //читает записи из файла
	html, err := template.ParseFiles("view.html")
	check(err)
	deathNote:= DeathNote{
		SignatureCount: len(signatures), //хранение количества записей
		Signatures: signatures, //хранение самих записей
	}
	err = html.Execute(writer, deathNote) //данные структуры DeathNote выставляются в шаблон, результат записывается в ResponseWriter
	check(err)
	}

//отображает форму для ввода записей
func newHandler(writer http.ResponseWriter, request *http.Request)  {
	html, err:= template.ParseFiles("new.html")
	check(err)
	err=html.Execute(writer, nil)
	check(err)
}

//получает запрос POST с добавляемой записью
//и присоединяет ее к файлу signatures
func createHandler(writer http.ResponseWriter, request *http.Request)  {
	signature:= request.FormValue("signature") //получаем значение поля формы signature
	option:= os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err:= os.OpenFile("signatures.txt",option,os.FileMode(0600)) //добавляем в файл, если нет создаем его
	check(err)
	_, err= fmt.Fprintln(file,signature) //добавляем содержимое поле в фалй
	check(err)
	err=file.Close()
	check(err)
	http.Redirect(writer,request,"/deathNote",http.StatusFound)
}

func main() {
	http.HandleFunc("/deathNote", viewHandler)
	http.HandleFunc("/deathNote/new",newHandler)
	http.HandleFunc("/deathNote/create",createHandler)
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("./image"))))
	//http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
