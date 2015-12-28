package models

import (
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
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

// Basic model struct
type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Dashboard struct {
	Model
	Slug        string    `json:"slug" sql:"unique_index" valid:"required,alphanum,length(5|10),lowercase"`
	Title       string    `json:"title" valid:"required,alphanum,length(1|50)"`
	Description string    `json:"description" valid:"alphanum"`
	Link        string    `json:"link" valid:"url"`
	Open        bool      `json:"open"`
	Projects    []Project `json:"projects"`
	UserID      uint      `sql:"index" valid:"required"`
}

type Cover struct {
	ID  uint   `gorm:"primary_key"`
	Url string `valid:"url"`
}

type User struct {
	Model
	Name   string `valid:"required,alphanum,length(1|50)"`
	Email  string `valid:"required,email"`
	Avatar string `valid:"url"`
	Bio    string

	Auth0Id    string `valid:"required,alphanum,length(1|150)"`
	Provider   string `valid:"required,alphanum,length(1|50)"`
	ProviderId string `valid:"required,alphanum,length(1|100)"`
}

type Project struct {
	Model
	Title        string `json:"title" valid:"required,alphanum,length(1|50)"`
	Description  string `json:"description"`
	UserID       uint   `sql:"index" valid:"required"`
	Status       string `json:"status"`
	Contributors []User `json:"contributors" gorm:"many2many:project_contributors;"`
	Followers    []User `json:"followers" gorm:"many2many:project_followers;"`
	Cover        string `json:"cover" valid:"url"`
	Link         string `json:"link" valid:"url"`
	Tags         []Tag  `json:"tags" gorm:"many2many:project_tags;"`
	DashboardID  uint   `json:"dashboard_id" sql:"index" valid:"required"`
	Showcase     uint
}

type Tag struct {
	ID    uint   `gorm:"primary_key"`
	Value string `valid:"required,alphanum"`
}

type Collection struct {
	Model
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

func (d *Dashboard) BeforeSave() (err error) {
	_, err = govalidator.ValidateStruct(d)
	return
}

func (p *Project) BeforeSave() (err error) {
	_, err = govalidator.ValidateStruct(p)
	return
}
