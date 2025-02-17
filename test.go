package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		method        string
		formData      url.Values
		expectedCode  int
		expectedBody  string
	}{
		{
			name:          "GET запрос должен вернуть форму",
			method:        "GET",
			formData:      url.Values{},
			expectedCode:  http.StatusOK,
			expectedBody:  "Угадай число!",
		},
		{
			name:          "POST запрос с некорректным значением",
			method:        "POST",
			formData:      url.Values{"guess": {"abc"}},
			expectedCode:  http.StatusOK,
			expectedBody:  "Угадай число!",
		},
		{
			name:          "POST запрос с перезапуском игры",
			method:        "POST",
			formData:      url.Values{"restart": {"true"}},
			expectedCode:  http.StatusOK,
			expectedBody:  "Угадай число!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "GET" {
				req = httptest.NewRequest("GET", "/", nil)
			} else {
				req = httptest.NewRequest("POST", "/", strings.NewReader(tt.formData.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			w := httptest.NewRecorder()
			handler(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("handler вернул неверный код состояния: получили %v, ожидали %v",
					w.Code, tt.expectedCode)
			}

			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("handler вернул неожиданное тело: получили %v, ожидали содержание %v",
					w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestGameTimeLimits(t *testing.T) {
	// Сброс глобальных переменных
	mu.Lock()
	gameStartTime = time.Now().Add(-61 * time.Second) // Игра длится больше минуты
	lastGuessTime = time.Now().Add(-6 * time.Second)  // Последняя попытка была более 5 секунд назад
	mu.Unlock()

	req := httptest.NewRequest("POST", "/", strings.NewReader(url.Values{"guess": {"50"}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	
	handler(w, req)

	if !strings.Contains(w.Body.String(), "Слишком долго думаете!") {
		t.Error("Ожидалось сообщение о превышении времени между попытками")
	}
}

func TestGameState(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func()
		guess       string
		wantResult  string
	}{
		{
			name: "Правильное угадывание",
			setupFunc: func() {
				mu.Lock()
				gameNumber = 50
				gameStartTime = time.Now()
				lastGuessTime = time.Now()
				mu.Unlock()
			},
			guess: "50",
			wantResult: "Поздравляем! Вы угадали число!",
		},
		{
			name: "Число меньше загаданного",
			setupFunc: func() {
				mu.Lock()
				gameNumber = 50
				gameStartTime = time.Now()
				lastGuessTime = time.Now()
				mu.Unlock()
			},
			guess: "30",
			wantResult: "Загаданное число больше!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()

			form := url.Values{"guess": {tt.guess}}
			req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler(w, req)

			if !strings.Contains(w.Body.String(), tt.wantResult) {
				t.Errorf("Ожидалось сообщение '%s', получено: %s", tt.wantResult, w.Body.String())
			}
		})
	}
}

func TestVisitCounter(t *testing.T) {
	mu.Lock()
	visitCount = 0
	mu.Unlock()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	handler(w, req)
	
	mu.Lock()
	if visitCount != 1 {
		t.Errorf("После первого посещения счетчик = %d, ожидалось 1", visitCount)
	}
	mu.Unlock()
	
	handler(w, req)
	
	mu.Lock()
	if visitCount != 2 {
		t.Errorf("После второго посещения счетчик = %d, ожидалось 2", visitCount)
	}
	mu.Unlock()
} 
