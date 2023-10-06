package tests

import (
	"bytes"
	"encoding/json"
	"fiber-backend/models"
	"fiber-backend/router"
	"fiber-backend/settings"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func initConnection() {
	op := settings.Settings{}.GetEnvVarTest()
	models.DialDb(op.DSN, "./../logs/error.log")
}

func doLogin(email, password, prex string) (string, int) {
	appTest := fiber.New()
	router.Router(appTest)
	var user = struct {
		ID    int    `json:"user_id"`
		Image string `json:"image"`
		Token string `json:"session_token"`
	}{}
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{email, password}
	payload, _ := json.Marshal(&data)
	path := fmt.Sprintf("/%s/signin", prex)
	req := httptest.NewRequest("POST", path, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&user); err != nil {
		panic(err)
	}
	return user.Token, user.ID
}

func doCheckCode(body []byte) int {
	appTest := fiber.New()
	router.Router(appTest)
	req := httptest.NewRequest("POST", "/users/check-code", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := appTest.Test(req, -1)
	return resp.StatusCode
}

func getPublishings(path,token string) (*http.Response,error) {
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	appTest := fiber.New()
	router.Router(appTest)
	resp, err := appTest.Test(req, -1)
	return resp, err
}

func TestSignUp(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	payload := []byte(`{"username":"userTestCase","email":"usertestcase@test.com","password":"1234567abc","tel":"+5353464587"}`)
	req := httptest.NewRequest("POST", "/users/signup", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 70 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 70)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if status := doCheckCode(respBody); status != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, status)
	}

}

func TestSignIn(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	payload := []byte(`{"email":"usertestcase@test.com","password":"1234567abc"}`)
	req := httptest.NewRequest("POST", "/users/signin", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 200 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 200)
	}
}

func TestGetUser(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/%d", pk)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
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

func TestChangeEmail(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	payload := struct {
		Id    int    `json:"user_id"`
		Email string `json:"email"`
	}{Id: pk, Email: "usertestcase@test.com"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PATCH", "/users/change-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 70 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 70)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if status := doCheckCode(respBody); status != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, status)
	}

}

