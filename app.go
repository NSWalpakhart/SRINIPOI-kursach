package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"html/template"
	"os"
	"strconv"
)

type PageData struct {
	Visits      int
	LastVisit   string
	Message     string
	GameState   string
	HiddenNum   string
	TimeLeft    int
	LastGuess   string
	GuessResult string
	GameOver    bool
}

var (
	visitCount int
	lastVisit  time.Time
	mu         sync.Mutex
	gameNumber     int
	gameStartTime  time.Time
	lastGuessTime  time.Time
)

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if r.Method == "GET" {
		gameNumber = time.Now().Nanosecond()%100 + 1
		gameStartTime = time.Now()
		lastGuessTime = time.Now()
	}

	data := PageData{
		Visits:    visitCount,
		LastVisit: lastVisit.Format("02.01.2006 15:04:05"),
		HiddenNum: "***",
		TimeLeft:  60 - int(time.Since(gameStartTime).Seconds()),
		GameOver:  false,
	}

	if data.TimeLeft < 0 {
		data.TimeLeft = 0
		data.GameOver = true
		data.GuessResult = "Время вышло! Игра окончена."
	}

	if r.Method == "POST" && !data.GameOver {
		if time.Since(lastGuessTime) > 5*time.Second {
			data.GuessResult = "Слишком долго думаете! Нужно отвечать быстрее 5 секунд."
		} else {
			guess := r.FormValue("guess")
			if guessNum, err := strconv.Atoi(guess); err == nil {
				if guessNum == gameNumber {
					data.GuessResult = "Поздравляем! Вы угадали число!"
					data.GameOver = true
				} else if guessNum < gameNumber {
					data.GuessResult = "Загаданное число больше!"
				} else {
					data.GuessResult = "Загаданное число меньше!"
				}
			}
		}
		lastGuessTime = time.Now()
	}

	visitCount++
	lastVisit = time.Now()

	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		log.Printf("Ошибка при парсинге шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
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

func init() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}
