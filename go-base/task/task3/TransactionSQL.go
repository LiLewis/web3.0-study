package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Account 表
type Account struct {
	ID      uint  `gorm:"primaryKey"`
	Balance int64 // 账户余额
}

// Transaction 表
type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID uint
	ToAccountID   uint
	Amount        int64
}

func main() {
	//数据库连接
	dsn := "root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	//创建表
	db.AutoMigrate(&Account{}, &Transaction{})

	//转账200米
	var accountA, accountB Account
	db.FirstOrCreate(&accountA, Account{ID: 1, Balance: 1000})
	db.FirstOrCreate(&accountB, Account{ID: 2, Balance: 30})

	//转账事务
	err = db.Transaction(func(tx *gorm.DB) error {
		//查询A账户
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&accountA, 1).Error; err != nil {
			return err
		}

		// 2. 检查余额是否足够
		if accountA.Balance < 100 {
			return fmt.Errorf("账户 A 余额不足")
		}

		// 3. 扣除 A 的余额
		if err := tx.Model(&Account{}).Where("id = ?", accountA.ID).
			Update("balance", accountA.Balance-100).Error; err != nil {
			return err
		}

		// 4. 查询账户 B（加锁，避免并发问题）
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&accountB, 2).Error; err != nil {
			return err
		}

		// 5. 增加 B 的余额
		if err := tx.Model(&Account{}).Where("id = ?", accountB.ID).
			Update("balance", accountB.Balance+100).Error; err != nil {
			return err
		}

		// 6. 记录交易流水
		transaction := Transaction{
			FromAccountID: accountA.ID,
			ToAccountID:   accountB.ID,
			Amount:        100,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// 提交事务
		return nil
	})

	if err != nil {
		fmt.Println("转账失败，事务已回滚:", err)
	} else {
		fmt.Println("转账成功！")
	}

}
