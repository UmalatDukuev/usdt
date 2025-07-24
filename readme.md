USDT Rate GRPC Service


Установка и запуск

git clone https://github.com/UmalatDukuev/usdt.git
cd usdt
make build
docker-compose up -d


Запуск вручную (без Docker):
DB_URL="postgres://postgres:pass@localhost:5432/dbname?sslmode=disable" \
API_URL="https://grinex.io/api/v2/depth" \
PORT=50051 \
./usdt-app


3 уровня конфигурации
1. флаги запуска
2. переменные окружения
3. yml-конфиг

