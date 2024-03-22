package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil) // здесь нужно создать запрос к сервису

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	// получаем ответ
	response := responseRecorder.Result()

	// здесь нужно добавить необходимые проверки
	//проверка корректности запроса
	require.NotNil(t, req, "Erorr creating request")
	//проверка возврата код ответа - 200
	require.Equal(t, 200, response.StatusCode, "Unexpected status code")
	//тело ответа не пустое
	require.NotEmpty(t, response.Body, "Response body is empty")
	//проверяем что код ответа 400
	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "Unexepected status code")
	//читаем тело ответа и сохраняем в переменную
	body, _ := io.ReadAll(response.Body)
	//проверяем что тело ответа содержитт сообщение об ошибке
	require.Contains(t, string(body), "Wrong city value")

	list := strings.Split(string(body), ",")
	//проверяем что длина списка кафе соответсыует ожидаемому общему количеству
	assert.Len(t, list, totalCount, "Unexpected number of cafes")

}
