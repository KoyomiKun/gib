package model

import (
	"Blog/util/errmsg"
	"errors"

	"github.com/jinzhu/gorm"
)

type Post struct {
	Category Category `gorm:"foreignKey:Cid" json:"category,omitempty"`
	gorm.Model
	Title string `gorm:"type:varchar(100);not null" json:"title,omitempty"`
	Cid   int    `json:"cid,omitempty"`
	Desc  string `gorm:"type:varchar(200);" json:"desc,omitempty"`
	// 试试[]string怎么处理
	// 会引发panic, 可恶
	Tags    string `gorm:"type:varchar(300);" json:"tags,omitempty"`
	Content string `gorm:"type:longtext;" json:"content,omitempty"`
}

// 新增文章
func InsertPost(post *Post) errmsg.Code {
	//post.Password = scryptPassword(post.Password)
	err := db.Create(&post).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 查询单个文章
func GetPost(id int) (*Post, errmsg.Code) {
	var p Post
	err := db.Preload("Category").First(&p, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &p, errmsg.ERROR_POST_NOT_FOUND
		} else {
			logger.Debug("Fail to get post:", err)
			return &p, errmsg.ERROR
		}
	}
	return &p, errmsg.SUCCESS
}

// 查询分类下所有文章
func GetPostByCate(cid, pageSize, pageNum int) ([]Post, errmsg.Code) {
	var posts []Post
	err := db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where("cid = ?", cid).Find(&posts).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return posts, errmsg.SUCCESS
}

// 查询文章表
func GetPostsList(pageSize int, pageNum int) ([]Post, errmsg.Code) {
	var posts []Post
	// limit: 指定需要多少条； offset：指定从哪一条开始
	err := db.Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&posts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR
	}
	return posts, errmsg.SUCCESS
}

// 编辑文章信息
func EditPost(post *Post) errmsg.Code {
	logger.Debug("Edit post:", post)
	err := db.Model(post).Updates(post).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除文章
func DeletePost(id int) errmsg.Code {
	if err = db.Delete(&Post{}, id).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
