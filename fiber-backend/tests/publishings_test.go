package tests

import (
	"bytes"
	"fiber-backend/router"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestTop9(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	req := httptest.NewRequest("GET", "/publishings/top9", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestGetPublishing(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings/%d", 20)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, er := appTest.Test(req, -1)
	if er != nil {
		t.Fatal(er)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestByCategories(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings/%s?page=%d&page_size=%d", "Technology", 1, 2)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestSubcategories(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings/%s/%s?page=%d&page_size=%d", "Technology", "Smartphone", 1, 4)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestFilters(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings/filters?category=%s&subcategory=%s&min_price=%f&max_price=%f&page_size=%d&page=%d",
		"Technology", "LapTop", 900.0, 1100.0, 2, 1)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestCountCategories(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	req := httptest.NewRequest("GET", "/publishings/ccount", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	

}

func TestCountSubcategories(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings/scount/%s", "Technology")
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestSaveView(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	payload := []byte(`{
		"publishing_id":3,
		"category":"Fashion",
		"subcategory":"Man",
		"count":1
	}
	`)
	req := httptest.NewRequest("POST", "/publishings/views", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}

}

func TestCreateComment(t *testing.T) {
	initConnection()
	token, _ := doLogin("userclient3@example.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	payload := []byte(`{
		"publishing_id":2,
		"rating":4,
		"description":"User test case is speaking",
		"email":"userclient3@example.com",
		"name":"userClient3-test",
		"owner_id":162
	}`)
	req := httptest.NewRequest("POST", "/publishings/comment", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	

}

func TestTop9ByCategory(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings/top9-category/%s?start_date=%s&end_date=%s", "Technology","2023-08-10","2023-08-26")
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestTop9BySubCategory(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings/top9-subcategory/%s?start_date=%s&end_date=%s", "Motherboard","2023-08-10","2023-08-26")
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestGetAll(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/publishings?page=%d&page_size=%d", 1,4)
	req := httptest.NewRequest("GET", path, nil)
	resp, err :=  appTest.Test(req,-1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestGetTopSeller(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	req := httptest.NewRequest("GET", "/publishings/top4-sellers", nil)
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
}