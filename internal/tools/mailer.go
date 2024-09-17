package tools

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/mattevans/postmark-go"
)

type ReceiptDetail struct {
	Description string `json:"description"`
	Amount      string `json:"amount"`
}

type Receipt struct {
	ProductURL              string          `json:"product_url"`
	ProductName             string          `json:"product_name"`
	Customer                string          `json:"customer"`
	PaymentStatus           string          `json:"payment_status"`
	CreditCardStatementName string          `json:"credit_card_statement_name"`
	OrderID                 string          `json:"order_id"`
	Date                    string          `json:"date"`
	PickupDate              string          `json:"pickup_date"`
	ReceiptDetails          []ReceiptDetail `json:"receipt_details"`
	Total                   string          `json:"total"`
	SupportURL              string          `json:"support_url"`
	CompanyName             string          `json:"company_name"`
	CompanyAddress          string          `json:"company_address"`
}

func SendReceipt(customerEmail string, receipt Receipt, attachment string) error {
	client := postmark.NewClient(
		postmark.WithClient(&http.Client{
			Transport: &postmark.AuthTransport{Token: os.Getenv("POSTMARK_API_TOKEN")},
		}),
	)

	log.Infof("Receipt for %s: %v", customerEmail, receipt)

	jsonReceipt, err := json.Marshal(receipt)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonReceipt, &data)
	if err != nil {
		return err
	}

	log.Infof("Data for receipt: %v", data)

	file, err := os.Open(attachment)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(attachment)

	attachmentContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	contentType := "application/pdf"

	emailReq := &postmark.Email{
		From:          os.Getenv("POSTMARK_SENDER"),
		To:            customerEmail,
		TemplateID:    36243053,
		TemplateModel: data,
		Tag:           "onboarding",
		TrackOpens:    true,
		Attachments: []postmark.EmailAttachment{
			{
				Name:        "receipt.pdf",
				ContentType: &contentType,
				Content:     attachmentContent,
			},
		},
	}

	_, _, err = client.Email.Send(emailReq)
	if err != nil {
		return err
	}

	return nil
}
