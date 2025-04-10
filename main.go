package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database = func() (db *gorm.DB) {
	//cargamos las variables de entorno
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic(errorVariables)
	}
	//domain server name

	//dsn := "pepe:lorito2004@tcp(localhost:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println("DB_SERVER: ", dsn)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Error de conexión")
		panic(err)
	} else {
		fmt.Println("Conexión a mysql exitosa")
		return db
	}
}()

type Categoria struct { //este es el que se tiene migrar, ojo pesteña ceja
	Id     uint   `json:"id"` //bigint
	Nombre string `gorm:"type:varchar(100)" json:"nombre"`
	Slug   string `gorm:"type:varchar(100)" json:"slug"`
}

type Categorias []Categoria

func Migraciones() {
	Database.AutoMigrate(&Categoria{})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//cors
	router.Use(corsMiddleware())
	//ejecutamos las migraciones
	Migraciones()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"estado":  "ok",
			"mensaje": "Hola desde Docker as",
		})
	})

	router.GET("/categorias", func(c *gin.Context) {
		datos := Categorias{}
		Database.Order("id desc").Find(&datos) //select id, nombre, slug from categoria
		c.JSON(http.StatusOK, datos)
	})

	router.Run(":" + os.Getenv("PORT"))

}
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

//go get -u github.com/gin-gonic/gin
//go clean -i github.com/gin-gonic/gin

//docker build . -t ejemplo-go-gin-docker
//docker run -e PORT=8080 -p 8080:8080 ejemplo-go-gin-docker

// docker compose up --build
