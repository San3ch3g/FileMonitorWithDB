package main

type DataForDB struct {
	ID       int    `json:"id"`
	FileName string `json:"fileName"`
	Dir      string `json:"dir"`
	FileData string `json:"fileData"`
	GUID     string `json:"guid"`
}
