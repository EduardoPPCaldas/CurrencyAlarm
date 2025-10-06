package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/EduardoPPCaldas/CurrencyAlarm/pkg/email"
)

type CurrencyResponse struct {
	EURBRL struct {
		Bid string `json:"bid"`
	} `json:"EURBRL"`
}

func CheckEuro() {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/EUR-BRL")
	if err != nil {
		fmt.Println("Error fetching:", err)
		return
	}
	defer resp.Body.Close()

	var result CurrencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	bid, err := strconv.ParseFloat(result.EURBRL.Bid, 64)
	if err != nil {
		fmt.Println("Error parsing bid:", err)
		return
	}

	if bid > 6.0 {
		fmt.Printf("Euro above 6.0")
		return
	}

	fmt.Printf("🚨 Euro dropped! Current EUR/BRL: %.4f\n", bid)

	subject := "Euro Alert: EUR < 6.0"
	body := fmt.Sprintf("The Euro just dropped below 6 BRL!\n\nCurrent rate: %.4f", bid)

	if err := email.SendEmail(subject, body); err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("✅ Alert email sent!")
	}
}
