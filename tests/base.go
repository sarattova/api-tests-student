package tests

import (
	"log"

	"api-tests-template/internal/utils"
)

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Загружаем переменные окружения")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Завершаем набор тестов")
}

// Precondition выводит текст в консоль с меткой «Предусловие».
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Предусловие", text)
}
