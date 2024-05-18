package external

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
)

// WARN: Для сохранения этой информации можно было создать файл .env,
// для упрощения оставляем так
const TOKEN = "48ab34464a5573519725deb5865cc74c"

// genReqAuth - генерирует запрос с авторизацией
func genReqAuth(endPoint string, data map[string]string) (*http.Request, error) {
	op := "genReqAuth"
	var buffer bytes.Buffer
	formData := multipart.NewWriter(&buffer)

	for k, v := range data {
		if err := formData.WriteField(k, v); err != nil {
			log.Println(op, err)
			return nil, err
		}
	}
	formData.Close()
	req, err := http.NewRequest("POST", endPoint, &buffer)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	req.Header.Set("Content-Type", formData.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+TOKEN)
	return req, nil
}
