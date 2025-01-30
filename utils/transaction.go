package utils

import (
    "context"
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
)

func HandleTransaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) (err error) {
    // เริ่มต้นทำการ transaction
    tx, err := db.BeginTxx(context.Background(), nil)
    // กรณีไม่สามารถเริ่มต้น transaction ได้
    if err != nil {
        log.Printf("❌ Error starting transaction: %v", err)
        return fmt.Errorf("❌ failed to start transaction: %w", err)
    }

    defer func() {
        // กรณีเกิด panic ระหว่างทำ transaction
        if r := recover(); r != nil {
            log.Printf("❌ Transaction panic: %v. Rolling back transaction.", r)
            tx.Rollback()
            panic(r)
        } else if err != nil {
            log.Printf("❌ Transaction error: %v. Rolling back transaction.", err)
            tx.Rollback()
        } else {
            err = tx.Commit()
            if err != nil {
                log.Printf("❌ Error committing transaction: %v", err)
                err = fmt.Errorf("❌ failed to commit transaction: %w", err)
            }
        }
    }()

    // เรียกใช้งานฟังก์ชันที่ผู้ใช้กำหนดด้วย -> (fn)
    if err = fn(tx); err != nil {
        // หากฟังก์ชันที่ส่งมาเกิดข้อผิดพลาด ให้ rollback transaction
        log.Printf("❌ Error during transaction execution: %v. Rolling back transaction.", err)
        return fmt.Errorf("❌ transaction failed: %w", err)
    }

    return nil
}