package handlers

import (
	"avitoInternAssignment/internal/domain"
	"avitoInternAssignment/internal/service"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

const (
	currency = "currency"
)

// Handler processes all http handlers and consists of service realization.
type Handler struct {
	GB  *service.GrossBook
	log *logrus.Logger
}

// NewHandler sets all Handler's values and returns Handler's pointer.
func NewHandler(gb *service.GrossBook, logger *logrus.Logger) *Handler {
	return &Handler{
		GB:  gb,
		log: logger,
	}
}

// InitRoutes initializes routes with necessary middlewares.
func (handler *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/users", func(r chi.Router) {
		r.Post("/balance", handler.balanceHandler)
		r.Post("/history", handler.historyHandler)
	})

	r.Route("/operations", func(r chi.Router) {
		r.Post("/deposit", handler.depositHandler)
		r.Post("/withdraw", handler.withdrawHandler)
		r.Post("/transfer", handler.transferHandler)
	})

	return r
}

// balanceHandler handles getting domain.User's balance.
func (handler *Handler) balanceHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	user := domain.User{}
	if err = json.Unmarshal(data, &user); err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	balance, err := handler.GB.Balance(user.ID)
	if err != nil {
		handler.log.Printf("BALANCE ERROR: <%s>", err)
		processError(w, http.StatusBadRequest, err)
		return
	}
	respBody, err := json.Marshal(balance)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(respBody); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
}

// depositHandler handles increasing domain.User's balance.
func (handler *Handler) depositHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	input := domain.OperationInput{}
	if err = json.Unmarshal(data, &input); err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	operationInfo, err := handler.GB.DepositMoney(
		input.InitiatorID, input.Amount)
	if err != nil {
		handler.log.Printf("DEPOSIT ERROR: <%s>", err)
		processError(w, http.StatusBadRequest, err)
		return
	}
	respBody, err := json.Marshal(operationInfo)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(respBody); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
}

// depositHandler handles decreasing domain.User's balance.
func (handler *Handler) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	var currencyValue string
	// query is necessary for currency type value
	value, ok := r.URL.Query()[currency]
	if ok && len(value) != 0 {
		currencyValue = value[0]
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	input := domain.OperationInput{}
	if err = json.Unmarshal(data, &input); err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	operationInfo, err := handler.GB.WithdrawMoney(
		input.InitiatorID, input.Amount, currencyValue)
	if err != nil {
		handler.log.Printf("WITHDRAW ERROR: <%s>", err)
		processError(w, http.StatusBadRequest, err)
		return
	}
	respBody, err := json.Marshal(operationInfo)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(respBody); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
}

// transferHandler handles transfer from one domain.User to another.
func (handler *Handler) transferHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	input := domain.OperationInput{}
	if err = json.Unmarshal(data, &input); err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	operationInfo, err := handler.GB.TransferMoney(
		input.InitiatorID, input.ReceiverID, input.Amount)
	if err != nil {
		handler.log.Printf("TRANSFER ERROR: <%s>", err)
		processError(w, http.StatusBadRequest, err)
		return
	}
	respBody, err := json.Marshal(operationInfo)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(respBody); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
}

// historyHandler handles getting info about all domain.User's Operations.
func (handler *Handler) historyHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	input := domain.HistoryInput{}
	if err = json.Unmarshal(data, &input); err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}
	operationInfo, err := handler.GB.History(
		input.ID, input.Quantity, input.Mode)
	if err != nil {
		handler.log.Printf("HISTORY ERROR: <%s>", err)
		processError(w, http.StatusBadRequest, err)
		return
	}
	respBody, err := json.Marshal(operationInfo)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(respBody); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
}

// processError sends status code with error text.
func processError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	respBody, err := json.Marshal(domain.ErrorJSON{Message: err.Error()})
	if err != nil {
		return
	}
	_, _ = w.Write(respBody)
}
