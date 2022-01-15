package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	timeout = 10 * time.Second

	apiURL  = "http://api.exchangeratesapi.io/v1/"
	latest  = "latest"
	symbols = "symbols"

	accessKeyTag = "access_key"
	baseTag      = "base"
	symbolsTag   = "symbols"

	rub = "RUB"
	eur = "EUR"
)

// ExchangeAPI implements Converter applying exchangerateapi v1.
type ExchangeAPI struct {
	client *http.Client
	apiKey string
}

// NewExchangeAPI sets timeout and returns pointer.
func NewExchangeAPI(apiKey string) *ExchangeAPI {
	return &ExchangeAPI{
		client: &http.Client{Timeout: timeout},
		apiKey: apiKey,
	}
}

// SupportedCurrencies represents slice of available currencies.
type SupportedCurrencies []string

// Contains implements fast search and returns true or false.
func (currencies SupportedCurrencies) Contains(currency string) bool {
	i := sort.SearchStrings(currencies, currency)
	return i < len(currencies) && currencies[i] == currency
}

// ContainsAll implements Contains for several currencies.
func (currencies SupportedCurrencies) ContainsAll(checkingCurrencies ...string) error {
	for _, currency := range checkingCurrencies {
		if !currencies.Contains(currency) {
			return fmt.Errorf("%s is unsupported", currency)
		}
	}
	return nil
}

type ConversionResponse struct {
	Success bool            `json:"success"`
	Error   ConversionError `json:"error"`
	ConversionResponseInfo
}

