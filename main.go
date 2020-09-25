package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const config = "config.txt"

type result struct {
	timeResult time.Duration
	errorInProcess error
}

func main() {

	var test rune = 'c'
	fmt.Printf("%c\n", test)
	for i := 0; i < 3; i++ {
		test++
		fmt.Printf("%c\n", test)
	}

	var boolean = true
	var number = 10
	var floatNumber = 10.1
	var comlexNumber = 10+3i
	var str = "GOOD DAY"
	fmt.Println("boolean " + )

	//cpuNumber := runtime.NumCPU()
	//fmt.Println("NumCPU", cpuNumber)
	//runtime.GOMAXPROCS(cpuNumber)
	////runtime.GOMAXPROCS(1)
	//
	//programTimer := time.Now()
	//
	//var logFilename = "logsFile.log"
	//file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//	log.Println("Инициалиация файла для логирования завершена. Сообщения будут выводиться в консоль!")
	//}
	//
	//defer file.Close()
	//log.SetOutput(file)
	//log.Println("Инициалиация файла для логирования завершена.")
	//
	//inputFile, outputFolder, err := setFlags()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//err = createFolder(outputFolder)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//urls, err := getTextFromFileLikeArray(inputFile)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//fmt.Println("ВРЕМЯ РАБОТЫ ДО СКАЧИВАНИЯ ФАЙЛОВ: ", time.Since(programTimer))
	//
	//var timesOfWork = time.Duration(0)
	//timeForCircle := time.Now()
	//
	//urlChan := make(chan string, len(urls))
	//timeResult := make(chan result, len(urls))
	//
	//for _, value := range urls {
	//	urlChan <- value
	//}
	//
	//for index, _ := range urls {
	//	filename :=  outputFolder + "\\file_" + fmt.Sprint(index) + "_.html"
	//	log.Println("Проверка ", filename, " на существование")
	//	if _, err := os.Stat(filename);
	//		os.IsNotExist(err) {
	//		log.Println("Начинается скачивание файла: ", filename)
	//		go download(filename, urlChan, timeResult)
	//	} else {
	//		log.Println(filename, " уже существует")
	//	}
	//}
	//
	//for range urls{
	//	value:=  <- timeResult
	//	if value.errorInProcess != nil {
	//		fmt.Println(value.errorInProcess)
	//	}
	//	timesOfWork += value.timeResult
	//}
	//
	//
	//
	//fmt.Println("НА ВЕСЬ ЦИКЛ УШЛО: ", time.Since(timeForCircle))
	//log.Println("Общее время по сохранению файлов: ", timesOfWork)
	//fmt.Println("ВСЯ ПРОГРАММА РАБОТАЛА: ", time.Since(programTimer))
}

func download(filename string, urlChan <- chan string, timeResult chan <- result) {
	url := <-urlChan
	log.Println("Начато подключение к URL: ", url)
	timeResult <- getResponseAndCreateFile(filename, url)
}

func getResponseAndCreateFile(filename string, url string) result {
	start := time.Now()
	response := getResponse(url)
	if response == nil {
		return result{timeResult: time.Duration(0),errorInProcess: errors.New("Плохой ответ по URL: " + url)}
	}
	defer response.Body.Close()
	file, err := createFile(filename, response)
	if err != nil {
		log.Println(err)
		return result{timeResult: time.Duration(0),errorInProcess: err}
	}
	timeForDownload := time.Since(start)
	log.Println("Сохранение файла ", file.Name(), " заняло по времени ", timeForDownload)
	return result{timeResult: timeForDownload, errorInProcess: nil}
}

func createFile(fileName string, response *http.Response) (*os.File,error){
	file, err := os.Create(fileName)
	if err != nil {
		return nil,err
	}
	bytes,err := io.Copy(file, response.Body)
	log.Println("Файл",file.Name(),"сохранён")
	log.Println("Получен файл на ", bytes, " байт")
	return file,nil
}

func getResponse(url string) *http.Response  {
	newClient := http.Client{Timeout: 10 * time.Second}
	response, err := newClient.Get(url)
	if err != nil{
		log.Println(err)
		return nil
	}
	if response == nil || response.StatusCode != 200{
		log.Println("Плохой ответ по адресу: ", url)
		return nil
	}
	log.Println("Подключение к URL: ", url, " прошло успешно. Получен ответ с сервера.")
	return response
}

func setLogOutput() {
	var logFilename = "logFile.log"
	file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
}

func createFolder(path string) error {
	log.Println("Начато создание папки для файлов")
	if _, err := os.Stat(path)
		os.IsNotExist(err) {
		err = os.MkdirAll(path, 0775)
		if err != nil{
			log.Println("Папка не была создана из-за ошибки: ",err)
			return err
		}
		log.Println("Папка была создана.")
		return nil
	} else {
		log.Println("Папка уже существует.")
		return nil
	}
}

func getTextFromFileLikeArray(filename string) ([]string, error) {
	log.Println("Попытка открыть файл: ", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Ошибка при попытке открыть файл: ", filename)
		log.Println(err)
		return nil,err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err();err != nil{
		return  nil,err
	}

	log.Println("Файл ",file.Name()," успешно прочитан.")
	return lines,nil
}

func setFlags() (string, string, error) {
	log.Println("Начато становление флагов.")
	log.Println("Начато чтение файла.")
	var configText, err = getTextFromFileLikeArray(config)
	if err != nil{
		log.Println("Ошибка при чтении фалов конфига становление флагов.")
		log.Println(err)
		return "","",err
	}
	inputFile := flag.String("inputFile", strings.TrimSuffix(configText[0], "\r"), "File with the URLs.")
	outputFolder := flag.String("outputFolder", strings.TrimSuffix(configText[1], "\r"), "Path to the folder, where the downloaded files are stored.")
	flag.Parse()
	log.Println("Становление флагов прошло успешно.")
	return *inputFile, *outputFolder,nil
}

func checkPath(path string){

}
