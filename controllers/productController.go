package controllers

import (
	"log"
	"net/http"

	"github.com/AlirezaAK2000/online-shop/initializers"
	"github.com/AlirezaAK2000/online-shop/repo"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertProductController(c *gin.Context) {

	var body struct {
		Name          string
		Description   string
		Price         float64
		StockQuantity int32
		Category      string
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body.",
		})
		return
	}

	product := repo.Product{
		Name:          body.Name,
		Description:   body.Description,
		Price:         body.Price,
		StockQuantity: body.StockQuantity,
		Category:      body.Category,
	}
	_, err := product.InsertProduct()
	if err != nil {
		log.Fatal(err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"product": product,
	})

}

func GetAllProductsController(c *gin.Context) {

	result, err := repo.FindAllProducts()

	if err != nil {
		log.Fatal(err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"products": result,
	})

}

func GetProductByIDController(c *gin.Context) {

	id := c.Param("id")
	product, err := repo.FindProductByID(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	err2 := initializers.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &initializers.ProductClickTopic, Partition: kafka.PartitionAny},
		Value:          []byte(id),
	}, nil)

	if err2 != nil {
		log.Fatalln(err2)
	}

	c.JSON(200, gin.H{
		"product": product,
	})

}

func DeleteProductByIDController(c *gin.Context) {
	id := c.Param("id")
	_, err := repo.DeleteProductByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, gin.H{})
}
