package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/agandreev/avito-intern-assignment/internal/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertTransferOperationSQL = "INSERT INTO operations(initiator_id, type, amount, " +
		"time, receiver_id) " +
		"VALUES(" +
		"(SELECT id from users WHERE user_id=$1), " +
		"$2, $3, $4, " +
		"(SELECT id from users WHERE user_id=$5))"
	insertNonTransferOperationSQL = "INSERT INTO operations(initiator_id, type, amount, " +
		"time, receiver_id) " +
		"VALUES(" +
		"(SELECT id from users WHERE user_id=$1), " +
		"$2, $3, $4, " +
		"NULL)"
)

var (
	ErrNotConnected = errors.New("there is no db connection")
	ErrNoSuchUser   = errors.New("user with this id doesn't exist")
	ErrNoOperations = errors.New("this user hasn't any operations")

	InitialAmountValue = 0
)

// GrossBookStorage is implementation of service.GrossBookRepository
type GrossBookStorage struct {
	pool   *pgxpool.Pool
	Config ConnectionConfig
}

// NewGrossBookStorage create an entity and returns pointer.
func NewGrossBookStorage(config ConnectionConfig) *GrossBookStorage {
	return &GrossBookStorage{
		Config: config,
	}
}

// ConnectionConfig contains all necessary parameters for db connection.
type ConnectionConfig struct {
	Username string
	Password string
	NameDB   string
	Port     string
}

// Connect creates connection and returns error otherwise.
func (storage *GrossBookStorage) Connect() error {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s"+
		"?sslmode=disable", storage.Config.Username, storage.Config.Password,
		storage.Config.Port, storage.Config.NameDB)
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("can't connect db <%w>", err)
	}
	storage.pool = pool
	return nil
}

// User return domain.User by id.
func (storage *GrossBookStorage) User(id int64) (*domain.User, error) {
	row := storage.pool.QueryRow(context.Background(),
		"SELECT * FROM users WHERE user_id=$1", id)
	var user domain.User
	var dbNumber int
	if err := row.Scan(&dbNumber, &user.ID, &user.Amount); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoSuchUser
		}
		return nil, fmt.Errorf("can't read from db <%w>", err)
	}
	return &user, nil
}

// AddUser initialize domain.User by id with initial amount value.
func (storage *GrossBookStorage) AddUser(id int64) error {
	if _, err := storage.pool.Exec(context.Background(),
		"INSERT INTO users(user_id, amount) VALUES($1, $2)",
		id, InitialAmountValue); err != nil {
		return fmt.Errorf("can't add to db <%w>", err)
	}
	return nil
}

