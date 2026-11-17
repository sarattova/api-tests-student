package advertisementSearch

import (
	"strings"
	"testing"

	"api-tests-template/internal/assertions"
	advertisementManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/support"
)

func TestSearchByTitleSubstring(t *testing.T) {
	user := support.LoginPrimary(t)
	request, _ := support.CreateAdvertisementFixture(t, user.Token)
	title := request.Title
	titleRunes := []rune(title)
	beginning := string(titleRunes[:len(titleRunes)-1])
	middle := string(titleRunes[1 : len(titleRunes)-1])
	end := string(titleRunes[1:])

	cases := []struct {
		name   string
		search string
	}{
		{"полное название", title},
		{"начало", beginning},
		{"середина", middle},
		{"конец", end},
		{"нижний регистр", strings.ToLower(middle)},
		{"верхний регистр", strings.ToUpper(middle)},
		{"смешанный регистр", utils.MixCase(middle)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Шаг: ищем объявление по запросу %q", tc.search)
			list := advertisementManager.GetAdvertisementsList(t, user.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(t, list.Items, title)
			t.Logf("Результат: объявление %q найдено", title)
		})
	}
}
