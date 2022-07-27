package auth

import (
	Cfg "couriergate/configs"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ----- получить авторизационный токен используя Ключ доступа к API -----
// 01
func GetAuthTokenByApiKey() (string, error) {

	url := Cfg.ServiceURL + "/api/auth/logonByApiKey/?key=" + Cfg.ApiKey

	data, err := Web.SendPostRequest(url, "", "application/json", "")
	if err != nil {
		return "", errors.New("Error (AT-010101): " + fmt.Sprintf("%s\n", err))
	}

	var LogonResponse Models.LogonResponseOptions

	if err = json.Unmarshal([]byte(data), &LogonResponse); err != nil {
		return "", errors.New("Error (AT-010102): " + fmt.Sprintf("%s\n", err))
	}

	token := LogonResponse.Token

	return token, nil
}
