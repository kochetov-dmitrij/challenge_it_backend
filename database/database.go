package database

import (
	"fmt"
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

	err = createSchema(db)
	if err != nil {
		panic(err)
	}
	examples(db)
	return db
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Story)(nil),
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

// Code snippets
func examples(db *pg.DB) {
	adminUser := &User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	_, err := db.Model(adminUser).Insert()
	if err != nil {
		panic(err)
	}

	coolStory := &Story{
		Title:  "Cool story",
		Author: adminUser,
	}
	_, err = db.Model(coolStory).Insert()
	if err != nil {
		panic(err)
	}

	// Select user by primary key.
	user := &User{Id: adminUser.Id}
	err = db.Model(user).WherePK().Select()
	if err != nil {
		panic(err)
	}

	// Select all users.
	var users []User
	err = db.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	// Select story and associated author in one query.
	story := new(Story)
	err = db.Model(story).
		Relation("Author").
		Where("story.id = ?", coolStory.Id).
		Select()
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	// Output: User<1 admin [admin1@admin admin2@admin]>
	fmt.Println(users)
	// [User<1 admin [admin1@admin admin2@admin]> User<2 root [root1@root root2@root]>]
	fmt.Println(story)
	// Story<1 Cool story User<1 admin [admin1@admin admin2@admin]>>
}
