package routes

import (
	"errors"
	"fmt"

	"github.com/Jin1iangYan/fiber-gorm/database"
	"github.com/Jin1iangYan/fiber-gorm/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	// This is not the model User, see this as the serializer
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// CreateResponseUser creates a user serializer.
func CreateResponseUser(userModel models.User) User {
	return User{
		ID:        userModel.ID,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
	}
}

// CreateUser creates a user and save back into database
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

// GetUsers gets all users from the database
func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users)

	responseUsers := []User{}

	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

// FindUser is a helper function finds user by id
func findUser(id string, user *models.User) error {
	database.Database.Db.Find(user, "id = ?", id)
	if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("user does not exist")
	}
	return nil
}

// GetUser gets the user from database by id
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	fmt.Println(user)

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

// UpdateUser updates the user information
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	fmt.Println(user)

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

// UpdateUser deletes the user
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	fmt.Println(user)

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Delete User")
}
