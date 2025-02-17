package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		method        string
		guessValue    string
		expectedCode  int
		expectedBody  string
	}{
		{
			name:          "GET запрос должен вернуть форму",
			method:        "GET",
			guessValue:    "",
			expectedCode:  http.StatusOK,
			expectedBody:  "Угадай число!",
		},
		{
			name:          "POST запрос с некорректным значением",
			method:        "POST",
			guessValue:    "abc",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "Некорректное значение",
		},
		{
			name:          "POST запрос со значением вне диапазона",
			method:        "POST",
			guessValue:    "101",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "Число должно быть от 1 до 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "GET" {
				req = httptest.NewRequest("GET", "/", nil)
			} else {
				form := url.Values{}
				form.Add("guess", tt.guessValue)
				req = httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
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

func TestGameLogic(t *testing.T) {
	tests := []struct {
		name     string
		target   int
		guess    int
		expected string
	}{
		{
			name:     "Угадано правильно",
			target:   50,
			guess:    50,
			expected: "Поздравляем! Вы угадали число!",
		},
		{
			name:     "Число меньше загаданного",
			target:   50,
			guess:    30,
			expected: "Загаданное число больше",
		},
		{
			name:     "Число больше загаданного",
			target:   50,
			guess:    70,
			expected: "Загаданное число меньше",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hiddenNumber = tt.target // устанавливаем тестовое значение
			result := checkGuess(tt.guess)
			if result != tt.expected {
				t.Errorf("checkGuess(%d) = %v, ожидалось %v", 
					tt.guess, result, tt.expected)
			}
		})
	}
}

func TestVisitCounter(t *testing.T) {
	// Сброс счетчика перед тестом
	visits = 0
	
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	// Первое посещение
	handler(w, req)
	if visits != 1 {
		t.Errorf("После первого посещения счетчик = %d, ожидалось 1", visits)
	}
	
	// Второе посещение
	handler(w, req)
	if visits != 2 {
		t.Errorf("После второго посещения счетчик = %d, ожидалось 2", visits)
	}
} 
