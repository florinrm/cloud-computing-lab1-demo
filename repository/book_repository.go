package repository

import (
	"context"
	"database/sql"
	"exercise2/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	insertStmt = `INSERT INTO books (id, title, author)
					VALUES ($1, $2, $3) 
					RETURNING id, title, author`

	getStmt = `SELECT id, title, author FROM books`
)

type BookRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewBookRepository(db *sql.DB, logger *logrus.Logger) *BookRepository {
	return &BookRepository{
		db:     db,
		logger: logger,
	}
}

func (b *BookRepository) AddBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	bookId := uuid.New()
	book.ID = bookId.String()
	err := b.db.QueryRowContext(ctx, insertStmt, bookId, book.Title, book.Author)
	if err != nil && err.Err() != nil {
		b.logger.WithError(err.Err()).Errorf("failed to insert book in database")
		return nil, err.Err()
	}
	return book, nil
}

func (b *BookRepository) GetBooks(ctx context.Context) ([]domain.Book, error) {
	rows, err := b.db.QueryContext(ctx, getStmt)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			b.logger.WithError(err).Errorf("failed to close the rows")
		}
	}()

	var books []domain.Book
	for rows.Next() {
		book := domain.Book{}

		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}
