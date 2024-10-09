package main

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/invoice"
	"github.com/stripe/stripe-go/v75/invoiceitem"
)

func main() {

	// Stripe API kalitini o'rnating
	stripe.Key = "key"

	// Foydalanuvchining Stripe customer_id si (mavjud bo'lishi kerak)
	customerID := "custumer_id"
	pm := "payment_method_id"
	// 1. Invoice item na summa 10$
	positiveItemParams := &stripe.InvoiceItemParams{
		Customer:    stripe.String(customerID),
		Amount:      stripe.Int64(1000), // 10.00 USD
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Total Sum"),
	}
	_, err := invoiceitem.New(positiveItemParams)
	if err != nil {
		log.Fatalf("Ошибка при создании invoice item: %v", err)
	}

	// 2. Invoice item na summa -5$ (kredit)
	negativeItemParams := &stripe.InvoiceItemParams{
		Customer:    stripe.String(customerID),
		Amount:      stripe.Int64(-500), // -5.00 USD (kredit)
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Inactivated truck"),
	}
	_, err = invoiceitem.New(negativeItemParams)
	if err != nil {
		log.Fatalf("Ошибка при создании invoice item: %v", err)
	}

	// 3. Invoice yaratish va pending itemlarni qo'shish
	invoiceParams := &stripe.InvoiceParams{
		Customer:                    stripe.String(customerID),
		AutoAdvance:                 stripe.Bool(false),
		PendingInvoiceItemsBehavior: stripe.String("include"), // Pending itemlarni invoiceda aks ettirish
	}
	inv, err := invoice.New(invoiceParams)
	if err != nil {
		log.Fatalf("Ошибка при создании инвойса: %v", err)
	}

	// 4. Finalizatsiya qilish (Invoice summasini hisoblash va tayyorlash)
	invoiceLast, err := invoice.FinalizeInvoice(inv.ID, nil)
	if err != nil {
		log.Fatalf("Ошибка при финализации инвойса: %v", err)
	}

	_, err = invoice.Pay(inv.ID, &stripe.InvoicePayParams{
		PaymentMethod: &pm,
	})
	if err != nil {
		log.Fatalf("Error paying invoice: %v", err)
	}

	fmt.Printf("Инвойс создан и оплачен: %s\n", invoiceLast.HostedInvoiceURL)
}
