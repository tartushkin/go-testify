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

var handler http.Handler

func createTestRequest(method, url string) (*http.Request, error) {
	req := httptest.NewRequest(method, url, nil)
	return req, nil
}

func TestMain(m *testing.M) {
	handler = http.HandlerFunc(MainHandle)
}

func TestMainHandlerWhenOk(t *testing.T) {
	req, err := createTestRequest("GET", "/cafe?count=4&city=moscow") // здесь нужно создать запрос к сервису
	//проверка корректности запроса
	require.NoError(t, err, "Erorr creating request")

	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)
	// получаем ответ
	response := responseRecorder.Result()
	//проверка возврата код ответа - 200
	require.Equal(t, 200, response.StatusCode, "Unexpected status code")
	//тело ответа не пустое
	require.NotEmpty(t, response.Body, "Response body is empty")
}

func TestMainHandlerWhereIsTheWrongCity(t *testing.T) {
	req, err := createTestRequest("GET", "/cafe?count=4&city=UnExistsCity") // здесь нужно создать запрос к сервису
	require.NoError(t, err, "Error creating request")

	responseRecorder := httptest.NewRecorder()
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
	req, err := createTestRequest("GET", "/cafe?count=5&city=moscow") // здесь нужно создать запрос к сервису
	require.NoError(t, err, "Error creating request")

	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)
	// получаем ответ
	response := responseRecorder.Result()

	body, _ := io.ReadAll(response.Body)
	list := strings.Split(string(body), ",")
	//проверяем что длина списка кафе соответсыует ожидаемому общему количеству
	assert.Len(t, list, totalCount, "Unexpected number of cafes")
	//проверка возврата код ответа - 200
	require.Equal(t, 200, response.StatusCode, "Unexpected status code")
}
