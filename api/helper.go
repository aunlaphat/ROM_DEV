package api

import (
	"fmt"
	"net/http"
	"strconv"
	res "boilerplate-backend-go/dto/response"
)

// Helper function: parsePagination
func parsePagination(r *http.Request) (int, int) {
	query := r.URL.Query()
	page := parseInt(query.Get("page"), 1)    // Default page = 1
	limit := parseInt(query.Get("limit"), 10) // Default limit = 10
	return page, limit
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// Helper function สำหรับดึง userID จาก claims
func getUserIDFromClaims(claims map[string]interface{}) (string, error) {
	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid user information in token")
	}
	return userID, nil
}

// Helper function of form response
func printOrderDetails(order *res.BeforeReturnOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", order.SrNo)
	fmt.Printf("📡 ChannelID: %d\n", order.ChannelID)
	fmt.Printf("🔙 Reason: %s\n", order.Reason)
	fmt.Printf("👤 CustomerID: %s\n", order.CustomerID)
	fmt.Printf("📦 TrackingNo: %s\n", order.TrackingNo)
	fmt.Printf("🚚 Logistic: %s\n", order.Logistic)
	fmt.Printf("🏢 WarehouseID: %d\n", order.WarehouseID)
	fmt.Printf("📄 SoStatusID: %v\n", order.SoStatusID)
	fmt.Printf("📊 MkpStatusID: %v\n", order.MkpStatusID)
	fmt.Printf("📅 ReturnDate: %v\n", order.ReturnDate)
	fmt.Printf("🔖 StatusReturnID: %d\n", order.StatusReturnID)
	fmt.Printf("✅ StatusConfID: %d\n", order.StatusConfID)
	fmt.Printf("👤 ConfirmBy: %v\n", order.ConfirmBy)
	fmt.Printf("👤 CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
	fmt.Printf("👤 UpdateBy: %v\n", order.UpdateBy)
	fmt.Printf("📅 UpdateDate: %v\n", order.UpdateDate)
	fmt.Printf("❌ CancelID: %v\n", order.CancelID)
}

func printOrderLineDetails(line *res.BeforeReturnOrderLineResponse) {
	fmt.Printf("🔢 SKU: %s\n", line.SKU)
	fmt.Printf("🔢 QTY: %d\n", line.QTY)
	fmt.Printf("🔢 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("💲 Price: %.2f\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %v\n", line.CreateDate)
}

func printSaleOrderDetails(order *res.SaleOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🔢 SoNo: %s\n", order.SoNo)
	fmt.Printf("📊 StatusMKP: %s\n", order.StatusMKP)
	fmt.Printf("📊 SalesStatus: %s\n", order.SalesStatus)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
}

func printSaleOrderLineDetails(line *res.SaleOrderLineResponse) {
	fmt.Printf("🔢 SKU: %s\n", line.SKU)
	fmt.Printf("🚩 ItemName: %s\n", line.ItemName)
	fmt.Printf("🔢 QTY: %d\n", line.QTY)
	fmt.Printf("💲 Price: %.2f\n", line.Price)
}

func printDraftDetails(draft *res.BeforeReturnOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", draft.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", draft.SoNo)
	fmt.Printf("👤 Customer: %s\n", draft.CustomerID)
	fmt.Printf("🔄 SrNo: %s\n", draft.SrNo)
	fmt.Printf("📦 TrackingNo: %s\n", draft.TrackingNo)
	fmt.Printf("📡 Channel: %d\n", draft.ChannelID)
	fmt.Printf("📅 CreateDate: %v\n", draft.CreateDate)
	fmt.Printf("🏢 Warehouse: %d\n", draft.WarehouseID)
}

func printDraftLineDetails(line *res.BeforeReturnOrderLineResponse) {
	fmt.Printf("🔢 SKU: %s\n", line.SKU)
	fmt.Printf("🔢 QTY: %d\n", line.QTY)
	fmt.Printf("🔢 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("💲 Price: %.2f\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %v\n", line.CreateDate)
}
