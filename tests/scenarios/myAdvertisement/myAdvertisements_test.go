package myAdvertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"

	base "api-tests-template/tests"
	"api-tests-template/tests/support"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Позитивный сценарий: проверяем получение собственных объявлений в профиле.
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	_, created := support.CreateAdvertisementFixture(s.T(), s.loginData.Token)
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50)
	})

	var items gjson.Result
	s.Run("Проверяем, что созданное объявление присутствует в списке", func() {
		items = gjson.Get(advertisementsBody, "items")
		require.NotEmpty(s.T(), items.Array())
	})

	s.Run("Проверяем владельца и параметры созданного объявления", func() {
		found := false
		for _, item := range items.Array() {
			if gjson.Get(item.String(), "id").String() != created.ID {
				continue
			}
			found = true
			require.Equal(s.T(), s.loginData.User.Id, gjson.Get(item.String(), "user_id").String(), "user_id не совпадает")

			require.Equal(s.T(), created.ID, gjson.Get(item.String(), "id").String(), "id не совпадает")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "created_at").String(), "created_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "updated_at").String(), "updated_at пустое")
			require.Equal(s.T(), created.Title, gjson.Get(item.String(), "title").String(), "Название не совпадает")
			require.Equal(s.T(), created.Description, gjson.Get(item.String(), "description").String(), "Описание не совпадает")
			require.True(s.T(), len(gjson.Get(item.String(), "photos").Array()) > 0, "Присутствуют пустые фотографии")
			break
		}
		require.Truef(s.T(), found, "созданное объявление %s отсутствует в списке своих объявлений", created.ID)
	})
}

// Негативный сценарий: проверяем ошибку 401 при запросе собственных объявлений без токена.
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
