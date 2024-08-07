package routes

import (
	"errors"
	"time"

	"github.com/Jin1iangYan/fiber-gorm/database"
	"github.com/Jin1iangYan/fiber-gorm/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	User      User      `json:"user"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"order_data"`
}

func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{
		ID:        order.ID,
		User:      user,
		Product:   product,
		CreatedAt: order.CreatedAt,
	}
}

func SerializeOrder(order models.Order) (Order, error) {
	var user models.User
	if err := findUser(order.UserID.String(), &user); err != nil {
		return Order{}, errors.New(err.Error())
	}
	if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return Order{}, errors.New("user does not exist")
	}

	var product models.Product
	if err := findProduct(order.ProductID.String(), &product); err != nil {
		return Order{}, errors.New(err.Error())
	}
	if product.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return Order{}, errors.New("product does not exist")
	}

	respUser := CreateResponseUser(user)
	respProduct := CreateResponseProduct(product)
	respOrder := CreateResponseOrder(order, respUser, respProduct)

	return respOrder, nil
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var user models.User
	if err := findUser(order.UserID.String(), &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(400).SendString("user does not exist")
	}

	var product models.Product
	if err := findProduct(order.ProductID.String(), &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if product.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(400).SendString("product does not exist")
	}

	database.Database.Db.Create(&order)

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(201).JSON(responseOrder)
}

func GetOrders(c *fiber.Ctx) error {
	orders := []models.Order{}

	database.Database.Db.Find(&orders)

	responseOrders := []Order{}

	for _, order := range orders {
		respOrder, err := SerializeOrder(order)
		if err != nil {
			return c.Status(400).JSON(err.Error())
		}
		responseOrders = append(responseOrders, respOrder)
	}

	return c.Status(200).JSON(responseOrders)
}

func findOrder(id string, order *models.Order) error {
	if err := database.Database.Db.Find(&order, "id = ?", id).Error; err != nil {
		return err
	}
	if order.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("order dose not exist")
	}
	return nil
}

func GetOrder(c *fiber.Ctx) error {
	var order models.Order
	id := c.Params("id")

	if err := findOrder(id, &order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseOrder, err := SerializeOrder(order)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(responseOrder)
}

// func UpdateOrder(c *fiber.Ctx) error {
// 	var order models.Order
// 	id := c.Params("id")

// 	if err := findOrder(id, &order); err != nil {
// 		return c.Status(400).JSON(err.Error())
// 	}

// 	type UpdateOrder struct {
// 		ProductID string `json:"product_id"`
// 		UserID    string `json:"user_id"`
// 	}
// 	var updateOrder UpdateOrder
// 	if err := c.BodyParser(&updateOrder); err != nil {
// 		return c.Status(500).JSON(err.Error())
// 	}

// 	order.ProductID, _ = uuid.Parse(updateOrder.ProductID)
// 	order.UserID, _ = uuid.Parse(updateOrder.UserID)

// 	database.Database.Db.Save(order)

// 	responseOrder, err := SerializeOrder(order)
// 	if err != nil {
// 		return c.Status(400).JSON(err.Error())
// 	}

// 	return c.Status(200).JSON(responseOrder)
// }

// func DeleteOrder(c *fiber.Ctx) error {
// 	var order models.Order
// 	id := c.Params("id")

// 	if err := findOrder(id, &order); err != nil {
// 		return c.Status(400).JSON(err.Error())
// 	}

// 	if err := database.Database.Db.Delete(&order).Error; err != nil {
// 		return c.Status(404).JSON(err.Error())
// 	}

// 	return c.Status(200).SendString("Successfully Delete Order")
// }