func TestChangePassword(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	payload := struct {
		Id       int    `json:"user_id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{Id: pk, Email: "usertestcase@test.com", Password: "1234567abc"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PATCH", "/users/change-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 70 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 70)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if status := doCheckCode(respBody); status != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, status)
	}
}

func TestRestorePassword(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	body := []byte(`{"email":"usertestcase@test.com"}`)
	req := httptest.NewRequest("POST", "/users/restorepassword", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 70 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 70)
	}
	bodyRs, _ := io.ReadAll(resp.Body)
	if status := doCheckCode(bodyRs); status != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, status)
	}
}

func TestChangeRole(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	payload := struct {
		Id   int    `json:"user_id"`
		Role string `json:"role"`
	}{pk, "seller"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PATCH", "/users/change-role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp.ContentLength < 20 {
		t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 20)
	}
}

func TestNewCode(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	body := []byte(`{"email":"usertestcase@test.com"}`)
	req := httptest.NewRequest("POST", "/users/restorepassword", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	bodyRs, _ := io.ReadAll(resp.Body)
	req2 := httptest.NewRequest("PATCH", "/users/new-code", bytes.NewBuffer(bodyRs))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := appTest.Test(req2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	if resp2.ContentLength < 70 {
		t.Errorf("Content length recibed %d < expected %d", resp2.ContentLength, 70)
	}

	bodyRs2, _ := io.ReadAll(resp2.Body)
	if status := doCheckCode(bodyRs2); status != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, status)
	}

}

func TestGetView(t *testing.T) {
	initConnection()
	token, _ := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/views?pk=%d&start_date=2023-08-07&end_date=2023-08-08", 1)
	req := httptest.NewRequest("GET", path, nil)
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

func TestInvalidEmail(t *testing.T) {
	initConnection()
	appTest := fiber.New()
	router.Router(appTest)
	payload := []byte(`{"username":"userTestCase2","email":"usertetscase2@testcase.com","password":"1234567abc","tel":"+5353464588"}`)
	req := httptest.NewRequest("POST", "/users/signup", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	req2 := httptest.NewRequest("POST", "/users/invalid-email", bytes.NewBuffer(body))
	req2.Header.Set("Content-type", "application/json")
	resp2, err := appTest.Test(req2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp2.StatusCode)
	}
	if resp2.ContentLength < 30 {
		t.Errorf("Content length recibed %d < expected %d", resp2.ContentLength, 30)
	}

}

func TestGetComments(t *testing.T) {
	initConnection()
	token, _ := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/comments?pub=%d&page=%d&page_size=%d", 1, 1, 1)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	
}

func TestAddTrolley(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	payload := struct {
		Id      int     `json:"user_id"`
		ItemId  int     `json:"item_id"`
		Count   int     `json:"count"`
		State   string  `json:"state"`
		Payment float64 `json:"payment"`
	}{
		Id: pk, ItemId: 37, Count: 1, State: "added", Payment: 120,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/users/trolley", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	for i := 0; i < 2; i++ {
		resp, err := appTest.Test(req, -1)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusAccepted {
			t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
		}
		if resp.ContentLength < 20 {
			t.Errorf("Content length recibed %d < expected %d", resp.ContentLength, 20)
		}
	}

}

func TestGetTrolley(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/trolley/%d", pk)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	if resp.ContentLength < 100 {
		t.Errorf("Content length recibed %d < expected %d",resp.ContentLength,100)
	}
}

func TestUpdateTrolley(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	dataItems := struct {
		Items []models.Trolley `json:"items"`
	}{}
	path := fmt.Sprintf("/users/trolley/%d", pk)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, _ := appTest.Test(req, -1)
	payload, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(payload, &dataItems); err != nil {
		t.Fatal(err)
	}
	item := &dataItems.Items[0]
	item.Count = 3
	item.Payment = item.Payment * float64(item.Count)
	body, _ := json.Marshal(item)
	req2 := httptest.NewRequest("PATCH", "/users/trolley", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", token)
	resp2, err := appTest.Test(req2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp2.StatusCode)
	}
	if resp2.ContentLength < 20 {
		t.Errorf("Content length recibed %d < expected %d",resp2.ContentLength,20)
	}

}

func TestDeleteTrolley(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	dataItems := struct {
		Items []models.Trolley `json:"items"`
	}{}
	path := fmt.Sprintf("/users/trolley/%d", pk)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, _ := appTest.Test(req, -1)
	payload, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(payload, &dataItems); err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(dataItems.Items[0])
	req2 := httptest.NewRequest("DELETE", "/users/trolley", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", token)
	resp2, err := appTest.Test(req2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != fiber.StatusAccepted {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp2.StatusCode)
	}
	
}

func TestDeleteAllTrolley(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/all-trolley/%d", pk)
	req := httptest.NewRequest("DELETE", path, nil)
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

func TestUpdateUser(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	userid, err := w.CreateFormField("user_id")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = userid.Write([]byte(fmt.Sprintf("%d", pk))); err != nil {
		t.Fatal(err)
	}
	address, err := w.CreateFormField("address")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = address.Write([]byte("Test de users update")); err != nil {
		t.Fatal(err)
	}
	w.Close()

	req := httptest.NewRequest("PATCH", "/users/", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", token)
	

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 202 {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
}

func TestCreatePublishings(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	var body bytes.Buffer 
	w := multipart.NewWriter(&body)
	userid, _ := w.CreateFormField("user_id")
	userid.Write([]byte(fmt.Sprintf("%d", pk)))
	category,_ := w.CreateFormField("category")
	category.Write([]byte("Fashion"))
	subca, _ := w.CreateFormField("subcategory")
	subca.Write([]byte("Man"))
	productname,_ := w.CreateFormField("product_name")
	productname.Write([]byte("T-shirt"))
	price, _ := w.CreateFormField("price")
	price.Write([]byte(fmt.Sprintf("%f", 40.0)))
	w.Close()
	req := httptest.NewRequest("POST", "/users/create-publishing", &body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 202 {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
}

func TestGetPublishings(t *testing.T)  {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/get-publishings/%d?page=1&page_size=2", pk)
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := appTest.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusOK, resp.StatusCode)
	}
	if resp.ContentLength < 100 {
		t.Errorf("Content length error expected %d > recived %d", resp.ContentLength, 100)
	}
	
}

func TestUpdatePublishing(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/get-publishings/%d?page=1&page_size=2", pk)
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
	pubId := mapRsp.Publishings[0].ID

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	publishingId,_ := w.CreateFormField("publishing_id")
	publishingId.Write([]byte(fmt.Sprintf("%d", pubId)))
	inOffer, _ := w.CreateFormField("in_offer")
	inOffer.Write([]byte("30"))
	w.Close()

	req2 := httptest.NewRequest("PATCH", "/users/update-publishing", &b)
	req2.Header.Set("Authorization", token)
	req2.Header.Set("Content-Type", w.FormDataContentType())

	resp2, err := appTest.Test(req2, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp2.StatusCode != 202 {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
}

func TestDeletePublishing(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	path := fmt.Sprintf("/users/get-publishings/%d?page=1&page_size=2", pk)
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
	pubId := mapRsp.Publishings[0].ID

	var payload = make(map[string]int)
	payload["publishing_id"] = int(pubId)

	data,_ := json.Marshal(payload)

	req2 := httptest.NewRequest("DELETE", "/users/delete-publishing", bytes.NewBuffer(data))
	req2.Header.Set("Authorization", token)
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := appTest.Test(req2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != 202 {
		t.Errorf("Status code error expected %d and recived %d", fiber.StatusAccepted, resp.StatusCode)
	}
}

func TestDeleteUser(t *testing.T) {
	initConnection()
	token, pk := doLogin("usertestcase@test.com", "1234567abc","users")
	appTest := fiber.New()
	router.Router(appTest)
	payload := struct {
		Id int `json:"user_id"`
	}{pk}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("DELETE", "/users/", bytes.NewBuffer(body))
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
