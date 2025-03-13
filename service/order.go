package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type OrderService interface {
	SearchOrder(ctx context.Context, req request.SearchOrder) (*response.SearchOrderResponse, error)
	CreateBeforeReturnOrder(ctx context.Context, req request.CreateBeforeReturnOrder, userID string) (*response.BeforeReturnOrderResponse, error)
	UpdateSrNo(ctx context.Context, orderNo string, srNo string, userID string) (*response.UpdateSrNoResponse, error)
	UpdateOrderStatus(ctx context.Context, orderNo string, userID string, roleID int) (*response.UpdateOrderStatusResponse, error)
	MarkOrderAsEdited(ctx context.Context, orderNo string, userID string) error
	CancelOrder(ctx context.Context, req request.CancelOrder, userID string) (*response.CancelOrderResponse, error)
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

	if len(req.Items) == 0 {
		err := errors.New("ต้องมีรายการสินค้าอย่างน้อย 1 รายการ")
		srv.logger.Warn("⚠️ No items provided", zap.Error(err))
		return nil, err
	}

	if req.ReturnDate.Before(time.Now()) {
		err := errors.New("วันที่คืนสินค้าต้องเป็นปัจจุบันหรืออนาคต")
		srv.logger.Warn("⚠️ Invalid ReturnDate", zap.Error(err))
		return nil, err
	}

	if req.SoStatus == "" {
		req.SoStatus = "open order"
	}
	if req.MkpStatus == "" {
		req.MkpStatus = "complete"
	}

	for i := range req.Items {
		req.Items[i].CreateBy = userID
	}

	err := srv.orderRepo.CreateBeforeReturnOrder(ctx, req, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to create BeforeReturnOrder", zap.Error(err))
		return nil, fmt.Errorf("failed to create return order: %w", err)
	}

	order, err := srv.orderRepo.GetBeforeReturnOrder(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch created BeforeReturnOrder", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve created order: %w", err)
	}

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

func (srv service) UpdateSrNo(ctx context.Context, orderNo string, srNo string, userID string) (*response.UpdateSrNoResponse, error) {
	srv.logger.Info("🔄 Updating SrNo...",
		zap.String("OrderNo", orderNo),
		zap.String("SrNo", srNo),
		zap.String("UpdatedBy", userID),
	)

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
		srNo, err = srv.requestSrNoFromAXAPI(ctx, orderNo) // 🔹 เรียก API ที่ AX
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

// 🔹 ฟังก์ชันเรียก API ไปที่ AX เพื่อขอ SrNo
func (srv service) requestSrNoFromAXAPI(ctx context.Context, orderNo string) (string, error) {
	// 🔹 เรียก API ภายในระบบ
	url := fmt.Sprintf("http://localhost:8080/api/order/generate-sr/%s", orderNo)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call generate SrNo API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("generate SrNo API failed with status: %s", resp.Status)
	}

	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("API returned error: %s", result.Message)
	}

	return result.Data, nil
}

