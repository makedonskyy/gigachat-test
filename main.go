package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// Config содержит конфигурационные параметры
type Config struct {
	AuthorizationKey string
}

// TokenResponse представляет структуру ответа API для получения токена
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// ChatRequest представляет структуру запроса к GigaChat API
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message представляет сообщение в диалоге
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse представляет структуру ответа от GigaChat API
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice представляет один из вариантов ответа
type Choice struct {
	Message struct {
		Content string `json:"content"`
		Role    string `json:"role"`
	} `json:"message"`
}

// ProductRequest представляет структуру запроса на анализ товара
type ProductRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Keywords string `json:"keywords"`
}

// TemplateData представляет данные для шаблона
type TemplateData struct {
	Result string
}

func main() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Ошибка загрузки .env файла: %v\n", err)
		os.Exit(1)
	}

	// Проверяем наличие ключа авторизации
	authKey := os.Getenv("GIGACHAT_AUTH_KEY")
	if authKey == "" {
		fmt.Println("Ошибка: не установлен GIGACHAT_AUTH_KEY")
		os.Exit(1)
	}

	config := Config{
		AuthorizationKey: authKey,
	}

	// Настраиваем CORS
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "http://localhost:5173" || origin == "http://localhost:5174" || origin == "http://localhost:5175" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	// Настраиваем маршруты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, TemplateData{})
	})

	http.HandleFunc("/analyze", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка чтения тела запроса: %v", err), http.StatusBadRequest)
			return
		}

		// Декодируем JSON
		var productReq ProductRequest
		if err := json.Unmarshal(body, &productReq); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Получаем токен доступа
		token, err := getAccessToken(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка получения токена: %v", err), http.StatusInternalServerError)
			return
		}

		// Формируем промпт для анализа
		prompt := fmt.Sprintf(`Проанализируй товар со следующими характеристиками:
Название: %s
Категория: %s
Ключевые слова: %s

Сделай анализ по следующим аспектам:
1. Целевая аудитория
2. Основные преимущества
3. Возможные недостатки
4. Рекомендации по улучшению
5. Потенциальные риски`, productReq.Name, productReq.Category, productReq.Keywords)

		// Отправляем запрос к GigaChat API
		response, err := sendChatRequest(token.AccessToken, prompt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при анализе: %v", err), http.StatusInternalServerError)
			return
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Choices[0].Message)
	}))

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func getAccessToken(config Config) (*TokenResponse, error) {
	// Генерируем UUID для RqUID
	rqUID := uuid.New().String()

	// Формируем URL и данные запроса
	apiURL := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	data := "scope=GIGACHAT_API_PERS"

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", rqUID)
	req.Header.Set("Authorization", "Basic "+config.AuthorizationKey)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: %s, тело ответа: %s", resp.Status, string(body))
	}

	// Декодируем JSON-ответ
	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &tokenResponse, nil
}

func sendChatRequest(token string, prompt string) (*ChatResponse, error) {
	// Формируем URL и данные запроса
	apiURL := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

	// Создаем структуру запроса
	chatRequest := ChatRequest{
		Model: "GigaChat",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Преобразуем запрос в JSON
	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("ошибка кодирования запроса: %v", err)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: %s, тело ответа: %s", resp.Status, string(body))
	}

	// Декодируем JSON-ответ
	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &chatResponse, nil
}
