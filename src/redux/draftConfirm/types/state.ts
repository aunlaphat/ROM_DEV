// src/redux/draftConfirm/types/state.ts
import { Order, CodeR } from "./response";

// DraftConfirmState - สถานะใน Redux store
export interface DraftConfirmState {
  orders: Order[];
  selectedOrder: Order | null;
  codeRList: CodeR[];
  loading: boolean;
  error: string | null;
  lastSearchDates: {
    startDate: string | null;
    endDate: string | null;
  };
}

// initialDraftConfirmState - ค่าเริ่มต้นของ state
export const initialDraftConfirmState: DraftConfirmState = {
  orders: [],
  selectedOrder: null,
  codeRList: [],
  loading: false,
  error: null,
  lastSearchDates: {
    startDate: null,
    endDate: null,
  },
};
