FROM golang:1.24

WORKDIR /app

COPY go.mod .
COPY main.go .
COPY .env .

#para que se instalen los paquetes
RUN go get github.com/joho/godotenv
RUN go get -u github.com/gin-gonic/gin 
RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/mysql


RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]