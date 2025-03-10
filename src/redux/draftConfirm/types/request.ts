// src/redux/draftConfirm/types/request.ts

// Request Types (สอดคล้องกับ GoLang request types)

// FetchOrders - เรียกดูรายการคำสั่งตาม StatusConfID และช่วงวันที่
export interface FetchOrdersRequest {
    statusConfID: number;
    startDate?: string;
    endDate?: string;
  }
  
  // FetchOrderDetails - เรียกดูรายละเอียดคำสั่งพร้อมรายการสินค้า
  export interface FetchOrderDetailsRequest {
    orderNo: string;
    statusConfID: number;
  }
  
  // AddItem - เพิ่มรายการสินค้า
  export interface AddItemRequest {
    orderNo: string;
    sku: string;
    itemName: string;
    qty: number;
    returnQty?: number;
    price: number;
  }
  
  // RemoveItem - ลบรายการสินค้า
  export interface RemoveItemRequest {
    orderNo: string;
    sku: string;
  }
  
  // ConfirmDraftOrder - ยืนยันคำสั่ง Draft เป็น Confirm
  export interface ConfirmDraftOrderRequest {
    orderNo: string;
  }
  
  // CancelOrder - ยกเลิกคำสั่ง
  export interface CancelOrderRequest {
    orderNo: string;
    cancelReason: string;
  }