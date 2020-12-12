package database

import "github.com/go-pg/pg/v10"

func examples(db *pg.DB) {
	// https://pg.uptrace.dev/queries/

	// insert user
	adminUser := &User{
		Name:        "admin",
		Email:       "admin1@admin",
		EncPassword: "qwerty",
	}
	_, err := db.Model(adminUser).Insert()
	if err != nil {
		panic(err)
	}

	// insert challenge
	challenge := &Challenge{
		Title:    "Cool story",
		AuthorId: adminUser.Id,
		Rating:   32,
		Taken:    Assigned,
	}
	_, err = db.Model(challenge).Insert()
	if err != nil {
		panic(err)
	}

	// get user by primary key
	user := &User{Id: adminUser.Id}
	err = db.Model(user).WherePK().Select()
	if err != nil {
		panic(err)
	}

	// get all users
	var users []User
	err = db.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	// get challenge and its author
	challenge1 := new(Challenge)
	err = db.Model(challenge1).
		Relation("Author").
		Where("challenge.id = ?", challenge.Id).
		Select()
	if err != nil {
		panic(err)
	}
}
