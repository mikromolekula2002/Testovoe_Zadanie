package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/utils"
)

func main() {
	fmt.Println("Запуск кода для создания Базы Данных(а также таблицы).")
	//Загрузка конфига
	cfg := config.LoadConfig("./config/config.yaml")

	utils.CreateDatabase(cfg)

	fmt.Println("Код завершил работу.")
}
