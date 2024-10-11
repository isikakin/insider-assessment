package repository

import (
	"assesment/internal/domain/entity"
	"assesment/internal/domain/enum"
	"assesment/pkg/sqlite"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type MessageRepository interface {
	Insert(message *entity.Message)
	RetrieveSentMessages(page int) []entity.Message
	RetrieveSentMessagesTotalCount() int
	RetrieveUnSentMessagesByLimit(count int) []entity.Message
	MarkAsSent(messageId uuid.UUID)
}

type messageRepository struct {
	sqlClient sqlite.Client
}

func (m *messageRepository) Insert(message *entity.Message) {

	var (
		db  *sql.DB
		err error
	)

	if db, err = m.sqlClient.OpenConnection(); err != nil {
		panic(fmt.Sprintf(OpenConnectionError, err.Error()))
	}

	defer m.sqlClient.CloseConnection(db)

	if _, err = db.Exec("insert into messages (message_id, recipient, content, status, sent_date, created_at) values($1, $2, $3, $4, $5, $6)",
		message.MessageId,
		message.Recipient,
		message.Content,
		message.Status,
		message.SentDate,
		message.CreatedAt.Format(time.DateTime)); err != nil {
		panic(fmt.Sprintf(InsertErr, err.Error()))
	}
}

func (m *messageRepository) RetrieveSentMessages(page int) (messages []entity.Message) {
	var (
		db  *sql.DB
		err error
	)

	if db, err = m.sqlClient.OpenConnection(); err != nil {
		panic(fmt.Sprintf(OpenConnectionError, err.Error()))
	}
	defer m.sqlClient.CloseConnection(db)

	limit := 10
	offset := limit * (page - 1)

	query := fmt.Sprintf("select message_id, recipient, content, status, sent_date, created_at  from messages where status=%d LIMIT %d OFFSET %d", enum.Sent, limit, offset)
	rows, err := db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		panic(fmt.Sprintf(RetrieveDataError, err.Error()))
	}

	defer rows.Close()

	for rows.Next() {
		item := entity.Message{}

		createdAtStr, sentDateStr := "", ""

		if err = rows.Scan(&item.MessageId, &item.Recipient, &item.Content, &item.Status, &sentDateStr, &createdAtStr); err != nil {
			panic(fmt.Sprintf(RetrieveDataError, err.Error()))
		}

		item.CreatedAt, _ = time.Parse(time.DateTime, createdAtStr)
		sentDate, _ := time.Parse(time.DateTime, sentDateStr)

		item.SentDate = &sentDate

		messages = append(messages, item)
	}

	return messages
}

func (m *messageRepository) RetrieveSentMessagesTotalCount() (count int) {

	var (
		db  *sql.DB
		err error
	)

	if db, err = m.sqlClient.OpenConnection(); err != nil {
		panic(fmt.Sprintf(OpenConnectionError, err.Error()))
	}
	defer m.sqlClient.CloseConnection(db)

	err = db.QueryRow("select count(message_id) from messages where status = $1", enum.Sent).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0
		}
		panic(fmt.Sprintf(RetrieveDataError, err.Error()))
	}

	return count
}

func (m *messageRepository) RetrieveUnSentMessagesByLimit(limit int) (messages []entity.Message) {
	var (
		db  *sql.DB
		err error
	)

	if db, err = m.sqlClient.OpenConnection(); err != nil {
		panic(fmt.Sprintf(OpenConnectionError, err.Error()))
	}
	defer m.sqlClient.CloseConnection(db)

	query := fmt.Sprintf("select message_id, recipient, content from messages where status = %d order by created_at LIMIT %d", enum.Pending, limit)
	rows, err := db.Query(query)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		panic(fmt.Sprintf(RetrieveDataError, err.Error()))
	}

	defer rows.Close()

	for rows.Next() {
		item := entity.Message{}
		if err = rows.Scan(&item.MessageId, &item.Recipient, &item.Content); err != nil {
			panic(fmt.Sprintf(RetrieveDataError, err.Error()))
		}
		messages = append(messages, item)
	}

	return messages
}

func (m *messageRepository) MarkAsSent(messageId uuid.UUID) {
	var (
		db  *sql.DB
		err error
	)

	if db, err = m.sqlClient.OpenConnection(); err != nil {
		panic(OpenConnectionError)
	}

	defer m.sqlClient.CloseConnection(db)

	if _, err = db.Exec("update messages set status=$1, sent_date=$2  where message_id = $3",
		enum.Sent,
		time.Now().Format(time.DateTime),
		messageId); err != nil {
		panic(MarkAsSentError)
	}
}

func NewMessageRepository(sqlClient sqlite.Client) MessageRepository {
	return &messageRepository{sqlClient: sqlClient}
}
