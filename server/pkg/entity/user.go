package entity

import (
	"time"
)

type Base struct {
	CreatedAt      time.Time `json:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	IsDeleted      bool      `json:"isDeleted"`
}

type User struct {
	Id             int    `gorm:"primaryKey" json:"id"`
	Name           string `json:"name"`
	Email          string `form:"unique" json:"email"`
	Password       string
	RoleId         *int
	Role           Role      `gorm:"constraint:OnDelete:SET NULL;`
	CreatedAt      time.Time `json:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	IsDeleted      bool      `json:"isDeleted"`
	CreatedBy      *string   `json:"createdBy"`
	UpdatedBy      *string   `json:"updatedBy"`
}

func (u *User) GetId() int {
	return u.Id
}

// type changeData struct {
// 	From string
// 	To   string
// }

// func (u *User) BeforeUpdate(scope *gorm.DB) (err error) {
// 	var originUser User
// 	key := fmt.Sprintf("%s.%v", utils.GetType(u),u.Id)
// 	test, ok := scope.Get()
// 	scope.Where("id = ?", u.Id).Find(&originUser)
// 	var test = make(map[string]changeData)
// 	for _, field := range scope.Statement.Schema.Fields {
// 		toValue := fmt.Sprintf("%v", reflect.ValueOf(*u).FieldByName(field.Name).Interface())
// 		newValue := fmt.Sprintf("%v", reflect.ValueOf(originUser).FieldByName(field.Name).Interface())
// 		if toValue != newValue {
// 			test[field.Name] = changeData{
// 				To:   toValue,
// 				From: newValue,
// 			}
// 		}
// 	}
// 	return nil
// }
