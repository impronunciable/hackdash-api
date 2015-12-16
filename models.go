package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

var db gorm.DB

func InitDB(dbConfig string) {
    var err error
    db, err = gorm.Open("postgres", dbConfig)
    if err != nil {
		panic(err)
	}

    db.AutoMigrate(&Dashboard{})
    db.AutoMigrate(&User{})
    db.AutoMigrate(&Project{})
    db.AutoMigrate(&Collection{})
}

type Dashboard struct {
    gorm.Model
    Slug        string  `json:"Slug" sql:"unique_index" valid:"required,alphanum,length(5|10),lowercase"`
    Title       string  `json:"title" valid:"required,alphanum,length(1|50)"`
    Description string  `json:"description" valid:"required,alphanum"`
    Link        string  `json:"link" valid:"url"`
    Open        bool    `json:"open"`
    UserID      uint    `sql:"index"`
}

type Cover struct {
    ID     uint    `gorm:"primary_key"`
	Url    string  `valid:"url"`
}

type User struct {
    gorm.Model
    Name    string  `valid:"required,alphanum,length(1|50)"`
    Email   string  `valid:"required,email"`
    Avatar  string  `valid:"url"`
    Bio     string
}

type Project struct {
    gorm.Model
    Title           string  `valid:"required,alphanum,length(1|50)"`
    Description     string
    UserID          uint    `sql:"index"`
    Status          string
    Contributors    []User  `gorm:"many2many:project_contributors;"`
    Followers       []User  `gorm:"many2many:project_followers;"`
    Cover           string  `valid:"url"`
    Link            string  `valid:"url"`
    Tags            []Tag   `gorm:"many2many:project_tags;"`
    DashboardID     uint    `sql:"index"`
    Showcase        uint
}

type Tag struct {
	ID     uint   `gorm:"primary_key"`
	Value  string  `valid:"required,alphanum"`
}

type Collection struct {
    gorm.Model
    UserID      uint        `sql:"index"`
    Title       string      `valid:"required,alphanum,length(1|50)"`
    Description string
    Dashboards  []Dashboard `gorm:"many2many:collection_dashboards;"`
}
