package model

import (
	"Blog/util/errmsg"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
)

type User struct {
	gorm.Model
	Username    string `gorm:"type:varchar(20);unique;not null" json:"username,omitempty" validate:"required,min=4,max=12" label:"用户名"`
	Email       string `gorm:"type:varchar(100)" json:"email,omitempty"`
	Phonenumber string `gorm:"type:varchar(30)" json:"phonenumber,omitempty"`
	Password    string `gorm:"type:varchar(20);not null" json:"password,omitempty"  validate:"required,min=6,max=20" label:"密码"`
	Role        int    `gorm:"type:int" json:"role,omitempty"`
	Avatar      string `gorm:"type:varchar(100)" json:"avatar,omitempty"`
}

// 查询用户是否存在
// 存在返回id 和success，不存在返回0和error
func ExistsUser(username string) (uint, errmsg.Code) {
	var u User
	db.Select("id").Where("username = ?", username).First(&u)
	if u.ID > 0 {
		return u.ID, errmsg.SUCCESS
	}
	return u.ID, errmsg.ERROR_USER_NOT_EXIST
}

// 查询单个用户信息
func SelectUser(id uint) (*User, errmsg.Code) {
	var u User
	if err := db.First(&u, id).Error; err != nil {
		return nil, errmsg.ERROR_USER_NOT_EXIST
	} else {
		return &u, errmsg.SUCCESS
	}
}

// 新增用户
func InsertUser(user *User) errmsg.Code {
	//user.Password = scryptPassword(user.Password)
	err := db.Create(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Password, err = scryptPassword(u.Password)
	logger.Info("Exec hook before create:", u.Username)
	if err != nil {
		return err
	}
	return nil
}

// 查询用户表
func GetUsersList(pageSize int, pageNum int) ([]User, errmsg.Code) {
	var users []User
	// limit: 指定需要多少条； offset：指定从哪一条开始
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR
	}
	return users, errmsg.SUCCESS

}

// 编辑用户信息
func EditUser(id uint, user *User) errmsg.Code {
	//dataMap := util.Struct2Map(*user)
	//delete(dataMap, "password")
	//logger.Debug("Edit user:", dataMap)
	logger.Debug("Edit user:", id, user)
	err := db.Model(&User{}).Where("id = ?", id).Update(user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除用户
func DeleteUser(id int) errmsg.Code {
	var u User
	db.First(&u, id)
	if u.ID == 0 {
		return errmsg.ERROR
	}
	now := time.Now().Unix()
	old_name := u.Username
	//在这里如果用GORM的update会卡住, 待研究
	//删除和更新必须都要在一个事务中执行，不然要是更新了结果删除失败名字就变了
	tx := db.Begin()
	err := tx.Exec("update user set username = ? where id = ?", old_name+"_"+strconv.Itoa(int(now)), u.ID).Error
	logger.Debug("Finish update")
	if err != nil {
		logger.Error("Fail to update when delete user:", err.Error())
	}
	err = tx.Delete(&u).Error
	if err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}
	tx.Commit()
	return errmsg.SUCCESS
}

// 密码加密
func scryptPassword(password string) (string, error) {
	const KEYLEN = 10
	salt := []byte{12, 32, 4, 6, 66, 22, 222, 11}

	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KEYLEN)
	if err != nil {
		logger.Error("Scrypte Error:", err.Error())
		return "", err
	}
	return base64.StdEncoding.EncodeToString(HashPw), nil
}

// 登录验证
func Login(username, password string) errmsg.Code {
	var user User
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	if secPass, _ := scryptPassword(password); secPass != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	return errmsg.SUCCESS
}
