package routes

import (
	"net/http"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/ivangurin/restful-api-go/models"
)

func PopulateDocument(c *fiber.Ctx) error {

	documnet := models.Document{}

	documnet.Number = strconv.Itoa(gofakeit.Number(1, 1000000))
	documnet.Description = gofakeit.Phrase()
	documnet.Date = gofakeit.Date()

	documnet.Items = append(documnet.Items, &models.DocumentItem{
		Number: 1,
		Description: gofakeit.NounCollectiveThing(),
		Quantity:    float64(gofakeit.Number(1, 100)),
		Unit:        "ШТ",
		Price:       gofakeit.Price(100, 1000),
		Currency:    gofakeit.CurrencyShort(),
	})

	documnet.Items = append(documnet.Items, &models.DocumentItem{
		Number: 2,
		Description: gofakeit.NounCollectiveThing(),
		Quantity:    float64(gofakeit.Number(1, 100)),
		Unit:        "ШТ",
		Price:       gofakeit.Price(100, 1000),
		Currency:    gofakeit.CurrencyShort(),
	})	

	if err := documnet.Save(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(documnet.GetDocumentResponse())

}

func GetDocuments(c *fiber.Ctx) error {

	documnets := models.GetDocuments(string(c.Request().URI().FullURI()))

	documnetsResponse := []models.DocumentResponse{}

	for _, document := range documnets {
		documnetsResponse = append(documnetsResponse, document.GetDocumentResponse())
	}

	return c.Status(http.StatusOK).JSON(documnetsResponse)

}

func GetDocument(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Please ensure that document id is correct")
	}

	document, err := models.GetDocumentById(id)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(document.GetDocumentResponse())

}

func CreateDocument(c *fiber.Ctx) error {

	documentResponse := models.DocumentResponse{}

	if err := c.BodyParser(&documentResponse); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	document := models.Document{}

	if err := document.SetDocumentResponse(&documentResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	if err := document.Save(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(document.GetDocumentResponse())

}

func UpdateDocument(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Please ensure that document id is correct")
	}

	document, err := models.GetDocumentById(id)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	documentResponse := models.DocumentResponse{}

	err = c.BodyParser(&documentResponse);
	
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	err = document.SetDocumentResponse(&documentResponse)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	err = document.Save()

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(document.GetDocumentResponse())

}

func DeleteDocument(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Please ensure that document id is correct")
	}

	err = models.DeleteDocumentById(id)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON("Document was deleted")

}
