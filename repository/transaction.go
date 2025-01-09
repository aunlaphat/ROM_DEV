package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// ใช้จัดการ transaction ให้สามารถ rollback หรือ commit ได้อย่างเหมาะสม
func handleTransaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	// เริ่มต้นทำการ transaction
	tx, err := db.BeginTxx(context.Background(), nil)
    // กรณี cant begin transaction
	if err != nil {
        log.Printf("Error starting transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
        // กรณีเกิด panic ระหว่างทำ transaction
		if r := recover(); r != nil {
            log.Printf("Transaction panic: %v. Rolling back transaction.", r)
			tx.Rollback()
			panic(r)
		}
	}()

    // เรียกใช้งานฟังก์ชันที่ผู้ใช้กำหนดด้วย -> (fn)
	if err := fn(tx); err != nil {
        // หากฟังก์ชันที่ส่งมาเกิดข้อผิดพลาด ให้ rollback transaction
        log.Printf("Error during transaction execution: %v. Rolling back transaction.", err)
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	// หากไม่มีข้อผิดพลาด ให้ commit transaction
	if err := tx.Commit(); err != nil {
        // กรณี commit ไม่สำเร็จ
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

    log.Println("Transaction committed successfully.")
	return nil
}
