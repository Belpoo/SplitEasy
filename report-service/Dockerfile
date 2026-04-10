# Archivo: Dockerfile para correr entorno Serverless emulado nativamente en Render WS
FROM python:3.11-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
EXPOSE 8080

# En modo Function Framework, emulamos la llamada
CMD ["functions-framework", "--target=generate_report", "--port=8080"]
