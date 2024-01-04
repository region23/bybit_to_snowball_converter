# Конвертер файла истории операций с сайта ByBit в формат [Snowball Income](https://snowball-income.com/register/ynhjcgcmjvbtscy)

### Чтобы выгрузить операции с ByBit

Перейди в раздел [Trade History](https://www.bybit.com/user/assets/order/fed/spot-history/active) > нажми кнопку "Export" в правом верхнем углу и выбери период, за который хочешь сделать выгрузку.

### Чтобы конвертировать полученный файл в формат Snowball Income

1. Если у тебя не установлен Golang, то перед следующим шагом обязательно [установи](https://go.dev/) его.
2. Склонируй этот репозиторий себе на компьютер.
3. В папке с проектом выполни команду `go mod download`, затем `go run main.go path/to/Bybit-History.xls`
4. После исполнения программы в папке появится файл `snowball_Bybit-History.csv`. Загрузи этот файл в новый портфель на Snowball Income.
