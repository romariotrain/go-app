#!/bin/bash

# Скрипт для запуска веб-приложения "Блог"
# Запускает бэкенд на Go и фронтенд на простом HTTP-сервере

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Запуск веб-приложения 'Блог'${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Проверяем, установлен ли Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go не установлен!${NC}"
    echo "Установите Go перед запуском: apt-get install golang-go"
    exit 1
fi

echo -e "${GREEN}✓ Go найден: $(go version)${NC}"
echo ""

# Переходим в директорию проекта
cd "$(dirname "$0")"

# Запуск бэкенда
echo -e "${BLUE}📦 Запуск бэкенда...${NC}"
cd backend

# Убиваем старый процесс на порту 8080, если он есть
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "Останавливаем старый процесс на порту 8080..."
    kill $(lsof -t -i:8080) 2>/dev/null || true
    sleep 1
fi

# Запускаем бэкенд в фоне
nohup go run main.go > /tmp/backend.log 2>&1 &
BACKEND_PID=$!

# Даём серверу время на запуск
sleep 2

# Проверяем, что бэкенд запустился
if ! curl -s http://localhost:8080/ > /dev/null; then
    echo -e "${RED}❌ Не удалось запустить бэкенд${NC}"
    echo "Проверьте логи: cat /tmp/backend.log"
    exit 1
fi

echo -e "${GREEN}✓ Бэкенд запущен на http://localhost:8080${NC}"
echo "  Логи: /tmp/backend.log"
echo "  PID: $BACKEND_PID"
echo ""

# Запуск фронтенда
echo -e "${BLUE}🌐 Запуск фронтенда...${NC}"
cd ../frontend

# Убиваем старый процесс на порту 3000, если он есть
if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null ; then
    echo "Останавливаем старый процесс на порту 3000..."
    kill $(lsof -t -i:3000) 2>/dev/null || true
    sleep 1
fi

# Запускаем фронтенд через Python HTTP-сервер
nohup python3 -m http.server 3000 > /tmp/frontend.log 2>&1 &
FRONTEND_PID=$!

sleep 2

echo -e "${GREEN}✓ Фронтенд запущен на http://localhost:3000${NC}"
echo "  Логи: /tmp/frontend.log"
echo "  PID: $FRONTEND_PID"
echo ""

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}🎉 Приложение успешно запущено!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "📝 Доступные сервисы:"
echo "   • Веб-интерфейс: http://localhost:3000"
echo "   • API:           http://localhost:8080"
echo ""
echo "Для остановки сервисов используйте:"
echo "   kill $BACKEND_PID $FRONTEND_PID"
echo ""
echo "Или создайте файл stop.sh:"
echo "   echo '#!/bin/bash' > stop.sh"
echo "   echo 'kill $BACKEND_PID $FRONTEND_PID' >> stop.sh"
echo "   chmod +x stop.sh"
echo ""