package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gte=0"`
	//SKU         string  `json:"sku" validate:"required, sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

type Products []*Product

var ErrProductNotFound = fmt.Errorf("product not found")

func (p *Product) Validate() error {
	validate := validator.New()
	//validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

//func validateSKU(fl validator.FieldLevel) bool {
//	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
//	matches := re.FindAllString(fl.Field().String(), -1)
//
//	if len(matches) != 1 {
//		return false
//	}
//
//	return true
//}

func FromJSON(i interface{}, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(i)
}

func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

func GetProductById(id int) (*Product, error) {
	i := findIndexByProductId(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}
	return productList[i], nil
}

func AddProduct(p *Product) {
	p.ID = generateNextId()
	productList = append(productList, p)
}

func UpdateProduct(p Product) error {
	pos := findIndexByProductId(p.ID)
	if pos == -1 {
		return ErrProductNotFound
	}
	productList[pos] = &p
	return nil
}

func findIndexByProductId(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func generateNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

// productList is a hard coded list of products for this example data source
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Made with espresso and steamed milk.",
		Price:       2.99,
		//SKU:         "abc323",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Mocaccino",
		Description: "A chocolate-flavoured warm beverage that is a variant of a caffè latte",
		Price:       1.99,
		//SKU:         "fjd34",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
}
