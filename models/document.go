package models

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ivangurin/restful-api-go/database"
	"github.com/timsolov/rest-query-parser"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	Number      string
	Description string
	Date        time.Time
	Items       []*DocumentItem
}

type DocumentItem struct {
	gorm.Model
	DocumentID  uint
	Number      int
	Description string
	Quantity    float64
	Unit        string
	Price       float64
	Currency    string
}

type DocumentResponse struct {
	ID          uint                    `json:"id"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
	Number      string                  `json:"number"`
	Description string                  `json:"description"`
	Total       float64                 `json:"total"`
	Date        time.Time               `json:"date"`
	Items       []*DocumentItemResponse `json:"items"`
}

type DocumentItemResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Number      int       `json:"number"`
	Description string    `json:"description"`
	Quantity    float64   `json:"quantity"`
	Unit        string    `json:"unit"`
	Price       float64   `json:"price"`
	Value       float64   `json:"value"`
	Currency    string    `json:"currency"`
}

func (document *Document) GetDocumentResponse() DocumentResponse {

	documentResponse := DocumentResponse{
		ID:          document.ID,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
		Number:      document.Number,
		Description: document.Description,
		Date:        document.Date,
	}

	for _, item := range document.Items {

		itemResponse := &DocumentItemResponse{
			ID:          item.ID,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			Number:      item.Number,
			Description: item.Description,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			Price:       item.Price,
			Value:       float64(int(item.Quantity * item.Price * 100))/100,
			Currency:    item.Currency,
		}

		documentResponse.Items = append(documentResponse.Items, itemResponse)

		documentResponse.Total += itemResponse.Value

	}

	return documentResponse

}

func (document *Document) SetDocumentResponse(documentResponse *DocumentResponse) (err error) {

	document.Number = documentResponse.Number
	document.Description = documentResponse.Description
	document.Date = documentResponse.Date

	for _, item := range document.Items {

		itemResponse := &DocumentItemResponse{}

		for _, itemResponse = range documentResponse.Items {

			if item.ID == itemResponse.ID {
				break
			}

		}

		if itemResponse.ID != 0 {

			item.Number = itemResponse.Number
			item.Description = itemResponse.Description
			item.Quantity = itemResponse.Quantity
			item.Unit = itemResponse.Unit
			item.Price = itemResponse.Price
			item.Currency = itemResponse.Currency

		} else {
			item.DeletedAt.Time = time.Now()
		}

	}

	for _, itemResponse := range documentResponse.Items {

		if itemResponse.ID != 0 {
			continue
		}

		item := &DocumentItem{}

		item.Number = itemResponse.Number
		item.Description = itemResponse.Description
		item.Quantity = itemResponse.Quantity
		item.Unit = itemResponse.Unit
		item.Price = itemResponse.Price
		item.Currency = itemResponse.Currency

		document.Items = append(document.Items, item)

	}

	return

}

func (document *Document) Save() (err error) {

	if document.ID == 0 {

		err = database.Db.Create(document).Error

		if err != nil {
			return
		}

	} else {

		itemsForDelete := []*DocumentItem{}

		for _, item := range document.Items {
			if !item.DeletedAt.Time.IsZero() {
				itemsForDelete = append(itemsForDelete, item)
			}
		}

		//err = database.Db.Save(&document).Error
		err = database.Db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&document).Error

		if err != nil {
			return
		}

		if len(itemsForDelete) > 0 {

			err = database.Db.Delete(&itemsForDelete).Error

			if err != nil {
				return
			}

		}

	}

	*document, err = GetDocumentById(int(document.ID))

	if err != nil {
		return
	}

	return

}

func GetDocuments(uri string) []Document {

	fmt.Println(uri)

	url, _ := url.Parse(uri)

	q, _ := rqp.NewParse(url.Query(), rqp.Validations{
		"limit": nil,
		"offset": nil,
		"sort": rqp.In("id", "number", "description"),
		"id:int": nil,
	})

	fmt.Println(q.SQL("table"))
	fmt.Println(q.Where())
	fmt.Println(q.Args())
	fmt.Println(q.Limit)
	fmt.Println(q.Offset)
	fmt.Println(q.Sorts)

	documnets := []Document{}

	database.Db.Where(q.Where(), q.Args()...).Limit(q.Limit).Offset(q.Offset).Preload("Items").Find(&documnets)

	return documnets

}

func GetDocumentById(id int) (document Document, err error) {

	err = database.Db.Preload("Items").Find(&document, id).Error

	if err != nil {
		return
	}

	if document.ID == 0 {
		err = errors.New("Document not found")
		return
	}

	return

}

func DeleteDocumentById(id int) (err error) {

	document, err := GetDocumentById(id)

	if err != nil {
		return
	}

	err = database.Db.Delete(&document).Error

	if err != nil {
		return
	}

	err = database.Db.Delete(&document.Items).Error

	if err != nil {
		return
	}

	return

}
