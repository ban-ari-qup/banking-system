🏦 Banking System - Полное руководство установки
📋 Оглавление
Что это за проект?

Предварительные требования

Установка PostgreSQL

Установка Go

Настройка проекта

Использование API

Примеры запросов

🤔 Что это за проект?
Это банковская система с веб-интерфейсом API, которая позволяет:

✅ Создавать банковские аккаунты

✅ Входить в систему (аутентификация)

✅ Пополнять баланс

✅ Снимать деньги

✅ Переводить другим пользователям

✅ Просматривать историю операций

Технологии: Go, PostgreSQL, Chi Router, JWT-подобные сессии

⚙️ Предварительные требования
Установить Go

Установить PostgreSQL

Установить Git

Установить curl (или использовать Postman)

🐘 Установка PostgreSQL
Скачай установщик: https://www.postgresql.org/download/windows/

Запусти установщик:

Выбери версию PostgreSQL

Укажи пароль: password (запомни его!)

Порт оставь: 5432

Заверши установку

⚡ Установка Go

Скачай: https://golang.org/dl/

Установи: Запусти .msi файл

Проверь: Открой Command Prompt и введи:

cmd
go version

🚀 Настройка проекта
1. Скачай проект
bash
# Открой терминал/командную строку
git clone https://github.com/твой-username/mfp.git
cd mfp
2. Настрой базу данных
bash
# Создай базу (если еще не создана)
createdb mybank

# Примени миграции (создаст таблицы)
psql -U postgres -d mybank -f migrations/001_create_tables.sql
psql -U postgres -d mybank -f migrations/002_add_indexes.sql
3. Запусти сервер
bash
go run cmd/server/main.go
Увидишь сообщение: Server started at http://localhost:8080

🛠 Использование API
📝 Основные понятия:
HTTP Методы:
GET - получить данные (например, баланс)

POST - создать что-то (например, аккаунт)

DELETE - удалить что-то (например, аккаунт)

Типы данных:
Все запросы в формате JSON

Все ответы в формате JSON

Авторизация:
После входа система дает session_id в cookies

Этот session_id автоматически отправляется в каждом запросе

📮 Примеры запросов
1. 🆕 Регистрация нового аккаунта
Метод: POST
URL: http://localhost:8080/register
Тело запроса (JSON):

json
{
  "first_name": "Иван",
  "phone": "77001234567",
  "password": "1234",
  "age": 25
}

2. 🔐 Вход в систему
Метод: POST
URL: http://localhost:8080/login
Тело запроса:

json
{
  "phone": "77001234567",
  "password": "1234"
}

3. 💰 Пополнение баланса
Метод: POST
URL: http://localhost:8080/accounts/me/deposit?amount=1000

Curl команда:

bash
curl -X POST "http://localhost:8080/accounts/me/deposit?amount=1000" \
  -H "Content-Type: application/json"

  4. 🏧 Снятие денег
Метод: POST
URL: http://localhost:8080/accounts/me/withdraw?amount=500

Curl команда:

bash
curl -X POST "http://localhost:8080/accounts/me/withdraw?amount=500" \
  -H "Content-Type: application/json"

  5. 🔄 Перевод другому пользователю
Метод: POST
URL: http://localhost:8080/accounts/me/transfer
Тело запроса:

json
{
  "to": "77009876543",
  "amount": 300
}

6. 📊 Просмотр моего аккаунта
Метод: GET
URL: http://localhost:8080/accounts/me

Curl команда:

bash
curl http://localhost:8080/accounts/me

7. 📋 История транзакций
Метод: GET
URL: http://localhost:8080/accounts/me/transactions

Curl команда:

bash
curl http://localhost:8080/accounts/me/transactions

