# Verwenden Sie das offizielle Go-Image als Basis
FROM golang:latest

# Setzen Sie das Arbeitsverzeichnis innerhalb des Containers
WORKDIR /app

# Kopieren Sie die go.mod- und go.sum-Dateien, um die Abhängigkeiten zu installieren
COPY go.mod .
COPY go.sum .

# Führen Sie 'go mod download' aus, um die Abhängigkeiten herunterzuladen
RUN go mod download

# Kopieren Sie den restlichen Projektcode in das Arbeitsverzeichnis
COPY . .

# Kompilieren Sie die Go-Anwendung
RUN go build -o main .

# Setzen Sie den Befehl zum Ausführen der Anwendung
CMD ["./main"]
