// src/redux/rootReducer.ts
import { combineReducers, Reducer } from 'redux';
import authReducer from './auth/reducer';
import orderReducer from './orders/reducer';
import alertReducer from './alert/reducer';
import draftConfirmReducer from './draftConfirm/reducer';
import userReducer from './users/reducer';

// Import state types
import { AuthState } from './auth/interface';
import { OrderState } from './orders/types';
import { AlertState } from './types';
import { DraftConfirmState } from './draftConfirm/types'; 
import { UserState } from './users/types';

const rootReducer = combineReducers({
  auth: authReducer as Reducer<AuthState>,
  order: orderReducer as Reducer<OrderState>,
  alert: alertReducer as Reducer<AlertState>,
  draftConfirm: draftConfirmReducer as Reducer<DraftConfirmState>,
  user: userReducer as Reducer<UserState>
});

export default rootReducer;