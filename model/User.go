package model

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"ginblog/utils/errmsg"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"log"
)

type User struct {
	//id 创建时间，更新时间
	gorm.Model
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(255);not null" json:"password" validate:"required,min=6,max=20" label:"密码""`
	Salt     string `gorm:"type:varchar(16);not null" json:"salt"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

//查询用户是否存在

func CheckUser(name string) (code int) {
	var users User
	//这段代码，以id来排序，where来设定查询的限制，first寻找第一个匹配的记录
	db.Select("id").Where("username = ?", name).First(&users)
	//如果名字有记录，则再Model函数里，id会设置成正整数，注册过的用户，如果被软删除了，它依然会被占用
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001用户名存在
	}
	return errmsg.SUCCESS
}

// 注册用户

func CreateUser(data *User) int {
	//加密密码
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		log.Fatal(err)
	}
	data.Salt = base64.StdEncoding.EncodeToString(salt) // 存储盐值
	data.Password = ScryptPw(data.Password, data.Salt)  // 使用盐值加密密码
	err := db.Create(&data).Error
	if err != nil {
		log.Printf("Error CreateUser: %v", err)
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//查询用户列表

func GetUsers(pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录，可以返回空的分类列表而不是nil
			return []User{}, 0
		}
		// 如果有其他错误，记录日志并处理错误
		log.Printf("Error GetCate: %v", err)
		return nil, 0
	}
	return users, total
}

//编辑用户(密码以外的信息，密码单独做一个功能)

func EditUser(id int, data *User) int {
	var user User
	// 首先检查用户ID是否存在
	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果用户不存在，返回错误码
			return errmsg.ERROR_USER_NOT_EXIST
		}
		// 如果查询过程中出现其他错误，也返回错误码
		return errmsg.ERROR
	}
	// 如果用户存在，继续更新操作
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err = db.Model(&user).Where("id= ? ", id).Updates(maps).Error
	if err != nil {
		log.Printf("Error EditUser: %v", err)
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除用户
func DeleteUser(id int) int {
	var user User
	db.Where("id = ?", id).Delete(&user)
	if err != nil {
		log.Printf("Error DeleteUser: %v", err)
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 密码加密

// 钩子函数
// 这是特定的方法，方法名必须按照规范的写，不然需要和之前的一样要调用这个方法
//func (u *User) BeforeSave() {
//	u.Password = ScryptPw(u.Password)
//}

func ScryptPw(password, salt string) string {
	const KeyLen = 32
	// 使用传入的盐值而不是生成新的盐值
	HashPw, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}

// 登录验证+

func CheckLogin(username string, password string) int {
	var user User
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	//判断密码
	if ScryptPw(password, user.Salt) != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	//判断管理权限
	if user.Role != 1 {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	return errmsg.SUCCESS
}
