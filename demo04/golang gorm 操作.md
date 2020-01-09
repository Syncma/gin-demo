# golang gorm 操作

[toc]

## gorm例子

```go
package main

import (
	"fmt"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	//这里导入包使用了 _ 前缀代表仅仅是导入这个包，但是我们在代码里面不会直接使用
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const DSN = "root:123456@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local"
const DRIVER = "mysql"

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(DRIVER, DSN)
	if err != nil {
		panic(err)
	}
}

//定义User模型，绑定users表，ORM库操作数据库，需要定义一个struct类型和MYSQL表进行绑定或者叫映射，struct字段和MYSQL表字段一一对应
//在这里User类型可以代表mysql users表
type User struct {
	ID       int    `gorm:"primary_key"` //表字段首字母要大写？
	Username string `gorm:"type:varchar(20);not null;"`
	Password string `gorm:"type:varchar(50);not null;"`
}

//设置表名，可以通过给struct类型定义 TableName函数，返回当前struct绑定的mysql表名是什么
//如果不设置表名，默认是users
// func (u User) TableName() string {
// 	//绑定MYSQL表名为users
// 	return "user_info"
// }

func main() {
	defer db.Close()

	//创建表
	if !db.HasTable(&User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
			panic(err)
		}
	}

	//增加记录
	// user := User{
	// 	Username: "tony",
	// 	Password: "123456",
	// }

	// if err := db.Create(&user).Error; err != nil {
	// 	fmt.Println("插入失败")
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("插入成功")
	// }

	//查询记录
	u := User{}
	isNotFound := db.Where("username=?", "tony").First(&u).RecordNotFound()
	if isNotFound {
		fmt.Println("找不到记录")
	} else {
		fmt.Println("找到记录")
	}

	//删除记录
	// var user User
	// db.Where("username=?", "tony").Delete(&user)

	//修改记录
	var user User
	// db.Model(&user).Update("username", "jacky")
	err := db.Model(&user).Where("username=?", "tony").Update("password", "111111")
	if err != nil {
		fmt.Println("修改成功")
	} else {
		fmt.Println("修改失败")
	}

	//事务操作

}
```

[更多例子](https://www.tizi365.com/archives/26.html)