type ConversionResponseInfo struct {
	Timestamp int                `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

type ConversionError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

func (conversionError ConversionError) Error() string {
	return fmt.Sprintf("Code %d: %s", conversionError.Code, conversionError.Info)
}

type SymbolsResponse struct {
	Success bool              `json:"success"`
	Error   ConversionError   `json:"error"`
	Symbols map[string]string `json:"symbols"`
}

// Currencies returns SupportedCurrencies for each request.
func (symbolsResponse SymbolsResponse) Currencies() (SupportedCurrencies, error) {
	if !symbolsResponse.Success {
		return nil, symbolsResponse.Error
	}
	if symbolsResponse.Symbols == nil {
		return nil, fmt.Errorf("there is no any supported currency")
	}
	currencies := make(SupportedCurrencies, 0)
	for currency := range symbolsResponse.Symbols {
		currencies = append(currencies, currency)
	}
	sort.Strings(currencies)
	return currencies, nil
}

type BadRequestResponse struct {
	BadRequestError BadRequestError `json:"error"`
}

type BadRequestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (badRequestError BadRequestError) Error() string {
	return fmt.Sprintf("Code <%s>: %s", badRequestError.Code,
		badRequestError.Message)
}

// SupportedSymbols return array of supported currencies.
func (exchange ExchangeAPI) SupportedSymbols() (
	SupportedCurrencies, error) {
	req, err := http.NewRequest("GET", apiURL+symbols, nil)
	if err != nil {
		return nil, fmt.Errorf("exchange convert error: <%w>", err)
	}
	q := req.URL.Query()
	q.Add(accessKeyTag, exchange.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := exchange.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchange request error: <%w>", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("exchange read body error: <%w>", err)
	}
	if resp.StatusCode == http.StatusOK {
		var symbolsResponse SymbolsResponse
		if err = json.Unmarshal(data, &symbolsResponse); err != nil {
			return nil, fmt.Errorf("exchange unmarshal conversion error: <%w>", err)
		}
		return symbolsResponse.Currencies()
	}
	if resp.StatusCode == http.StatusBadRequest {
		var badRequestResponse BadRequestResponse
		if err = json.Unmarshal(data, &badRequestResponse); err != nil {
			return nil, fmt.Errorf("exchange unmarshal bad request error: <%w>", err)
		}
		return nil, badRequestResponse.BadRequestError
	}
	return nil, fmt.Errorf("unexpected status code received from exchager: %d",
		resp.StatusCode)
}

// Convert converts any currency to rubbles.
func (exchange ExchangeAPI) Convert(from string, amount float64) (
	float64, error) {
	// for each request to avoid mutexes
	supportedCurrencies, err := exchange.SupportedSymbols()
	if err != nil {
		return 0, fmt.Errorf("can't get supported symbols: <%w>", err)
	}
	if err = supportedCurrencies.ContainsAll(rub, eur, from); err != nil {
		return 0, fmt.Errorf("can't convert: <%w>", err)
	}
	// request creation
	req, err := http.NewRequest("GET", apiURL+latest, nil)
	if err != nil {
		return 0, fmt.Errorf("exchange convert error: <%w>", err)
	}
	symbols := []string{rub, from}
	q := req.URL.Query()
	q.Add(accessKeyTag, exchange.apiKey)
	q.Add(baseTag, eur)
	q.Add(symbolsTag, strings.Join(symbols, ","))
	req.URL.RawQuery = q.Encode()

	// sending request
	resp, err := exchange.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("exchange request error: <%w>", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("exchange read body error: <%w>", err)
	}
	// process status codes
	switch resp.StatusCode {
	case http.StatusOK:
		var conversionResponse ConversionResponse
		if err = json.Unmarshal(data, &conversionResponse); err != nil {
			return 0, fmt.Errorf("exchange unmarshal conversion error: <%w>", err)
		}
		convertedAmount, err := conversionResponse.Amount(amount, from)
		if err != nil {
			return 0, fmt.Errorf("exchange calculation error: <%w>", err)
		}
		return convertedAmount, nil
	case http.StatusBadRequest:
		var badRequestResponse BadRequestResponse
		if err = json.Unmarshal(data, &badRequestResponse); err != nil {
			return 0, fmt.Errorf("exchange unmarshal bad request error: <%w>", err)
		}
		return 0, badRequestResponse.BadRequestError
	default:
		return 0, fmt.Errorf("unexpected status code received from exchager: %d",
			resp.StatusCode)
	}
}

// Validate is necessary to prevent an error from the API side.
func (conversion ConversionResponse) Validate(currency string) error {
	// check if body exists
	if !conversion.Success {
		return fmt.Errorf("conversion response body is empty")
	}
	// check base currency
	if conversion.Base != eur {
		return fmt.Errorf("returned unsupported currency from base: %s",
			conversion.Base)
	}
	// check rates map
	if err := conversion.checkRatesLen(currency); err != nil {
		return fmt.Errorf("something wrong with rates map in response: <%w>", err)
	}
	// check rate's values
	for curr, rate := range conversion.Rates {
		if rate <= 0 {
			return fmt.Errorf("returned negative rate: %f, by this currency: %s",
				rate, curr)
		}
	}
	return nil
}

// Amount counts final value by formula.
func (conversion ConversionResponse) Amount(amount float64, currency string) (float64, error) {
	if !conversion.Success {
		return 0, conversion.Error
	}
	if err := conversion.Validate(currency); err != nil {
		return 0, fmt.Errorf("validation is failed: <%w>", err)
	}
	return amount / conversion.Rates[currency] * conversion.Rates[rub], nil
}

// checkRatesLen is a part of validation.
func (conversion ConversionResponse) checkRatesLen(currency string) error {
	switch len(conversion.Rates) {
	case 0:
		return fmt.Errorf("returned empty rates map")
	case 1:
		// case, when users input == base or input currency == output currency
		if currency != conversion.Base && currency != rub {
			return fmt.Errorf("not enough rates returned")
		}
		if _, ok := conversion.Rates[rub]; !ok {
			return fmt.Errorf("there is no RUB is response")
		}
	case 2:
		// case with 3 currencies, including RUB and EURO
		if _, ok := conversion.Rates[currency]; !ok {
			return fmt.Errorf("lack of currency in map: %s", currency)
		}
		if _, ok := conversion.Rates[rub]; !ok {
			return fmt.Errorf("there is no RUB is response")
		}
	default:
		return fmt.Errorf("too much currencies")
	}
	return nil
}
