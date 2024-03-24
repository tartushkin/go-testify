package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil) // здесь нужно создать запрос к сервису

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)
	// получаем ответ
	response := responseRecorder.Result()
	//проверка корректности запроса
	require.NotNil(t, req, "Erorr creating request")
	//проверка возврата код ответа - 200
	require.Equal(t, 200, response.StatusCode, "Unexpected status code")
	//тело ответа не пустое
	require.NotEmpty(t, response.Body, "Response body is empty")
}

func TestMainHandlerWhereIsTheWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=Ryazan", nil) // здесь нужно создать запрос к сервису

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)
	// получаем ответ
	response := responseRecorder.Result()
	//проверяем что код ответа 400
	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "Unexepected status code")
	//читаем тело ответа и сохраняем в переменную
	body, _ := io.ReadAll(response.Body)
	//проверяем что тело ответа содержитт сообщение об ошибке
	require.Contains(t, string(body), "wrong city value")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil) // здесь нужно создать запрос к сервису

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)
	// получаем ответ
	response := responseRecorder.Result()

	body, _ := io.ReadAll(response.Body)
	list := strings.Split(string(body), ",")
	//проверяем что длина списка кафе соответсыует ожидаемому общему количеству
	assert.Len(t, list, totalCount, "Unexpected number of cafes")
}
