# 🏦 Banking System

Полное руководство по установке и использованию банковской системы

## 📋 Оглавление

1. [Что это за проект?](#-что-это-за-проект)
2. [Предварительные требования](#️-предварительные-требования)
3. [Установка PostgreSQL](#-установка-postgresql)
4. [Установка Go](#-установка-go)
5. [Настройка проекта](#-настройка-проекта)
6. [Использование API](#-использование-api)
7. [Примеры запросов](#-примеры-запросов)
8. [Решение проблем](#-решение-проблем)

---

## 🤔 Что это за проект?

Это **банковская система** с веб-интерфейсом API, которая позволяет:

- ✅ Создавать банковские аккаунты
- ✅ Входить в систему (аутентификация)
- ✅ Пополнять баланс
- ✅ Снимать деньги
- ✅ Переводить другим пользователям
- ✅ Просматривать историю операций

**Технологии:** Go, PostgreSQL, Chi Router, JWT-подобные сессии

---

## ⚙️ Предварительные требования

### Необходимое ПО:
- [ ] Установить **Go** (версия 1.19+)
- [ ] Установить **PostgreSQL** (версия 14+)
- [ ] Установить **Git**
- [ ] Установить **curl** (или использовать Postman)

---

## 🐘 Установка PostgreSQL

1. **Скачай установщик:** [https://www.postgresql.org/download/windows/](https://www.postgresql.org/download/windows/)
2. **Запусти установщик:**
   - Выбери версию PostgreSQL
   - Укажи пароль: `password` (запомни его!)
   - Порт оставь: `5432`
3. **Заверши установку**

## ⚡ Установка Go

1. **Скачай:** [https://golang.org/dl/]

2. **Установи:** Запусти .msi файл

3. **Проверь:** Открой Command Prompt и введи:

`go version`

## 🚀 Настройка проекта
1. **Скачай проект**
git clone [https://github.com/твой-username/mfp.git]
cd mfp
2. **Настрой базу данных**
`createdb mybank`

`psql -U postgres -d mybank -f migrations/001_create_tables.sql`
`psql -U postgres -d mybank -f migrations/002_add_indexes.sql`
3. **Запусти сервер**

`go run cmd/server/main.go`
Успешный запуск: Server started at [http://localhost:8080]

# 🛠 Использование API
## 📝 Основные понятия
### HTTP Методы:
1. GET - получить данные (например, баланс)

2. POST - создать что-то (например, аккаунт)

3. DELETE - удалить что-то (например, аккаунт)

### Типы данных:
1. Все запросы в формате JSON

2. Все ответы в формате JSON

### Авторизация:
1. После входа система дает session_id в cookies

2. Этот session_id автоматически отправляется в каждом запросе

## 📮 Примеры запросов
1. 🆕 Регистрация нового аккаунта
**Метод:** POST
**URL:** [http://localhost:8080/register]

**Тело запроса:**

`json`
{
  "first_name": "Иван",
  "phone": "77001234567",
  "password": "1234",
  "age": 25
}


2. 🔐 Вход в систему
**Метод:** POST
**URL:** [http://localhost:8080/login]

**Тело запроса:**

`json`
{
  "phone": "77001234567",
  "password": "1234"
}

3. 💰 Пополнение баланса
**Метод:** POST
**URL:** [http://localhost:8080/accounts/me/deposit?amount=1000]

**Curl команда:**

`curl -X POST "[http://localhost:8080/accounts/me/deposit?amount=1000]" \`
 ` -H "Content-Type: application/json" `

4. 🏧 Снятие денег
**Метод:** POST
**URL:** [http://localhost:8080/accounts/me/withdraw?amount=500]

**Curl команда:**

`curl -X POST "http://localhost:8080/accounts/me/withdraw?amount=500" \`
 ` -H "Content-Type: application/json"`

5. 🔄 Перевод другому пользователю
**Метод:** POST
**URL:** [http://localhost:8080/accounts/me/transfer]

**Тело запроса:**

`json`
`{`
`  "to": "77009876543",`
  `"amount": 300`
`}`

**Curl команда:**

`curl -X POST http://localhost:8080/accounts/me/transfer \ `
  `-H "Content-Type: application/json" \`
 ` -d '{"to":"77009876543","amount":300}'`

6. 📊 Просмотр моего аккаунта
**Метод:** GET
**URL:** [http://localhost:8080/accounts/me]

**Curl команда:**

`curl http://localhost:8080/accounts/me`

7. 📋 История транзакций
**Метод:** GET
**URL:** [http://localhost:8080/accounts/me/transactions]

**Curl команда:**

`curl http://localhost:8080/accounts/me/transactions`