package controllers

import (
	"context"
	"dirobash-api/configs"
	"dirobash-api/models"
	"dirobash-api/responses"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var citationCollection *mongo.Collection = configs.GetCollection(configs.DB, "quotes")
var validate = validator.New()

func CreateCitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var citation models.Citation
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&citation); err != nil {
			c.JSON(http.StatusBadRequest, responses.CitationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&citation); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CitationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newCitation := models.Citation{
			Id:     primitive.NewObjectID(),
			Date:   citation.Date,
			Number: citation.Number,
			Text:   citation.Text,
		}

		result, err := citationCollection.InsertOne(ctx, newCitation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CitationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.CitationResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetACitationById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		citationId := c.Param("citationId")
		var citation models.Citation
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(citationId)

		err := citationCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&citation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CitationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.CitationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": citation}})
	}
}

func GetACitationByNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		number, paramErr := strconv.Atoi(c.Param("number"))

		if paramErr != nil {
			// ... handle error
			panic(paramErr)
		}

		var citation models.Citation
		defer cancel()

		err := citationCollection.FindOne(ctx, bson.M{"number": number}).Decode(&citation)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CitationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.CitationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": citation}})
	}
}

func GetAllCitations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var citations []models.Citation
		defer cancel()

		filter := bson.D{}
		opts := options.Find().SetSort(bson.D{{"number", -1}})

		results, err := citationCollection.Find(ctx, filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CitationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleCitation models.Citation
			if err = results.Decode(&singleCitation); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CitationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			citations = append(citations, singleCitation)
		}

		c.JSON(http.StatusOK,
			responses.CitationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": citations}},
		)
	}
}

func GetPagedCitations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		page, paramErr := strconv.ParseInt(c.Param("page"), 10, 64)

		if paramErr != nil {
			// ... handle error
			panic(paramErr)
		}

		var citations []models.Citation
		defer cancel()

		filter := bson.D{}

		opts := options.Find().SetSort(bson.D{{Key: "number", Value: -1}}).SetSkip((page - 1) * 20).SetLimit(20)

		results, err := citationCollection.Find(ctx, filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CitationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleCitation models.Citation
			if err = results.Decode(&singleCitation); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CitationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			citations = append(citations, singleCitation)
		}

		c.JSON(http.StatusOK,
			responses.CitationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": citations}},
		)
	}
}
