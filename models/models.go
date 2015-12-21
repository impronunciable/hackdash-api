package models

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"strconv"
)

var DB gorm.DB

func InitDB(dbConfig string) {
	var err error
	DB, err = gorm.Open("postgres", dbConfig)
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&Dashboard{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Project{})
	DB.AutoMigrate(&Collection{})
}

type Dashboard struct {
	gorm.Model
	Slug        string `json:"Slug" sql:"unique_index" valid:"required,alphanum,length(5|10),lowercase"`
	Title       string `json:"title" valid:"required,alphanum,length(1|50)"`
	Description string `json:"description" valid:"required,alphanum"`
	Link        string `json:"link" valid:"url"`
	Open        bool   `json:"open"`
	UserID      uint   `sql:"index"`
}

type Cover struct {
	ID  uint   `gorm:"primary_key"`
	Url string `valid:"url"`
}

type User struct {
	gorm.Model
	Name   string `valid:"required,alphanum,length(1|50)"`
	Email  string `valid:"required,email"`
	Avatar string `valid:"url"`
	Bio    string

	Auth0Id   string `valid:"required,alphanum,length(1|150)"`
	Provider  string `valid:"required,alphanum,length(1|50)"`
	ProviderId string `valid:"required,alphanum,length(1|100)"`
}

type Project struct {
	gorm.Model
	Title        string `valid:"required,alphanum,length(1|50)"`
	Description  string
	UserID       uint `sql:"index"`
	Status       string
	Contributors []User `gorm:"many2many:project_contributors;"`
	Followers    []User `gorm:"many2many:project_followers;"`
	Cover        string `valid:"url"`
	Link         string `valid:"url"`
	Tags         []Tag  `gorm:"many2many:project_tags;"`
	DashboardID  uint   `sql:"index"`
	Showcase     uint
}

type Tag struct {
	ID    uint   `gorm:"primary_key"`
	Value string `valid:"required,alphanum"`
}

type Collection struct {
	gorm.Model
	UserID      uint   `sql:"index"`
	Title       string `valid:"required,alphanum,length(1|50)"`
	Description string
	Dashboards  []Dashboard `gorm:"many2many:collection_dashboards;"`
}

// Model pagination
func Paginate(s *gorm.DB, c *echo.Context) *gorm.DB {
	page, _ := strconv.ParseUint(c.Query("page"), 10, 64)
	limit, _ := strconv.ParseUint(c.Query("limit"), 10, 64)

	// Check for edge case of page 0
	if page == 0 {
		page = 1
	}
	offset := page - 1
	offset = limit * offset
	return s.Limit(limit).Offset(offset)
}
