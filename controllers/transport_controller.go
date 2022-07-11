package controllers

import (
	"context"
	"go-trans/configs"
	"go-trans/models"
	"go-trans/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var transCollection *mongo.Collection = configs.GetCollection(configs.DB, "transports")
var validate = validator.New()

//Create a transport
func CreateTransport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var transport models.Transport

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&transport); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TransportResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&transport); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TransportResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newTransport := models.Transport{
		Id:          primitive.NewObjectID(),
		Name:        transport.Name,
		Description: transport.Description,
		Modality:    transport.Modality,
	}

	result, err := transCollection.InsertOne(ctx, newTransport)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TransportResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.TransportResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})

}

//Get transport by ID
func GetTransport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	transportId := c.Params("transportId")
	var transport models.Transport
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(transportId)
	err := transCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&transport)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TransportResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.TransportResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": transport}})

}

//Edit a transport of given ID
func EditTransport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	transportId := c.Params("transportId")
	var transport models.Transport
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(transportId)

	//validate the request body
	if err := c.BodyParser(&transport); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TransportResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&transport); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TransportResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": transport.Name, "description": transport.Description, "modality": transport.Modality}

	result, err := transCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TransportResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated transport details
	var updatedTrans models.Transport
	if result.MatchedCount == 1 {
		err := transCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedTrans)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.TransportResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.TransportResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedTrans}})
}

//Delete a transport with given ID
func DeleteTransport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	transportId := c.Params("transportId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(transportId)

	result, err := transCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TransportResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.TransportResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Transport with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.TransportResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Transport successfully deleted!"}},
	)
}

//Return list of all transports in DB
func GetAllTransports(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var transports []models.Transport
	defer cancel()

	results, err := transCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TransportResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleTrans models.Transport
		if err = results.Decode(&singleTrans); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.TransportResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		transports = append(transports, singleTrans)
	}

	return c.Status(http.StatusOK).JSON(
		responses.TransportResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": transports}},
	)
}
