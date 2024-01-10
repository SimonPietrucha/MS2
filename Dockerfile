# Erstellen des Image
FROM golang:latest

# Arbeitsverzeichnis im Container
WORKDIR /app

# kopieren um Abhängigkeiten zu installieren
COPY go.mod .
COPY go.sum .

# Abhängigkeiten herunterladen
RUN go mod download

# Kopieren des restlichen Projektcodes
COPY . .

# Kompilieren der Anwendung
RUN go build -o main .

# Ausführen der Anwendung
CMD ["./main"]
