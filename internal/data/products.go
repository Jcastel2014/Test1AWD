package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jcastel2014/test1/internal/validator"
)

type Product struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	Image_url      string    `json:"image_url"`
	Average_rating float64   `json:"average_rating"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}

type ProductModel struct {
	DB *sql.DB
}

func (p ProductModel) Insert(product *Product) error {

	query := `
	INSERT INTO images (image_url)
	VALUES ($1)
	RETURNING id
	`

	args := []any{product.Image_url}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&id)

	if err != nil {
		return err
	}

	query = `
	INSERT INTO products (name, description, category, image_id, average_rating) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, created_at, updated_at`

	//0 is default value for average_rating
	args = []any{product.Name, product.Description, product.Category, id, 0}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.Created_at, &product.Updated_at)
}

func (p ProductModel) Get(id int64) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT P.id, P.name, P.description, P.category, I.image_url, P.average_rating, P.created_at, P.updated_at
	FROM products AS P
	INNER JOIN images AS I ON P.image_id = I.id
	WHERE P.id = $1;
	`

	var product Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Image_url, &product.Average_rating, &product.Created_at, &product.Updated_at)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &product, nil
}

// func (c CommentModel) GetAll(content string, author string, filters Filters) ([]*Comment, error) {
// 	query := `
// 	SELECT id, created_at, content, author, version
// 	FROM comments
// 	WHERE (to_tsvector('simple', content) @@
// 		plainto_tsquery('simple', $1) OR $1 = '')
//     AND (to_tsvector('simple', author) @@
// 		plainto_tsquery('simple', $2) OR $2 = '')
// 	ORDER BY id
// 	LIMIT $3 OFFSET $4
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	rows, err := c.DB.QueryContext(ctx, query, content, author, filters.limit(), filters.offset())

// 	if err != nil {
// 		return nil, err
// 	}

// 	comments := []*Comment{}

// 	for rows.Next() {
// 		var comment Comment
// 		err := rows.Scan(&comment.ID, &comment.CreatedAt, &comment.Content, &comment.Author, &comment.Version)

// 		if err != nil {
// 			return nil, err
// 		}

// 		comments = append(comments, &comment)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return comments, nil
// }

// func (c CommentModel) Update(comment *Comment) error {
// 	query := `
// 	UPDATE comments
// 	SET content = $1, author = $2, version = version + 1
// 	WHERE id = $3
// 	RETURNING version
// 	`

// 	args := []any{comment.Content, comment.Author, comment.ID}
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
// }

// func (c CommentModel) Delete(id int64) error {
// 	if id < 1 {
// 		return ErrRecordNotFound
// 	}

// 	query := `
// 	DELETE FROM comment
// 	WHERE id =$1
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	result, err := c.DB.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return ErrRecordNotFound
// 	}

// 	return nil
// }

func ValidateProduct(v *validator.Validator, product *Product, handlerId int) {

	switch handlerId {
	case 1:
		v.Check(product.Name != "", "name", "must be provided")
		v.Check(product.Description != "", "description", "must be provided")
		v.Check(product.Category != "", "category", "must be provided")
		v.Check(product.Image_url != "", "image_url", "must be provided")

		v.Check(len(product.Name) <= 100, "name", "must not be more than 100 byte long")
		v.Check(len(product.Description) <= 100, "description", "must not be more than 100 byte long")
		v.Check(len(product.Category) <= 100, "category", "must not be more than 100 byte long")
	default:
		log.Println("Unable to locate handler ID: %s", handlerId)
		v.AddError("default", "Handler ID not provided")
	}
}
