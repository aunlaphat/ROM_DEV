package utils

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func HandleTransaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) (err error) {
	// *️⃣ เริ่มต้นทำการ transaction
	tx, err := db.BeginTxx(context.Background(), nil)
	// *️⃣ กรณีไม่สามารถเริ่มต้น transaction ได้
	if err != nil {
		return fmt.Errorf("[ Failed to start transaction: %w ]", err)
	}

	defer func() {
		// *️⃣ กรณีเกิด panic ระหว่างทำ transaction
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("[ Failed to commit transaction: %w ]", err)
			}
		}
	}()

	// *️⃣ เรียกใช้งานฟังก์ชันที่ผู้ใช้กำหนดด้วย -> (fn)
	if err = fn(tx); err != nil {
		// *️⃣ หากฟังก์ชันที่ส่งมาเกิดข้อผิดพลาด ให้ rollback transaction
		return fmt.Errorf("[ transaction failed: %w ]", err)
	}

	return nil
}