func (srv service) UpdateOrderStatus(ctx context.Context, orderNo string, userID string, roleID int) (*response.UpdateOrderStatusResponse, error) {
	srv.logger.Info("🔄 Updating Order Status...",
		zap.String("OrderNo", orderNo),
		zap.String("RequestedBy", userID),
		zap.Int("RoleID", roleID),
	)

	// 🔹 ดึงข้อมูล BeforeReturnOrder
	order, err := srv.orderRepo.GetBeforeReturnOrder(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch BeforeReturnOrder", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	// ✅ ตรวจสอบสิทธิ์ของผู้ใช้
	switch roleID {
	case 2: // **Accounting**
		srv.logger.Info("🔹 Role: Accounting - Checking isCNCreated",
			zap.String("OrderNo", orderNo),
			zap.Bool("isCNCreated", order.IsCNCreated),
		)

		if !order.IsCNCreated {
			// 🔸 ถ้ายังไม่ได้สร้าง CN ให้สร้าง CN และอัปเดตสถานะ
			err = srv.CreateCNForOrder(ctx, orderNo, userID)
			if err != nil {
				srv.logger.Error("❌ Failed to create CN", zap.Error(err))
				return nil, fmt.Errorf("failed to create CN: %w", err)
			}

			return &response.UpdateOrderStatusResponse{
				OrderNo:        orderNo,
				StatusReturnID: 1, // Pending
				StatusConfID:   1, // Draft
				ConfirmBy:      userID,
				ConfirmDate:    time.Now(),
			}, nil
		}

		// 🔸 ถ้า CN ถูกสร้างแล้ว ให้อัปเดตเป็น Booking/Confirm
		err = srv.orderRepo.UpdateOrderStatus(ctx, orderNo, 3, 2, userID)
		if err != nil {
			srv.logger.Error("❌ Failed to update order status", zap.Error(err))
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}

		return &response.UpdateOrderStatusResponse{
			OrderNo:        orderNo,
			StatusReturnID: 3, // Booking
			StatusConfID:   2, // Confirm
			ConfirmBy:      userID,
			ConfirmDate:    time.Now(),
		}, nil

	case 3: // **Warehouse**
		srv.logger.Info("🔹 Role: Warehouse - Checking isEdited",
			zap.String("OrderNo", orderNo),
			zap.Bool("isEdited", order.IsEdited),
		)

		if order.IsEdited {
			// 🔸 ถ้ามีการแก้ไขข้อมูลก่อนยืนยัน ให้เปลี่ยนเป็น Pending/Draft
			err = srv.orderRepo.UpdateOrderStatus(ctx, orderNo, 1, 1, userID)
			if err != nil {
				srv.logger.Error("❌ Failed to update order status", zap.Error(err))
				return nil, fmt.Errorf("failed to update order status: %w", err)
			}

			return &response.UpdateOrderStatusResponse{
				OrderNo:        orderNo,
				StatusReturnID: 1, // Pending
				StatusConfID:   1, // Draft
				ConfirmBy:      userID,
				ConfirmDate:    time.Now(),
			}, nil
		}

		// 🔸 ถ้าไม่มีการแก้ไข ให้เปลี่ยนเป็น Booking/Confirm
		err = srv.orderRepo.UpdateOrderStatus(ctx, orderNo, 3, 2, userID)
		if err != nil {
			srv.logger.Error("❌ Failed to update order status", zap.Error(err))
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}

		return &response.UpdateOrderStatusResponse{
			OrderNo:        orderNo,
			StatusReturnID: 3, // Booking
			StatusConfID:   2, // Confirm
			ConfirmBy:      userID,
			ConfirmDate:    time.Now(),
		}, nil

	default: // **Role อื่น ๆ** → Default เป็น Pending/Draft
		srv.logger.Warn("⚠️ Unrecognized Role - Assigning Default Pending/Draft",
			zap.String("OrderNo", orderNo),
			zap.Int("RoleID", roleID),
		)

		err = srv.orderRepo.UpdateOrderStatus(ctx, orderNo, 1, 1, userID)
		if err != nil {
			srv.logger.Error("❌ Failed to update order status for unrecognized role", zap.Error(err))
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}

		return &response.UpdateOrderStatusResponse{
			OrderNo:        orderNo,
			StatusReturnID: 1, // Pending
			StatusConfID:   1, // Draft
			ConfirmBy:      userID,
			ConfirmDate:    time.Now(),
		}, nil
	}
}

func (srv service) CreateCNForOrder(ctx context.Context, orderNo string, userID string) error {
	srv.logger.Info("🔄 Creating CN...",
		zap.String("OrderNo", orderNo),
		zap.String("RequestedBy", userID),
	)

	// process create CN here...

	// 🔸 อัปเดต isCNCreated = true และเปลี่ยนสถานะเป็น Pending/Draft
	err := srv.orderRepo.UpdateCNForOrder(ctx, orderNo, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to update CN status", zap.Error(err))
		return fmt.Errorf("failed to update CN status: %w", err)
	}

	srv.logger.Info("✅ CN Created Successfully",
		zap.String("OrderNo", orderNo),
	)

	return nil
}

func (srv service) MarkOrderAsEdited(ctx context.Context, orderNo string, userID string) error {
	srv.logger.Info("✏️ Marking order as edited...",
		zap.String("OrderNo", orderNo),
		zap.String("UpdatedBy", userID),
	)

	err := srv.orderRepo.MarkOrderAsEdited(ctx, orderNo, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to mark order as edited", zap.Error(err))
		return fmt.Errorf("failed to mark order as edited: %w", err)
	}

	srv.logger.Info("✅ Order marked as edited", zap.String("OrderNo", orderNo))
	return nil
}

func (srv service) CancelOrder(ctx context.Context, req request.CancelOrder, userID string) (*response.CancelOrderResponse, error) {
	srv.logger.Info("🛑 Processing CancelOrder...",
		zap.String("RefID", req.RefID),
		zap.String("SourceTable", req.SourceTable),
		zap.String("CancelReason", req.CancelReason),
		zap.String("RequestedBy", userID),
	)

	if req.SourceTable != "BeforeReturnOrder" && req.SourceTable != "ReturnOrder" {
		srv.logger.Warn("⚠️ Invalid SourceTable", zap.String("SourceTable", req.SourceTable))
		return nil, fmt.Errorf("invalid SourceTable: %s", req.SourceTable)
	}

	statusReturnID, err := srv.orderRepo.GetReturnOrderStatus(ctx, req.RefID, req.SourceTable)
	if err != nil {
		srv.logger.Error("❌ Failed to retrieve order status",
			zap.String("RefID", req.RefID),
			zap.String("SourceTable", req.SourceTable),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve order status for RefID %s: %w", req.RefID, err)
	}

	const (
		StatusCancel    = 2 // ยกเลิก
		StatusUnsuccess = 5 // ไม่สำเร็จ
		StatusSuccess   = 6 // สำเร็จ
	)

	if statusReturnID == StatusCancel || statusReturnID == StatusUnsuccess || statusReturnID == StatusSuccess {
		srv.logger.Warn("⚠️ Order cannot be canceled due to current status",
			zap.String("RefID", req.RefID),
			zap.Int("StatusReturnID", statusReturnID),
		)
		return nil, fmt.Errorf("order cannot be canceled due to current status: %d", statusReturnID)
	}

	cancelID, err := srv.orderRepo.CancelOrder(ctx, req, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to cancel order",
			zap.String("RefID", req.RefID),
			zap.String("SourceTable", req.SourceTable),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to cancel order RefID %s: %w", req.RefID, err)
	}

	cancelDate := time.Now() // ควรจะดึงเวลาจริงๆ ใน db แต่ทำแบบนี้ไว้ก่อน

	srv.logger.Info("✅ Order canceled successfully",
		zap.Int("CancelID", cancelID),
		zap.String("RefID", req.RefID),
		zap.String("SourceTable", req.SourceTable),
		zap.String("CanceledBy", userID),
		zap.Time("CancelDate", cancelDate),
	)

	return &response.CancelOrderResponse{
		RefID:        req.RefID,
		SourceTable:  req.SourceTable,
		CancelReason: req.CancelReason,
		CancelBy:     userID,
		CancelDate:   cancelDate,
	}, nil
}
