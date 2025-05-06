package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strings"
)

// GetTelegramUserID возвращает Telegram User ID из initData (tgWebAppData).
// botToken - токен вашего бота от @BotFather.
// initData - строка вида "query_id=...&user={...}&hash=...".
func GetTelegramUserID(botToken, initData string) (int64, error) {
	// 1. Проверяем подпись данных
	ok, data, err := verifyTelegramWebAppData(botToken, initData)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("invalid Telegram data signature")
	}

	// 2. Извлекаем JSON пользователя
	userJSON, exists := data["user"]
	if !exists {
		return 0, errors.New("user data not found in initData")
	}

	// 3. Парсим JSON и возвращаем ID
	var user struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

// verifyTelegramWebAppData проверяет подпись и возвращает разобранные данные.
func verifyTelegramWebAppData(botToken string, initData string) (bool, map[string]string, error) {
	parsed, err := url.ParseQuery(initData)
	if err != nil {
		return false, nil, err
	}

	// Извлекаем хэш и удаляем его из данных для проверки
	receivedHash := parsed.Get("hash")
	parsed.Del("hash")

	// Сортируем ключи и формируем строку для проверки
	var keys []string
	for key := range parsed {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var dataToCheck []string
	for _, key := range keys {
		dataToCheck = append(dataToCheck, key+"="+parsed.Get(key))
	}
	dataString := strings.Join(dataToCheck, "\n")

	// Вычисляем HMAC-SHA256
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	// Создаем новый HMAC для проверки данных
	checker := hmac.New(sha256.New, secretKey)
	checker.Write([]byte(dataString))
	computedHash := hex.EncodeToString(checker.Sum(nil))

	// Если хэши не совпадают, данные невалидны
	if computedHash != receivedHash {
		return false, nil, nil
	}

	// Возвращаем разобранные данные
	result := make(map[string]string)
	for key, values := range parsed {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}

	return true, result, nil
}