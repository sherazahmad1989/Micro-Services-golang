package seed

import (
	"fullstack/api/models"
	"github.com/jinzhu/gorm"
	"log"
)

var users = []models.User{
	models.User{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}



var addresses = []models.Address{
	models.Address{
		AddressName:   "Title 1",
		LocationPath: "Hello world 1",
	},
	models.Address{
		AddressName:   "Title 2",
		LocationPath: "Hello world 2",
	},
}



func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.Address{} ,&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Address{}  ,&models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
		addresses[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Address{}).Create(&addresses[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Address table: %v", err)
		}
	}
}