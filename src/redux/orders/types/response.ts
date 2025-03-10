// src/redux/orders/types/responses.ts

// Response Types (สอดคล้องกับ GoLang response types)

// SearchOrderItem - รายการสินค้าจากการค้นหา
export interface SearchOrderItem {
  sku: string;
  itemName: string;
  qty: number;
  price: number;
}

// SearchOrderResponse - ผลลัพธ์การค้นหา
export interface SearchOrderResponse {
  soNo: string;
  orderNo: string;
  statusMKP: string;
  salesStatus: string;
  createDate: string;
  items: SearchOrderItem[];
}

// ReturnOrderItem - รายการสินค้าที่คืน
export interface ReturnOrderItem {
  orderNo: string;
  sku: string;
  itemName: string;
  qty: number;
  returnQty: number;
  price: number;
  createBy: string;
  createDate: string;
  trackingNo?: string;
  alterSKU?: string;
}

// ReturnOrderResponse - ผลลัพธ์การสร้าง return order
export interface ReturnOrderResponse {
  orderNo: string;
  soNo: string;
  srNo: string | null;
  channelId: number;
  reason: string;
  customerId: string;
  trackingNo: string;
  logistic: string;
  warehouseId: number;
  soStatus: string | null;
  mkpStatus: string | null;
  returnDate: string | null;
  statusReturnId: number | null;
  statusConfId: number | null;
  confirmBy: string | null;
  confirmDate: string | null;
  createBy: string;
  createDate: string;
  updateBy: string | null;
  updateDate: string | null;
  cancelId: number | null;
  isCNCreated: boolean;
  isEdited: boolean;
  items: ReturnOrderItem[];
}

// UpdateSrResponse - ผลลัพธ์การอัพเดท SR
export interface UpdateSrResponse {
  orderNo: string;
  srNo: string;
  statusReturnID?: number;
  statusConfID?: number;
  updateBy: string;
  updateDate: string;
}

// UpdateStatusResponse - ผลลัพธ์การอัพเดทสถานะ
export interface UpdateStatusResponse {
  orderNo: string;
  statusReturnID: number;
  statusConfID: number;
  confirmBy: string;
  confirmDate: string;
}

// CancelOrderResponse - ผลลัพธ์การยกเลิก order
export interface CancelOrderResponse {
  refID: string;
  sourceTable: string;
  cancelReason: string;
  cancelBy: string;
  cancelDate: string;
}