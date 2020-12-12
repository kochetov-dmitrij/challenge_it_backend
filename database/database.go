package database

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"os"
)

func InitDB() *pg.DB {
	connectOptions, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	db := pg.Connect(connectOptions)
	defer db.Close()

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	err = createSchema(db)
	if err != nil {
		panic(err)
	}

	return db
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Challenge)(nil),
		(*UserChallenge)(nil),
		(*Tag)(nil),
		(*ChallengeTag)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
