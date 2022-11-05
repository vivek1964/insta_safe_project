package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/assert"
	"github.com/vivek1964/insta_safe_project/controllers"
)

func TestTransactionsRoute(t *testing.T) {
	tests := []struct {
		Amount       string `json:"amount"`
		Description  string `json:"description"`
		ExpectedCode int    `json:"expectedCode"`
		Timestamp    string `json:"timestamp"`
	}{
		{
			Description:  "case of success",
			Amount:       "10.3",
			Timestamp:    "2022-11-05T09:44:00.312Z",
			ExpectedCode: 201,
		},
		{
			Description:  " transaction is older than 60 seconds",
			Amount:       "10.3",
			Timestamp:    "2022-11-05T01:09:00.312Z",
			ExpectedCode: 204,
		},
		{
			Description:  "Json is invalid",
			Amount:       "",
			Timestamp:    "",
			ExpectedCode: 400,
		},
		{
			Description:  "transaction date is in the future",
			Amount:       "10.3",
			Timestamp:    "2023-11-05T01:03:00.312Z",
			ExpectedCode: 422,
		},
	}

	// Define Fiber app.
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	// Create route
	app.Get("/statistics", controllers.Statistics)
	app.Post("/transactions", controllers.Transactions)
	app.Delete("/transactions", controllers.Transactions)

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case

		var requestBodyMap map[string]string = map[string]string{
			"amount":    test.Amount,
			"timestamp": test.Timestamp,
		}

		requestBodyjson, err := json.Marshal(requestBodyMap)
		if err != nil {
			assert.Equalf(t, test.ExpectedCode, 0, "error converting requestbody to json")
		}
		requestBodyByte := bytes.NewBuffer(requestBodyjson)
		req := httptest.NewRequest("POST", "/transactions", requestBodyByte)
		req.Header.Add("Content-Type", "application/json")
		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, err := app.Test(req, -1)
		if err != nil {
			assert.Equalf(t, test.ExpectedCode, 0, fmt.Sprintf("Unable to send request, %q", err.Error()))
		} else {
			// Verify, if the status code is as expected
			assert.Equalf(t, test.ExpectedCode, resp.StatusCode, test.Description)
		}

	}
type Response struct {
	Sum         string `json:"sum"`
	Avg         string `json:"avg"`
	Count       string `json:"count"`
	Max         string `json:"max"`
	Min         string `json:"min"`
}
var test_statistics [] Response
test_statistics = append(test_statistics, Response{	Sum: "10.30", Avg: "10.30", Count: "1", Max:  "10.30", Min:  "10.30",})


	// Iterate through test single test cases
	for _, test := range test_statistics {
		// Create a new http request with the route from the test case

		req := httptest.NewRequest("GET", "/statistics", nil)
		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, -1)
		var response Response
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equalf(t, test.Sum, response.Sum, "Mis match in sum")
		assert.Equalf(t, test.Min, response.Min, "Mis match in Min")
		assert.Equalf(t, test.Max, response.Max, "Mis match in Max")
		assert.Equalf(t, test.Count, response.Count, "Mis match in Count")
		assert.Equalf(t, test.Avg, response.Avg, "Mis match in Avg")

	}

	type Response2 struct {
		Description  string `json:"description"`
		ExpectedCode int    `json:"expectedCode"`
	}
	var test_transactions_delete [] Response2
	test_transactions_delete = append(test_transactions_delete, Response2{	Description: "Statistics000", ExpectedCode: 204})
		// Iterate through test single test cases
		for _, test := range test_transactions_delete {
			// Create a new http request with the route from the test case
	
			req := httptest.NewRequest("DELETE", "/transactions", nil)
			// Perform the request plain with the app,
			// the second argument is a request latency
			// (set to -1 for no latency)
			resp, _ := app.Test(req, -1)
			assert.Equalf(t, test.ExpectedCode, resp.StatusCode, test.Description)
	
		}
	
}
