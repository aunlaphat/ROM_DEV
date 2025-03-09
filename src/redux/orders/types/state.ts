// src/redux/orders/types/state.ts
import { OrderStep } from './action';
import { SearchOrderResponse, ReturnOrderResponse } from './response';

// เพิ่ม OrderData interface
export interface OrderData {
  orderNo: string;
  soNo: string;
  srNo: string | null;
  isCNCreated: boolean;
  isEdited: boolean;
  head: {
    orderNo: string;
    soNo: string;
    srNo: string | null;
    salesStatus: string;
    mkpStatus: string;
    locationTo: string;
    statusReturnID?: number;
    statusConfID?: number;
    confirmBy?: string;
    confirmDate?: string;
  };
  lines: OrderLineItem[];
}

// เพิ่ม OrderLineItem interface
export interface OrderLineItem {
  sku: string;
  itemName: string;
  qty: number;
  price: number;
  warehouse?: string;
  returnQty?: number;
}

// แก้ไข OrderState ให้มี orderData
export interface OrderState {
  // เพิ่ม orderData
  orderData: OrderData | null;
  
  searchResult: SearchOrderResponse | null;
  returnOrder: ReturnOrderResponse | null;
  currentStep: OrderStep;
  loading: boolean;
  error: string | null;
  srCreated: boolean;
  isEdited: boolean;
}

// แก้ไข initialOrderState ให้มี orderData
export const initialOrderState: OrderState = {
  orderData: null,
  searchResult: null,
  returnOrder: null,
  currentStep: 'search',
  loading: false,
  error: null,
  srCreated: false,
  isEdited: false
};