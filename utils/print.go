package utils

import (
	res "boilerplate-backend-go/dto/response"
	"fmt"
)

func PrintOrderDetails(order *res.BeforeReturnOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", order.SrNo)
	fmt.Printf("📡 ChannelID: %d\n", order.ChannelID)
	fmt.Printf("🎐 Reason: %s\n", order.Reson)
	fmt.Printf("👤 CustomerID: %s\n", order.CustomerID)
	fmt.Printf("📦 TrackingNo: %s\n", order.TrackingNo)
	fmt.Printf("🚚 Logistic: %s\n", order.Logistic)
	fmt.Printf("🏢 WarehouseID: %d\n", order.WarehouseID)
	fmt.Printf("📄 SoStatusID: %v\n", order.SoStatusID)
	fmt.Printf("📊 MkpStatusID: %v\n", order.MkpStatusID)
	fmt.Printf("📅 ReturnDate: %v\n", order.ReturnDate)
	fmt.Printf("🪃 StatusReturnID: %d\n", order.StatusReturnID)
	fmt.Printf("✅ StatusConfID: %d\n", order.StatusConfID)
	fmt.Printf("👤 ConfirmBy: %v\n", order.ConfirmBy)
	fmt.Printf("👤 CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
	fmt.Printf("👤 UpdateBy: %v\n", order.UpdateBy)
	fmt.Printf("📅 UpdateDate: %v\n", order.UpdateDate)
	fmt.Printf("❌ CancelID: %v\n", order.CancelID)
}

func PrintOrderLineDetails(line *res.BeforeReturnOrderLineResponse) {
	fmt.Printf("🔖 SKU: %s\n", line.SKU)
	fmt.Printf("🏷️  ItemName: %s\n", line.ItemName)
	fmt.Printf("📱 QTY: %d\n", line.QTY)
	fmt.Printf("📲 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("💲 Price: %.2f ฿\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %v\n", line.CreateDate)
}

func PrintSaleOrderDetails(order *res.SaleOrderResponse) {
	fmt.Printf("🧾 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("📊 StatusMKP: %s\n", order.StatusMKP)
	fmt.Printf("📄 SalesStatus: %s\n", order.SalesStatus)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
}

func PrintSaleOrderLineDetails(line *res.SaleOrderLineResponse) {
	fmt.Printf("🛒 SoNo: %s\n", line.SoNo)
	fmt.Printf("🧾 OrderNo: %s\n", line.OrderNo)
	fmt.Printf("🔖 SKU: %s\n", line.SKU)
	fmt.Printf("🏷️  ItemName: %s\n", line.ItemName)
	fmt.Printf("📱 QTY: %d\n", line.QTY)
	fmt.Printf("💲 Price: %.2f ฿\n", line.Price)
}

// ************************ Draft & Confirm ************************ //

// Draft & Confirm Head
func PrintDraftConfirmOrderDetails(draft *res.ListDraftConfirmOrdersResponse) {
	fmt.Printf("📦 OrderNo: %s\n", draft.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", draft.SoNo)
	fmt.Printf("👤 Customer: %s\n", draft.CustomerID)
	fmt.Printf("🔄 SrNo: %s\n", draft.SrNo)
	fmt.Printf("📦 TrackingNo: %s\n", draft.TrackingNo)
	fmt.Printf("🚚 Logistic: %s\n", draft.Logistic)
	fmt.Printf("📡 Channel: %d\n", draft.ChannelID)
	fmt.Printf("📅 CreateDate: %v\n", draft.CreateDate)
	fmt.Printf("🏢 Warehouse: %d\n", draft.WarehouseID)
}

// Modal Edit Draft & Modal Show Confirm
func PrintDraftOrderDetails(draft *res.DraftHeadResponse) {
	fmt.Printf("📦 OrderNo: %s\n", draft.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", draft.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", draft.SrNo)
}

func PrintDraftOrderLineDetails(draft *res.DraftLineResponse) {
	fmt.Printf("🔖 SKU: %s\n", draft.SKU)
	fmt.Printf("🏷️  ItemName: %s\n", draft.ItemName)
	fmt.Printf("📱 QTY: %d\n", draft.QTY)
	fmt.Printf("💲 Price: %.2f\n", draft.Price)
}