package support

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	shared "api-tests-template/internal/managers/models"
	base "api-tests-template/tests"

	"github.com/stretchr/testify/require"
)

const (
	fixtureDescription = "Объявление создано автоматическим API-тестом"
	fixturePrice       = 1234.5
	fixtureQuantity    = 2
)

// LoginPrimary загружает настройки и авторизует основного тестового пользователя.
func LoginPrimary(t *testing.T) authModels.LoginOkResponse {
	t.Helper()
	base.SetupSuite()
	username := os.Getenv("TEST_LOGIN")
	password := os.Getenv("TEST_PASSWORD")
	require.NotEmpty(t, username, "укажите TEST_LOGIN в .env.override")
	require.NotEmpty(t, password, "укажите TEST_PASSWORD в .env.override")
	t.Log("Подготовка: авторизуем основного тестового пользователя")
	return auth.Login(t, username, password)
}

// LoginSecondary авторизует второго пользователя для проверки избранного и владельца объявления.
func LoginSecondary(t *testing.T) authModels.LoginOkResponse {
	t.Helper()
	base.SetupSuite()
	username := os.Getenv("TEST_OTHER_LOGIN")
	password := os.Getenv("TEST_OTHER_PASSWORD")
	require.NotEmpty(t, username, "укажите TEST_OTHER_LOGIN в .env.override; см. README")
	require.NotEmpty(t, password, "укажите TEST_OTHER_PASSWORD в .env.override; см. README")
	require.NotEqual(t, os.Getenv("TEST_LOGIN"), username, "для тестов нужны две разные учётные записи TestBoard")
	t.Log("Подготовка: авторизуем второго тестового пользователя")
	return auth.Login(t, username, password)
}

// UniqueTitle создаёт короткое название, которое не повторит названия из прошлых запусков.
func UniqueTitle(t *testing.T) string {
	t.Helper()
	var randomBytes [8]byte
	_, err := rand.Read(randomBytes[:])
	require.NoError(t, err, "не удалось создать уникальное название объявления")
	return "Тестовое объявление " + hex.EncodeToString(randomBytes[:])
}

// PhotoPath возвращает путь к тестовому PNG независимо от текущего каталога.
func PhotoPath(t *testing.T) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	require.True(t, ok, "не удалось найти тестовую фотографию")
	photoPath := filepath.Join(filepath.Dir(source), "testdata", "photo.png")
	_, err := os.Stat(photoPath)
	require.NoError(t, err, "тестовая фотография отсутствует")
	return photoPath
}

// CreateAdvertisementFixture создаёт объявление со всеми полями и регистрирует его удаление после теста.
func CreateAdvertisementFixture(t *testing.T, token string) (adModels.CreateAdvertisementRequest, adModels.Advertisement) {
	t.Helper()
	request := adModels.CreateAdvertisementRequest{
		Title:       UniqueTitle(t),
		Description: fixtureDescription,
		Price:       fixturePrice,
		Quantity:    fixtureQuantity,
	}
	t.Logf("Подготовка: создаём объявление %q", request.Title)
	result := advertisements.CreateAdvertisement(t, token, request, []string{PhotoPath(t)})
	require.Equalf(t, http.StatusCreated, result.StatusCode, "ошибка создания объявления: %+v", result.Error)
	require.NotEmpty(t, result.Data.ID, "у созданного объявления отсутствует ID")
	t.Cleanup(func() {
		existing := advertisements.GetAdvertisement(t, result.Data.ID)
		if existing.StatusCode == http.StatusNotFound {
			t.Logf("Очистка: объявление %s уже удалено", result.Data.ID)
			return
		}
		if existing.StatusCode != http.StatusOK {
			t.Errorf("очистка: не удалось проверить объявление %s: статус %d, ошибка %+v",
				result.Data.ID, existing.StatusCode, existing.Error)
			return
		}
		cleanup := advertisements.DeleteAdvertisement(t, token, result.Data.ID)
		if cleanup.StatusCode != http.StatusNoContent {
			t.Errorf("очистка: не удалось удалить объявление %s: статус %d, ошибка %+v",
				result.Data.ID, cleanup.StatusCode, cleanup.Error)
		}
		t.Logf("Очистка: объявление %s удалено", result.Data.ID)
	})
	return request, result.Data
}

// RequireAPIError проверяет HTTP-статус и содержательное сообщение об ошибке API.
func RequireAPIError[T any](t *testing.T, result shared.Result[T], status int, messagePart string) {
	t.Helper()
	require.Equalf(t, status, result.StatusCode, "неожиданный ответ с ошибкой: %+v", result.Error)
	require.NotEmpty(t, result.Error.Error, "код ошибки API пустой")
	require.NotEmpty(t, result.Error.Message, "сообщение об ошибке API пустое")
	if messagePart != "" {
		require.Contains(t, strings.ToLower(result.Error.Message), strings.ToLower(messagePart),
			fmt.Sprintf("неожиданное сообщение об ошибке API: %q", result.Error.Message))
	}
}
