package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vivek1964/insta_safe_project/models"
)

var transactions []models.Transaction

func Transactions(c *fiber.Ctx) error {
	if c.Method() == "POST" {
		user_data := new(models.Transaction)

		if err := c.BodyParser(user_data); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"message": "fields are not parsable",
			})
		}
		converted_amount, err := strconv.ParseFloat(user_data.Amount, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "cannot convert amount remove",
			})
		}
		if converted_amount <= 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "amount<=0 remove",
			})
		}

		_, err = time.Parse(time.RFC3339, user_data.Timestamp)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "invalid timestamp remove",
			})
		}
		if time.Now().Format(time.RFC3339) < user_data.Timestamp {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"message": "transaction date is in the future",
			})
		}

		date := time.Now().Add(time.Second * -60).Format(time.RFC3339)
		switch {
		case date < user_data.Timestamp:
			transactions = append(transactions, models.Transaction{Amount: user_data.Amount, Timestamp: user_data.Timestamp})
			c.Status(fiber.StatusCreated)
			return c.JSON(fiber.Map{})

		case date >= user_data.Timestamp:
			c.Status(fiber.StatusNoContent)
			return c.JSON(fiber.Map{})
		}
	} else if c.Method() == "DELETE" {

		transactions = nil
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{})
	}

	return c.JSON(fiber.Map{
		"message": "Invalid",
	})
}

func Statistics(c *fiber.Ctx) error {
	count := 0
	var min float64
	var max float64
	var sum float64
	var avg float64
	for i := 0; i < len(transactions); i++ {
		if time.Now().Add(time.Second*-60).Format(time.RFC3339) < transactions[i].Timestamp {

			converted_amount, err := strconv.ParseFloat(transactions[i].Amount, 64)
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": "cannot convert amount",
				})
			}
			if count == 0 {
				min = converted_amount
				max = converted_amount
			} else if converted_amount < min {
				min = converted_amount
			}

			if converted_amount > max {
				max = converted_amount
			}
			sum += converted_amount
			count++
		}
	}
	if count != 0 {
		avg = sum / float64(count)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"sum":   fmt.Sprintf("%.2f", sum),
		"avg":   fmt.Sprintf("%.2f", avg),
		"max":   fmt.Sprintf("%.2f", max),
		"min":   fmt.Sprintf("%.2f", min),
		"count": fmt.Sprintf("%d", count),
	})
}
