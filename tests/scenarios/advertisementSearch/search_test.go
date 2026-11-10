package advertisementSearch

import (
	"os"
	"testing"

	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"

	"api-tests-template/internal/utils"
	"strings"
)

const (
	existingAdTitle = "Носки"
	searchBeginning = "Но"
	searchMiddle    = "ос"
	searchEnd       = "ки"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Проверяем, что название и подстроки заданы")

	require.NotEmpty(s.T(), existingAdTitle)
	require.NotEmpty(s.T(), searchBeginning)
	require.NotEmpty(s.T(), searchMiddle)
	require.NotEmpty(s.T(), searchEnd)

	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}
func (s *TestSuite) TestSearchByTitle() {
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)

	s.Run("В ответе есть объявления", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}
func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}
