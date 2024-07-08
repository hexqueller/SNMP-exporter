FROM python:3.11-slim

# Устанавливаем необходимые библиотеки и зависимости
RUN apt-get update && apt-get install -y \
    libsnmp-dev \
    build-essential \
    && pip install pysnmp

# Создаем и устанавливаем рабочий каталог
WORKDIR /app

# Копируем файлы проекта в рабочий каталог
COPY proxy /app/proxy
COPY configs/default.yaml /app/default.yaml
COPY agent.py /app/agent.py

# Указываем команду запуска
CMD ["/app/proxy", "-c", "/app/default.yaml"]
