package routes

import (
	"errors"

	"github.com/Jin1iangYan/fiber-gorm/database"
	"github.com/Jin1iangYan/fiber-gorm/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serial_number"`
}

func CreateResponseProduct(product models.Product) Product {
	return Product{
		ID:           product.ID,
		Name:         product.Name,
		SerialNumber: product.SerialNumber,
	}
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&product)

	responseProduct := CreateResponseProduct(product)

	return c.Status(201).JSON(responseProduct)
}

func GetProducts(c *fiber.Ctx) error {
	products := []models.Product{}

	database.Database.Db.Find(&products)

	responseProducts := []Product{}

	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}

	return c.Status(200).JSON(responseProducts)
}

func findProduct(id string, product *models.Product) error {
	if err := database.Database.Db.Find(&product, "id = ?", id).Error; err != nil {
		return err
	}
	if product.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("product does not exist")
	}
	return nil
}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product

	id := c.Params("id")

	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseProduct := CreateResponseProduct(product)
	return c.Status(200).JSON(responseProduct)
}

func UpdateProduct(c *fiber.Ctx) error {
	var product models.Product

	id := c.Params("id")

	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateProduct struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}
	updateProduct := UpdateProduct{}

	if err := c.BodyParser(&updateProduct); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	product.Name = updateProduct.Name
	product.SerialNumber = updateProduct.SerialNumber

	database.Database.Db.Save(&product)

	responseProduct := CreateResponseProduct(product)
	return c.Status(200).JSON(responseProduct)
}

func DeleteProduct(c *fiber.Ctx) error {
	var product models.Product

	id := c.Params("id")

	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Delete Product")
}
