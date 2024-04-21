package service

import (
	"bytes"
	"encoding/json"
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

	return postRequest(payload, "/accounts")
}

func (handler *ServiceHandler) CreateInvoice(studentId string, fees float32) error {
	payload, _ := json.Marshal(InvoicePayload{
		Amount:  fees,
		DueDate: time.Now().AddDate(0, 1, 0).String(),
		Type:    "TUITION_FEES",
		Account: Account{
			StudentID: studentId,
		},
	})

	return postRequest(payload, "/invoices")
}

func postRequest(payload []byte, endpoint string) error {
	body := bytes.NewBuffer(payload)
	financeURL := os.Getenv("FINANCE_URL")

	resp, err := http.Post(financeURL+endpoint, "application/json", body)

	if err != nil {
		log.Println(err)
		return err
	}

	if resp.StatusCode != 201 {
		log.Println(resp.StatusCode)
		return exception.ServiceException(financeURL + endpoint)
	}
	defer resp.Body.Close()

	return nil
}
