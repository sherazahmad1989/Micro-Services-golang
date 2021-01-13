package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Address struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	AddressName     string    `gorm:"size:255;not null;unique" json:"address_name"`
	LocationPath   string    `gorm:"size:255;not null;" json:"location_path"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (a *Address) Prepare() {
	a.ID = 0
	a.AddressName = html.EscapeString(strings.TrimSpace(a.AddressName))
	a.LocationPath = html.EscapeString(strings.TrimSpace(a.LocationPath))
	a.Author = User{}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *Address) Validate() error {

	if a.AddressName == "" {
		return errors.New("Required AddressName")
	}
	if a.LocationPath == "" {
		return errors.New("Required LocationPath")
	}
	return nil
}

func (a *Address) SaveAddress(db *gorm.DB) (*Address, error) {
	var err error
	err = db.Debug().Model(&Address{}).Create(&a).Error
	if err != nil {
		return &Address{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.AuthorID).Take(&a.Author).Error
		if err != nil {
			return &Address{}, err
		}
	}
	return a, nil
}

func (a *Address) FindAllAddress(db *gorm.DB) (*[]Address, error) {
	var err error
	posts := []Address{}
	err = db.Debug().Model(&Address{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Address{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Address{}, err
			}
		}
	}
	return &posts, nil
}

func (a *Address) FindAddressByID(db *gorm.DB, pid uint64) (*Address, error) {
	var err error
	err = db.Debug().Model(&Address{}).Where("id = ?", pid).Take(&a).Error
	if err != nil {
		return &Address{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.AuthorID).Take(&a.Author).Error
		if err != nil {
			return &Address{}, err
		}
	}
	return a, nil
}

func (a *Address) UpdateAAddress(db *gorm.DB) (*Address, error) {

	var err error

	err = db.Debug().Model(&Address{}).Where("id = ?", a.ID).Updates(Address{AddressName: a.AddressName, LocationPath: a.LocationPath, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Address{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.AuthorID).Take(&a.Author).Error
		if err != nil {
			return &Address{}, err
		}
	}
	return a, nil
}

func (a *Address) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Address{}).Where("id = ? and author_id = ?", pid, uid).Take(&Address{}).Delete(&Address{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Address not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}