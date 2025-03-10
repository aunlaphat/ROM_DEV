// src/redux/draftConfirm/types/state.ts
import { Order, CodeR } from './response';

// DraftConfirmState - สถานะใน Redux store
export interface DraftConfirmState {
  orders: Order[];             // รายการคำสั่งทั้งหมด
  selectedOrder: Order | null; // คำสั่งที่เลือก
  codeRList: CodeR[];          // รายการ CodeR ทั้งหมด
  loading: boolean;            // สถานะการโหลด
  error: string | null;        // ข้อผิดพลาด
}

// initialDraftConfirmState - ค่าเริ่มต้นของ state
export const initialDraftConfirmState: DraftConfirmState = {
  orders: [],
  selectedOrder: null,
  codeRList: [],
  loading: false,
  error: null
};