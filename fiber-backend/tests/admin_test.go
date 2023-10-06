package tests

import (
	"bytes"
	"encoding/json"
	"fiber-backend/models"
	"fiber-backend/router"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

)

func helperSearchUser(token, username string) (models.Users, error) {
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/admin/search-users?username=%s&page=1&page_size=5", username)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := appTest.Test(req, -1)

	body, _ := io.ReadAll(resp.Body)
	var user = struct {
		Users []models.Users `json:"users"`
	}{}

	var result models.Users
	if err := json.Unmarshal(body, &user); err != nil {
		return user.Users[0], err
	}
	for i := 0; i < len(user.Users); i++ {
		if user.Users[i].Username == username {
			result = user.Users[i]
		}
	}
	return result, nil

}

func helperGetViews() (models.Views, error) {
	var result = struct{
		Views []models.Views `json:"views"`
	}{}
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/get-views?page=1&page_size=50&category=Fashion", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		return models.Views{}, err
	}

	payload, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return models.Views{}, err
	}

	if err = json.Unmarshal(payload, &result); err != nil {
		return models.Views{}, err
	}
	return result.Views[len(result.Views) - 1],nil
}

func TestLogin(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	payload := []byte(`{"email":"admin@admin.com","password":"admin1234567"}`)
	req := httptest.NewRequest("POST", "/admin/signin", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 180 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 180)
	}

}

func TestGetUsers(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/get-users?page=1&page_size=5", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 200 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 200)
	}
}

func TestSearchUsers(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/search-users?username=user&page=1&page_size=5", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	if resp.ContentLength < 200 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 200)
	}
}

func TestCreateUser(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	username, _ := w.CreateFormField("username")
	username.Write([]byte("TestCaseUserAdmin"))
	email, _ := w.CreateFormField("email")
	email.Write([]byte("testcaseadmin@testcase.com"))
	tel, _ := w.CreateFormField("tel")
	tel.Write([]byte("+12345678"))
	pass, _ := w.CreateFormField("password")
	pass.Write([]byte("1234567abc"))
	role, _ := w.CreateFormField("role")
	role.Write([]byte("client"))
	w.Close()

	req := httptest.NewRequest("POST", "/admin/create-users", &body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}

}

func TestEditUser(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	user, err := helperSearchUser(token,"TestCaseUserAdmin")
	if err != nil {
		t.Fatal(err)
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	id, _ := w.CreateFormField("user_id")
	id.Write([]byte(fmt.Sprintf("%d", user.ID)))
	address, _ := w.CreateFormField("address")
	address.Write([]byte("Test address around the world"))
	w.Close()

	req := httptest.NewRequest("PATCH", "/admin/update-users", &body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}

}


func TestAdminGetPublishings(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/get-publishings?page=1&page_size=5", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	if resp.ContentLength < 150 {
		t.Errorf("Content length expected %d and recived %d", 150, resp.ContentLength)
	}

}

func TestAdminCreatePublishing(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	users , err := helperSearchUser(token,"TestCaseUserAdmin")
	if err != nil {
		t.Fatal(err)
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	userid, _ := w.CreateFormField("user_id")
	userid.Write([]byte(fmt.Sprintf("%d", users.ID)))
	category,_ := w.CreateFormField("category")
	category.Write([]byte("Fashion"))
	subca, _ := w.CreateFormField("subcategory")
	subca.Write([]byte("Man"))
	productname,_ := w.CreateFormField("product_name")
	productname.Write([]byte("TestProductName"))
	price, _ := w.CreateFormField("price")
	price.Write([]byte(fmt.Sprintf("%f", 40.0)))
	w.Close()
	req := httptest.NewRequest("POST", "/admin/create-publishing", &body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type",w.FormDataContentType())

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
}

func TestAdminUpdatePublishing(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	path := "/admin/get-publishings?product_name=TestProductName&page=1&page_size=5"
	resp, err := getPublishings(path, token)
	if err != nil {
		t.Fatal(err)
	}
	var mapRsp = struct {
		Publishings []models.Publishings `json:"publishings"`
	}{}
	
	body,_ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &mapRsp)
	if err != nil {
		t.Fatal(err)
	}
	var pubId int
	for i := 0; i < len(mapRsp.Publishings); i++ {
		if mapRsp.Publishings[i].ProductName == "TestProductName" {
			pubId = int(mapRsp.Publishings[i].ID)
		}
	}
	var payload bytes.Buffer
	w := multipart.NewWriter(&payload)
	publishingId,_ := w.CreateFormField("publishing_id")
	publishingId.Write([]byte(fmt.Sprintf("%d", pubId)))
	inOffer, _ := w.CreateFormField("in_offer")
	inOffer.Write([]byte("30"))
	w.Close()

	req := httptest.NewRequest("PATCH", "/admin/update-publishing", &payload)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp2, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
} 

func TestAdminDeletePublishing(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	path := "/admin/get-publishings?product_name=TestProductName&page=1&page_size=5"
	resp, err := getPublishings(path, token)
	if err != nil {
		t.Fatal(err)
	} 
	var mapRsp = struct {
		Publishings []models.Publishings `json:"publishings"`
	}{}

	var pubId int
	body,_ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &mapRsp)
	if err != nil {
		t.Fatal(err)
	}
	
	for i := 0; i < len(mapRsp.Publishings); i++ {
		if mapRsp.Publishings[i].ProductName == "TestProductName" {
			pubId = int(mapRsp.Publishings[i].ID)
		}
	}
	var payload = make(map[string]int)
	payload["publishing_id"] = pubId

	data,_ := json.Marshal(payload)

	req2 := httptest.NewRequest("DELETE", "/admin/delete-publishing", bytes.NewBuffer(data))
	req2.Header.Set("Authorization", token)
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := appTest.Test(req2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
}


func TestAdminGetByRole(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/get-byrole", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	} 
}

func TestAdminGetByCategory(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/get-bycategory", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	} 
}

func TestAdminGetBySubcategory(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/get-bysubcategory/:Technology", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	} 
	
}


func TestGetViews(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	req := httptest.NewRequest("GET", "/admin/get-views?page=1&page_size=10&category=Fashion", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	if resp.ContentLength < 100 {
		t.Errorf("Content length expected %d and recived %d", 100, resp.ContentLength)
	}
}


func TestAdminSaveView(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	body := make(map[string]interface{})
	body["publishing_id"] = 6
	body["category"] = "Fashion"
	body["subcategory"] = "Woman"
	body["count"] = 1
	payload, err := json.Marshal(body)

	req := httptest.NewRequest("POST", "/admin/save-views", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}

}


func TestDeleteViews(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")

	view, err := helperGetViews()
	if err != nil {
		t.Fatal(err)
	}

	body := make(map[string]int)
	body["id"] = int(view.ID)
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("DELETE", "/admin/delete-views", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}

}


func TestCleanUser(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, _ := doLogin("admin@admin.com", "admin1234567","admin")
	user, err := helperSearchUser(token,"TestCaseUserAdmin")
	if err != nil {
		t.Fatal(err)
	}

	body := make(map[string]int)
	body["user_id"] = int(user.ID)

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("DELETE", "/admin/delete-users", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}

}
