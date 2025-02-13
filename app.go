package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"html/template"
)

type PageData struct {
	Visits    int
	LastVisit string
	Message   string
}

var (
	visitCount int
	lastVisit  time.Time
	mu         sync.Mutex
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Привет, Мир!</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f0f0f0;
        }
        .container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .stats {
            color: #666;
            font-size: 0.9em;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Message}}</h1>
        <div class="stats">
            <p>Количество посещений: {{.Visits}}</p>
            <p>Последнее посещение: {{.LastVisit}}</p>
        </div>
    </div>
</body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	visitCount++
	lastVisit = time.Now()
	mu.Unlock()

	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		log.Printf("Ошибка при парсинге шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Visits:    visitCount,
		LastVisit: lastVisit.Format("02.01.2006 15:04:05"),
		Message:   "Привет, 世界",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Ошибка при выполнении шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Запуск демо-приложения. Нажмите Ctrl+C для выхода...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