// AddOperation adds domain.Operation to the storage and updates domain.User from it.
func (storage *GrossBookStorage) AddOperation(ctx context.Context, operation domain.Operation) error {
	if storage.pool == nil {
		return ErrNotConnected
	}
	// start transaction to add operations and update users
	tx, err := storage.pool.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	// check that operation is correct
	if err = operation.Validate(); err != nil {
		return fmt.Errorf("can't add operation: <%w>", err)
	}
	// check if users are existed
	if _, err = storage.User(operation.Initiator.ID); err != nil {
		return fmt.Errorf("error while adding operation "+
			"(can't get initiator): <%w>", err)
	}
	if operation.IsTransfer() {
		if _, err = storage.User(operation.Receiver.ID); err != nil {
			return fmt.Errorf("error while adding operation "+
				"(can't get receiver): <%w>", err)
		}
	}
	// try to execute queries
	if err = processOperation(ctx, tx, operation); err != nil {
		return fmt.Errorf("can't execute transaction: <%w>", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("can't commit operation transaction: <%w>", err)
	}
	return nil
}

// Operations returns domain.Operation's slice by domain.User's id,
// sorted as domain.SortingMode and limited as offset
func (storage *GrossBookStorage) Operations(id int64, offset int64,
	mode domain.SortingMode) ([]domain.RepositoryOperation, error) {
	if offset <= 0 {
		return nil, fmt.Errorf("incorrect offset value")
	}
	rows, err := storage.pool.Query(context.Background(),
		"SELECT * FROM operations WHERE initiator_id="+
			"(SELECT id FROM users WHERE user_id=$1) ORDER BY time DESC LIMIT $2", id, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get operations: <%w>", err)
	}
	var operationQuantity int64
	operations := make([]domain.RepositoryOperation, 0)
	var dbNumber int
	var optionalID sql.NullInt64
	for rows.Next() && operationQuantity < offset {
		var operation domain.RepositoryOperation
		if err := rows.Scan(&dbNumber, &operation.InitiatorID, &operation.Type,
			&operation.Amount, &operation.Timestamp, &optionalID); err != nil {
			if err == pgx.ErrNoRows {
				return nil, ErrNoOperations
			}
			return nil, fmt.Errorf("can't read from db <%w>", err)
		}
		if optionalID.Valid {
			operation.ReceiverID = optionalID.Int64
		}
		operationQuantity++
		operations = append(operations, operation)
	}
	sortOperations(operations, mode)
	return operations, nil
}

// Shutdown closes connection. It blocks while all current queries are processing.
func (storage GrossBookStorage) Shutdown() {
	if storage.pool != nil {
		storage.pool.Close()
	}
}

// sortOperations sorts slice of domain.RepositoryOperation by domain.SortingMode.
func sortOperations(operations []domain.RepositoryOperation, mode domain.SortingMode) {
	sort.SliceStable(operations, func(i, j int) bool {
		switch mode {
		case domain.AmountMode:
			return operations[i].Amount > operations[j].Amount
		case domain.DateMode:
			return operations[i].Timestamp.After(operations[j].Timestamp)
		default:
			return true
		}
	})
}

// processOperation executes pgx.Tx by domain.Operation.
func processOperation(ctx context.Context, tx pgx.Tx, operation domain.Operation) error {
	// update initiator
	if err := updateUser(ctx, tx, *operation.Initiator); err != nil {
		return fmt.Errorf("transaction initiator error: <%w>", err)
	}
	// update receiver if it's existed
	if operation.IsTransfer() {
		if err := updateUser(ctx, tx, *operation.Receiver); err != nil {
			return fmt.Errorf("transaction receiver error: <%w>", err)
		}
	}
	// add operation info to db
	if err := addOperation(ctx, tx, operation); err != nil {
		return fmt.Errorf("can't add operation to db: <%w>", err)
	}
	return nil
}

// updateUser updates domain.User's balance in db.
func updateUser(ctx context.Context, tx pgx.Tx, user domain.User) error {
	if _, err := tx.Exec(ctx, "UPDATE users SET amount=$1 WHERE user_id=$2",
		user.Amount, user.ID); err != nil {
		return fmt.Errorf("can't execute user updation: <%w>", err)
	}
	return nil
}

// addOperation adds domain.Operation's info to db.
func addOperation(ctx context.Context, tx pgx.Tx, operation domain.Operation) error {
	// add non-duplex transaction
	if !operation.IsTransfer() {
		if _, err := tx.Exec(ctx,
			insertNonTransferOperationSQL,
			operation.Initiator.ID, operation.Type, operation.Amount,
			operation.Timestamp); err != nil {
			return fmt.Errorf("can't add operation to db <%w>", err)
		}
	} else {
		// add duplex transaction (transfer-out)
		if err := addDuplexOperation(ctx, tx, operation); err != nil {
			return fmt.Errorf("origin operation error: <%w>", err)
		}
		// add reversed duplex transaction (transfer-in)
		reversed, err := operation.Reverse()
		if err != nil {
			return fmt.Errorf("can't add reversed transaction: <%w>", err)
		}
		if err = addDuplexOperation(ctx, tx, *reversed); err != nil {
			return fmt.Errorf("origin operation error: <%w>", err)
		}
	}
	return nil
}

// addDuplexOperation adds domain.Operation's info with receiver id to db.
func addDuplexOperation(ctx context.Context, tx pgx.Tx, operation domain.Operation) error {
	if _, err := tx.Exec(ctx,
		insertTransferOperationSQL,
		operation.Initiator.ID, operation.Type, operation.Amount,
		operation.Timestamp, operation.Receiver.ID); err != nil {
		return fmt.Errorf("can't add operation to db <%w>", err)
	}
	return nil
}
