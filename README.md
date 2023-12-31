# FileMonitorWithDB
Этот проект представляет собой инструмент для мониторинга указанного каталога на наличие файлов с расширением ".tsv" (Если вам необходим другой тип данных, вы можете его изменить), взаимодействия с базой данных MySQL и предоставления простого HTTP API для извлечения данных на основе уникального идентификатора (GUID). Программа выполняет следующие ключевые функции:
## Основные возможности
1. Взаимодействие с базой данных:
    * Подключается к базе данных MySQL для проверки существования файлов и их вставки в базу данных при отсутствии.
    * Использует таблицу с именем `allFiles` и колонками: `ID`, `fileName`, `directory`, `fileData` и `guid`.

2. Анализ файлов:
    * Анализирует файлы с расширением ".tsv" в мониторируемом каталоге, преобразуя их в строковый формат, и сохраняет содержимое в базе данных.
3. Мониторинг файлов:
    * Регулярно отслеживает указанный каталог на наличие изменений и обновляет базу данных соответственно.
4. HTTP API:
    * Предоставляет простой HTTP API-эндпоинт ("/data") для извлечения информации о файлах на основе GUID.
## Начало работы
1. Настройка базы данных
    * Убедитесь, что у вас есть база данных MySQL с указанными учетными данными (username, password, hostname, port, dbName).
    * Создайте таблицу с именем allFiles и колонками: `ID`, `fileName`, `directory`, `fileData` и `guid`.
2. Настройка программы
    * Обновите константы в коде (username, password, hostname, port, dbName, directory, results, PortForServer) в соответствии с вашей конфигурацией.
3. Запуск программы
    ```go
    go run main.go
    ```
    * Программа начнет мониторинг указанного каталога и обслуживание HTTP API на настроенном порту.
## API-эндпоинт
* Извлечение данных:
    * Эндпоинт: `http://localhost:8080/data?guid=<GUID>`
    * Метод: GET
    * Параметры: guid (обязательно) - Уникальный идентификатор, связанный с каталогом.
## Структура Кода
Программа организована в несколько функций:
* `CheckExist` : Проверяет существование файлов в базе данных, вставляет новые файлы и обновляет соответствующие файлы ".doc".
* `Parser` : Анализирует файлы ".tsv" и преобразует их в строковый формат.
* `Visit` : Проходит по структуре каталога и вызывает функцию CheckExist для каждого подкаталога.
* `WriteToDoc` : Записывает данные в файл ".doc", обновляя существующее содержимое.
* `ReadFromDoc`: Считывает текущее содержимое файла ".doc".
* `data` : обслуживает HTTP API-эндпоинт /data, извлекая информацию из базы данных на основе переданного GUID и возвращая результат в формате JSON.
* `main` : инициализирует HTTP-сервер и запускает мониторинг файлов в указанной директории. Она также обслуживает HTTP API для запросов к данным о файлах.
## Константы конфигурации 
* `username`, `password`, `hostname`, `port`, `dbName`: Подробности подключения к базе данных.
* `directory` : Каталог для мониторинга файлов с расширением ".tsv".
* `results` : Каталог для сохранения результатов.
* `PortForServer` : Порт, на котором будет обслуживаться HTTP API.
## Структура Данных
Программа использует структуру DataForDB для представления информации о файлах как в базе данных, так и в ответах JSON.
```go
type DataForDB struct {
	ID       int    `json:"id"`
	FileName string `json:"fileName"`
	Dir      string `json:"dir"`
	FileData string `json:"fileData"`
	GUID     string `json:"guid"`
}
```
## Примечание
Убедитесь, что у вас есть необходимые разрешения и настройки для доступа к базе данных и мониторингу каталога.
