// src/redux/types.ts
import { OrderState } from './orders/types';
import { AuthState } from './auth/interface';
import { DraftConfirmState } from './draftConfirm/types';

// Define AlertState interface
export interface AlertState {
  open: boolean;
  loading: boolean;
  type: string;
  model: string;
  title: string;
  message: string;
}

// Define RootState interface for type-safe access to Redux store
export interface RootState {
  auth: AuthState;
  order: OrderState;
  alert: AlertState;
  draftConfirm: DraftConfirmState;
}