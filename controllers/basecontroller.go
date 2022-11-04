package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserData struct {
	Amount    string `json:"amount" xml:"amount" form:"amount"`
	Timestamp string `json:"timestamp" xml:"timestamp" form:"timestamp"`
}

var transactions []UserData

func Transactions(c *fiber.Ctx) error {

	user_data := new(UserData)

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

	today := time.Now().Add(time.Second * -120).Format(time.RFC3339)
	switch {
	case today < user_data.Timestamp:
		transactions = append(transactions, UserData{Amount: user_data.Amount, Timestamp: user_data.Timestamp})
		c.Status(fiber.StatusCreated)
		return c.JSON(fiber.Map{
			"message1": fmt.Sprintf("TODAY<given(not older)     %s(current), %s", today, user_data.Timestamp),
			"message":  "less than or equal to",
		})

	case today >= user_data.Timestamp:
		//c.Status(fiber.StatusNoContent)here
		return c.JSON(fiber.Map{
			"message": "transaction is older than 60 seconds",
			"delete":  fmt.Sprintf("TODAY>given     %s(current), %s", today, user_data.Timestamp),
		})
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
					"message": "cannot convert amount remove",
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
	return c.JSON(fiber.Map{
		"sum":   sum,
		"avg":   avg,
		"max":   max,
		"min":   min,
		"count": count,
	})
}

/*
ISO 8601 Formats : 2020-07-10 15:00:00.000
in golang it is                                2022-11-04T16:51:17+05:30
*/
