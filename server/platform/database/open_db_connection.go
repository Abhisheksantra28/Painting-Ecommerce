package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/niladri2003/PaintingEcommerce/app/queries"
	"os"
)

type Queries struct {
	*queries.UserQueries
	*queries.CategoryQueries
	*queries.ProductQueries
	*queries.ProductImageQueries
	*queries.ContactQueries // Add ContactQueries
	*queries.AddressQueries // Add AddressQueries
	*queries.CartQueries    // Add CartQueries
	*queries.OrderQueries
	*queries.OrderItemsQuery
	*queries.CouponQueries
	*queries.ProductSubCategoryQuery
	*queries.ProductSizeQuery
}

// OpenDbConnection  func for opening database connection
func OpenDbConnection() (*Queries, error) {
	//Define Database connection variables.
	var (
		db  *sqlx.DB
		err error
	)
	//Get DB_TYPE value from .env file.
	dbType := os.Getenv("DB_TYPE")

	//Define a new Database connection with right DB type.
	switch dbType {
	case "pgx":
		db, err = PostgresSQLConnection()
	}
	if err != nil {
		return nil, err
	}
	return &Queries{
		UserQueries:             &queries.UserQueries{DB: db},
		CategoryQueries:         &queries.CategoryQueries{DB: db},
		ProductQueries:          &queries.ProductQueries{DB: db},
		ProductImageQueries:     &queries.ProductImageQueries{DB: db},
		ContactQueries:          &queries.ContactQueries{DB: db}, // Add ContactQueries
		AddressQueries:          &queries.AddressQueries{DB: db}, // Add AddressQueries
		CartQueries:             &queries.CartQueries{DB: db},
		OrderQueries:            &queries.OrderQueries{DB: db},
		OrderItemsQuery:         &queries.OrderItemsQuery{DB: db},
		CouponQueries:           &queries.CouponQueries{DB: db},
		ProductSubCategoryQuery: &queries.ProductSubCategoryQuery{DB: db},
		ProductSizeQuery:        &queries.ProductSizeQuery{DB: db},
	}, nil
}
