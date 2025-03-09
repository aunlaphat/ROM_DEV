// src/screens/Orders/Marketplace/types.ts
import { OrderStep } from '../../../redux/orders/types';
import { OrderData } from '../../../redux/orders/types/state';

/**
 * เพิ่ม ReturnOrderState interface ที่ถูกใช้ในหลาย component
 * แต่ไม่ได้ถูกนิยาม
 */
export interface ReturnOrderState {
  orderData: OrderData | null;
  currentStep: OrderStep;
  searchResult: any | null;
  returnOrder: any | null;
  loading: boolean;
  error: string | null;
  srCreated: boolean;
  isEdited: boolean;
}