package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CheckExist(path string) error {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, hostname, port, dbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	// Проверяем наличие данных в базе данных для данной папки
	rows, err := db.Query("SELECT * FROM allFiles WHERE guid = ?", filepath.Base(path))
	if err != nil {
		return err
	}
	defer rows.Close()

	var dataList []DataForDB

	for rows.Next() {
		var column1 int
		var column2 string
		var column3 string
		var column4 string
		var column5 string
		err := rows.Scan(&column1, &column2, &column3, &column4, &column5)
		if err != nil {
			return err
		}
		data := DataForDB{
			ID:       column1,
			FileName: column2,
			Dir:      column3,
			FileData: column4,
			GUID:     column5,
		}
		dataList = append(dataList, data)
	}

	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		dirOfFile := path + "/" + file.Name()
		fileParts := strings.Split(file.Name(), ".")
		typeOfFile := ""
		if len(fileParts) > 1 {
			typeOfFile = fileParts[1]
		}
		if typeOfFile == "tsv" {
			var count int
			for i := 0; i < len(dataList); i++ {
				if file.Name() == dataList[i].FileName {
					count++
					break
				}
			}
			if count == 0 {
				// Вставка данных в базу данных
				insertQuery := fmt.Sprintf("INSERT INTO `allFiles` (fileName, directory, fileData, guid) VALUES ('%s', '%s', '%s', '%s');", file.Name(), dirOfFile, parser(dirOfFile), filepath.Base(path))

				_, err := db.Exec(insertQuery)
				if err != nil {
					return err
				}

				// Создание или обновление файла .doc
				docFileName := results + "/" + filepath.Base(path) + ".doc"
				err = writeToDoc(docFileName, DataForDB{
					FileName: file.Name(),
					Dir:      dirOfFile,
					FileData: parser(dirOfFile),
					GUID:     filepath.Base(path),
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func parser(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Создаем читатель CSV
	reader := csv.NewReader(file)
	reader.Comma = '\t' // Устанавливаем символ табуляции как разделитель

	// Читаем все записи из файла
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Преобразуем записи в строку
	var result string
	for _, record := range records {
		for _, field := range record {
			result += field + "\t" // Используйте тот же символ табуляции для объединения полей
		}
		result += "\n" // Добавляем новую строку после каждой записи
	}

	// Выводим результат
	return result
}

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		dirPath := string(path)
		err := CheckExist(dirPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeToDoc(fileName string, data DataForDB) error {
	// Читаем текущее содержимое файла
	existingContent, err := readFromDoc(fileName)
	if err != nil {
		return err
	}

	// Создаем или пересоздаем файл
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем существующее и новое содержимое
	_, err = file.WriteString(fmt.Sprintf("File Name: %s\nDirectory: %s\nFile Content:\n%s\n\n", data.FileName, data.Dir, data.FileData))
	if err != nil {
		return err
	}

	// Дописываем новые данные
	_, err = file.WriteString(existingContent)
	return err
}

func readFromDoc(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		// Если файл не существует, возвращаем пустую строку без ошибки
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func data(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Извлекаем значение параметра guid из запроса
	guid := r.URL.Query().Get("guid")

	// Если guid не указан, возвращаем ошибку
	if guid == "" {
		http.Error(w, "Parameter 'guid' is required", http.StatusBadRequest)
		return
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, hostname, port, dbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Проверяем наличие данных в базе данных для данной папки и указанного guid
	rows, err := db.Query("SELECT * FROM allFiles WHERE guid = ?", guid)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dataList []DataForDB

	for rows.Next() {
		var column1 int
		var column2 string
		var column3 string
		var column4 string
		var column5 string
		err := rows.Scan(&column1, &column2, &column3, &column4, &column5)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data := DataForDB{
			ID:       column1,
			FileName: column2,
			Dir:      column3,
			FileData: column4,
			GUID:     column5,
		}
		dataList = append(dataList, data)
	}

	// Обработка ошибки при json.Marshal
	bytes, err := json.Marshal(dataList)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
