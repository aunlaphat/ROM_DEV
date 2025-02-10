package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type OrderService interface {
	SearchOrder(ctx context.Context, req request.SearchOrder) (*response.SearchOrderResponse, error)
	CreateBeforeReturnOrder(ctx context.Context, req request.CreateBeforeReturnOrder, userID string) (*response.BeforeReturnOrderResponse, error)
	UpdateSrNo(ctx context.Context, orderNo string, userID string) (*response.UpdateSrNoResponse, error)
}

func (srv service) SearchOrder(ctx context.Context, req request.SearchOrder) (*response.SearchOrderResponse, error) {
	srv.logger.Info("🔎 Searching for Order",
		zap.String("SoNo", req.SoNo),
		zap.String("OrderNo", req.OrderNo),
	)

	if req.SoNo == "" && req.OrderNo == "" {
		err := errors.New("either SoNo or OrderNo must be provided")
		srv.logger.Warn("⚠️ Invalid request - Missing parameters", zap.Error(err))
		return nil, err
	}

	order, err := srv.orderRepo.SearchOrder(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.logger.Warn("⚠️ No Sale Order found",
				zap.String("SoNo", req.SoNo),
				zap.String("OrderNo", req.OrderNo),
			)
			return nil, sql.ErrNoRows
		}

		srv.logger.Error("❌ Failed to search Order",
			zap.String("SoNo", req.SoNo),
			zap.String("OrderNo", req.OrderNo),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	srv.logger.Info("✅ Order found",
		zap.String("SoNo", order.SoNo),
		zap.String("OrderNo", order.OrderNo),
		zap.Int("TotalItems", len(order.Items)),
	)

	return order, nil
}

func (srv service) CreateBeforeReturnOrder(ctx context.Context, req request.CreateBeforeReturnOrder, userID string) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("📝 Creating BeforeReturnOrder",
		zap.String("OrderNo", req.OrderNo),
		zap.String("SoNo", req.SoNo),
		zap.Int("TotalItems", len(req.Items)),
		zap.String("CreateBy", userID),
	)

	// ✅ ตรวจสอบว่ามีสินค้าคืนหรือไม่
	if len(req.Items) == 0 {
		err := errors.New("ต้องมีรายการสินค้าอย่างน้อย 1 รายการ")
		srv.logger.Warn("⚠️ No items provided", zap.Error(err))
		return nil, err
	}

	// ✅ ตรวจสอบ `ReturnDate` ต้องไม่เป็นอดีต
	if req.ReturnDate.Before(time.Now()) {
		err := errors.New("วันที่คืนสินค้าต้องเป็นปัจจุบันหรืออนาคต")
		srv.logger.Warn("⚠️ Invalid ReturnDate", zap.Error(err))
		return nil, err
	}

	// ✅ ตั้งค่า Default `SoStatus` และ `MkpStatus`
	if req.SoStatus == "" {
		req.SoStatus = "open order"
	}
	if req.MkpStatus == "" {
		req.MkpStatus = "complete"
	}

	// ✅ กำหนด `CreateBy`
	for i := range req.Items {
		req.Items[i].CreateBy = userID
	}

	// 🔹 เรียกใช้ Repository Layer
	err := srv.orderRepo.CreateBeforeReturnOrder(ctx, req, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to create BeforeReturnOrder", zap.Error(err))
		return nil, fmt.Errorf("failed to create return order: %w", err)
	}

	// ✅ ดึงข้อมูลที่เพิ่งสร้าง
	order, err := srv.orderRepo.GetBeforeReturnOrder(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch created BeforeReturnOrder", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve created order: %w", err)
	}

	// ✅ ดึงข้อมูลรายการสินค้า
	items, err := srv.orderRepo.GetBeforeReturnOrderItems(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch created BeforeReturnOrderItems", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve created order items: %w", err)
	}

	order.Items = items

	srv.logger.Info("✅ BeforeReturnOrder created successfully",
		zap.String("OrderNo", order.OrderNo),
		zap.String("SoNo", order.SoNo),
		zap.Int("TotalItems", len(order.Items)),
	)

	return order, nil
}

func (srv service) UpdateSrNo(ctx context.Context, orderNo string, userID string) (*response.UpdateSrNoResponse, error) {
	srv.logger.Info("🔄 Requesting SrNo from AX...",
		zap.String("OrderNo", orderNo),
		zap.String("RequestedBy", userID),
	)

	srNo, err := srv.GenerateSrNoFromAX(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to generate SrNo", zap.Error(err))
		return nil, fmt.Errorf("failed to generate SrNo: %w", err)
	}

	resp, err := srv.orderRepo.UpdateSrNo(ctx, orderNo, srNo, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to update SrNo in DB", zap.Error(err))
		return nil, fmt.Errorf("failed to update SrNo in DB: %w", err)
	}

	srv.logger.Info("✅ SrNo updated successfully",
		zap.String("OrderNo", resp.OrderNo),
		zap.String("SrNo", resp.SrNo),
	)

	return resp, nil
}

// 🔹 ฟังก์ชันที่จำลองการสร้าง SrNo จาก AX พร้อม Retry 3 ครั้ง
func (srv service) GenerateSrNoFromAX(ctx context.Context, orderNo string) (string, error) {
	maxRetries := 3 // 🔄 กำหนดจำนวนครั้งที่ retry
	var srNo string
	var err error

	for i := 1; i <= maxRetries; i++ {
		srNo, err = requestSrNoFromAX(orderNo) // 🔹 เรียก API ที่ AX
		if err == nil {
			return srNo, nil // ✅ สำเร็จ ออกจากลูป
		}

		srv.logger.Warn("⚠️ Failed to request SrNo from AX",
			zap.String("OrderNo", orderNo),
			zap.Int("RetryAttempt", i),
			zap.Error(err),
		)

		// ❌ ถ้าเป็นรอบสุดท้ายของ retry ให้คืน error
		if i == maxRetries {
			break
		}

		// ⏳ รอ 2 วินาทีก่อน retry ใหม่
		time.Sleep(2 * time.Second)
	}

	return "", fmt.Errorf("failed to request SrNo from AX after %d retries", maxRetries)
}

// 🔹 ฟังก์ชันจำลองการเรียก API ไปที่ AX เพื่อขอ SrNo
func requestSrNoFromAX(orderNo string) (string, error) {
	// 🔹 จำลองเลข SrNo (จริง ๆ ต้องเรียก API AX)
	fakeSrNo := fmt.Sprintf("SR-%s-%d", orderNo, time.Now().Unix())

	// ❌ จำลองความล้มเหลวแบบสุ่ม 5%
	if rand.Intn(100) < 5 {
		return "", errors.New("AX API error - SrNo request failed")
	}

	return fakeSrNo, nil
}
