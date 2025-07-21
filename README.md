# Marketplace REST API

REST API для маркетплейса с авторизацией, регистрацией, размещением объявлений и лентой объявлений.

## Технологии
- **Язык**: Go 1.23
- **Фреймворк**: Gorilla Mux
- **База данных**: PostgreSQL
- **Аутентификация**: JWT
- **Упаковка**: Docker

## Функционал
- ✅ Регистрация пользователей
- ✅ Авторизация с выдачей JWT токена
- ✅ Размещение объявлений (только для авторизованных)
- ✅ Лента объявлений с фильтрацией и сортировкой
- ✅ Проверка принадлежности объявления текущему пользователю

## Запуск проекта

### Требования
- Docker
- Docker Compose

### Инструкция
1. Клонировать репозиторий:
```bash
   git clone https://github.com/yourusername/marketplace.git
   cd marketplace
 ```

2. Создать файл .env (на основе example.env):
PORT=8080
POSTGRES_USER=vk
POSTGRES_PASSWORD=vdemok
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
POSTGRES_DB=postgres
JWT_SECRET=your_strong_secret_here

3. Запустить сервисы:
```bash
docker-compose up --build
```

### 1. 📝 Регистрация Пользователя

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login":"newuser", "password":"StrongPass123"}'
```

### 2. 🔐 Авторизация

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"login":"newuser", "password":"StrongPass123"}'
```

### 3. 📢 Создание Объявления

```bash
curl -X POST http://localhost:8080/ads \
  -H "Authorization: Bearer <Токен_пользователя>"\
  -H "Content-Type: application/json; charset=utf-8" \
  -d '{
    "title": "Новый макбук",
    "description": "Отличный макбук с 512 SSD и 16 ГБ ОЗУ",
    "image_url": "https://example.com/laptop.jpg",
    "price": 45000.50
  }'
```


### 4. 📄 Получение Ленты Объявлений

```bash
curl http://localhost:8080/ads
```
### Параметры ленты объявлений
Параметр	Описание	По умолчанию
page	Номер страницы	1
size	Количество объявлений на странице	10 (max 100)
sort_by	Поле сортировки (created_at, price)	created_at
sort_order	Направление сортировки (ASC, DESC)	DESC
min_price	Минимальная цена	0
max_price	Максимальная цена	без ограничений

### Структура прокета
```text
marketplace/
├── database/        # Работа с базой данных
├── handlers/        # Обработчики запросов
├── internal/        # Внутренние структуры данных
│   ├── ad/          # Модели объявлений
│   └── user/        # Модели пользователей
├── middleware/      # Промежуточное ПО
├── migrations/      # Миграции базы данных
├── .env             # Переменные окружения
├── docker-compose.yml # Конфигурация Docker
├── Dockerfile       # Сборка приложения
├── go.mod           # Зависимости Go
└── main.go          # Точка входа
```

### Валидация данных
##Регистрация
Логин: 3-30 символов (только буквы и цифры)

Пароль: мин. 8 символов, 1 заглавная буква, 1 цифра

##Объявление
Заголовок: 3-100 символов

Описание: 10-1000 символов

Цена: положительное число

Изображение: URL с расширением .jpg, .jpeg, .png, .gif

## Технические особенности
Пароли хранятся в виде bcrypt-хешей

JWT токены с сроком действия 24 часа

Автоматические миграции PostgreSQL при запуске

Пагинация и фильтрация на стороне БД

Оптимизированные запросы к базе данных
