package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/bmtrann/sesc-component/internal/exception"
)

type FinanceService interface {
	CreateFinanceAccount()
	CreateInvoice()
	GetGraduationStatus()
}

type LibraryService interface {
	CreateLibraryAccount()
}

type InvoicePayload struct {
	Amount  float32 `json:"amount"`
	DueDate string  `json:"dueDate"`
	Type    string  `json:"type"`
	Account Account `json:"account"`
}

type Account struct {
	StudentID string `json:"studentId"`
}

type ServiceHandler struct{}

// Constants
var FINANCE_URL string = os.Getenv("FINANCE_URL")
var LIBRARY_URL string = os.Getenv("LIBRARY_URL")

// Singleton pattern
var (
	handler *ServiceHandler
	once    sync.Once
)

func GetInstance() *ServiceHandler {
	once.Do(func() {
		handler = &ServiceHandler{}
	})
	return handler
}

func (handler *ServiceHandler) CreateFinanceAccount(studentId string) error {
	payload, _ := json.Marshal(map[string]string{
		"studentId": studentId,
	})

	return postRequest(payload, FINANCE_URL+"/accounts")
}

func (handler *ServiceHandler) CreateInvoice(studentId string, fees float32) error {
	payload, _ := json.Marshal(InvoicePayload{
		Amount:  fees,
		DueDate: time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		Type:    "TUITION_FEES",
		Account: Account{
			StudentID: studentId,
		},
	})

	return postRequest(payload, FINANCE_URL+"/invoices")
}

func (handler *ServiceHandler) GetGraduationStatus(studentId string) (bool, error) {
	url := FINANCE_URL + "/accounts/student/" + studentId
	response, err := getRequest(url)

	if err != nil {
		return false, err
	}

	hasOutstandingBalance, ok := response["hasOutstandingBalance"].(bool)
	if !ok {
		return false, exception.ServiceException(url)
	}

	return hasOutstandingBalance, nil
}

func (handler *ServiceHandler) CreateLibraryAccount(studentId string) error {
	payload, _ := json.Marshal(map[string]string{
		"studentId": studentId,
	})

	return postRequest(payload, LIBRARY_URL+"/api/register")
}

func postRequest(payload []byte, url string) error {
	body := bytes.NewBuffer(payload)

	resp, err := http.Post(url, "application/json", body)

	if err != nil {
		log.Println(err)
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Println(resp.StatusCode)
		return exception.ServiceException(url)
	}
	defer resp.Body.Close()

	return nil
}

func getRequest(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	respBody, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(respBody, &data); err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}
