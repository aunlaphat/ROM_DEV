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
	fmt.Printf("🎐 Reason: %s\n", order.Reason)
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

func PrintReturnOrderDetails(order *res.ReturnOrder) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %v\n", order.SrNo)
	fmt.Printf("📡 ChannelID: %d\n", order.ChannelID)
	fmt.Printf("🎐 Reason: %v\n", order.Reason)
	// fmt.Printf("👤 CustomerID: %s\n", order.CustomerID)
	fmt.Printf("📦 TrackingNo: %v\n", order.TrackingNo)
	// fmt.Printf("🚚 Logistic: %s\n", order.Logistic)
	// fmt.Printf("🏢 WarehouseID: %d\n", order.WarehouseID)
	// fmt.Printf("📄 SoStatusID: %v\n", order.SoStatusID)
	// fmt.Printf("📊 MkpStatusID: %v\n", order.MkpStatusID)
	// fmt.Printf("📅 ReturnDate: %v\n", order.ReturnDate)
	// fmt.Printf("🪃 StatusReturnID: %d\n", order.StatusReturnID)
	// fmt.Printf("✅ StatusConfID: %d\n", order.StatusConfID)
	// fmt.Printf("👤 ConfirmBy: %v\n", order.ConfirmBy)
	fmt.Printf("👤 CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
	fmt.Printf("👤 UpdateBy: %v\n", order.UpdateBy)
	fmt.Printf("📅 UpdateDate: %v\n", order.UpdateDate)
	// fmt.Printf("❌ CancelID: %v\n", order.CancelID)
}

func PrintCreateReturnOrder(order *res.CreateReturnOrder) {
    fmt.Printf("OrderNo: %s\n", order.OrderNo)
    fmt.Printf("SoNo: %s\n", order.SoNo)
    fmt.Printf("SrNo: %s\n", order.SrNo)
    fmt.Printf("TrackingNo: %s\n", order.TrackingNo)
    fmt.Printf("PlatfID: %v\n", order.PlatfID)
    fmt.Printf("ChannelID: %v\n", order.ChannelID)
    fmt.Printf("OptStatusID: %v\n", order.OptStatusID)
    fmt.Printf("AxStatusID: %v\n", order.AxStatusID)
    fmt.Printf("PlatfStatusID: %v\n", order.PlatfStatusID)
    fmt.Printf("Reason: %v\n", order.Reason)
    fmt.Printf("CancelID: %v\n", order.CancelID)
    fmt.Printf("StatusCheckID: %v\n", order.StatusCheckID)
    fmt.Printf("CheckBy: %v\n", order.CheckBy)
    fmt.Printf("Description: %v\n", order.Description)
    fmt.Printf("CreateBy: %s\n", order.CreateBy)
    fmt.Printf("CreateDate: %s\n", order.CreateDate)
}

func PrintUpdateReturnOrder(order *res.UpdateReturnOrder) {
    fmt.Printf("OrderNo: %s\n", order.OrderNo)
    fmt.Printf("SoNo: %s\n", order.SoNo)
    fmt.Printf("SrNo: %s\n", order.SrNo)
    fmt.Printf("TrackingNo: %s\n", order.TrackingNo)
    fmt.Printf("PlatfID: %d\n", order.PlatfID)
    fmt.Printf("ChannelID: %d\n", order.ChannelID)
    fmt.Printf("OptStatusID: %d\n", order.OptStatusID)
    fmt.Printf("AxStatusID: %d\n", order.AxStatusID)
    fmt.Printf("PlatfStatusID: %d\n", order.PlatfStatusID)
    fmt.Printf("Reason: %s\n", order.Reason)
    fmt.Printf("CancelID: %d\n", order.CancelID)
    fmt.Printf("StatusCheckID: %d\n", order.StatusCheckID)
    fmt.Printf("CheckBy: %s\n", order.CheckBy)
    fmt.Printf("Description: %s\n", order.Description)
    fmt.Printf("UpdateBy: %s\n", order.UpdateBy)
    fmt.Printf("UpdateDate: %s\n", order.UpdateDate.Format("2006-01-02 15:04:05"))
}


func PrintReturnOrderLineDetails(line *res.ReturnOrderLine) {
    fmt.Printf("OrderNo: %s\n", line.OrderNo)
    fmt.Printf("TrackingNo: %s\n", line.TrackingNo)
    fmt.Printf("SKU: %s\n", line.SKU)
    fmt.Printf("ReturnQTY: %d\n", line.ReturnQTY)
    fmt.Printf("QTY: %d\n", line.QTY)
    fmt.Printf("Price: %.2f\n", line.Price)
    // fmt.Printf("AlterSKU: %s\n", line.AlterSKU)
    fmt.Printf("CreateBy: %s\n", line.CreateBy)
    fmt.Printf("CreateDate: %s\n", line.CreateDate)
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