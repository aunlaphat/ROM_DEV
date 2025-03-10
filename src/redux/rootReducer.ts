// src/redux/rootReducer.ts
import { combineReducers, Reducer } from 'redux';
import authReducer from './auth/reducer';
import orderReducer from './orders/reducer';
import alertReducer from './alert/reducer';

// Import or define state types
// You may need to adjust these imports based on your actual file structure
import { AuthState } from './auth/interface';
import { OrderState } from './orders/types'; // or './orders/types/state' if you're using the nested structure
import { AlertState } from './types'; // or './alert/types' if defined there

// Use type assertions to tell TypeScript these are reducer functions
const rootReducer = combineReducers({
  auth: authReducer as Reducer<AuthState>,
  order: orderReducer as Reducer<OrderState>,
  alert: alertReducer as Reducer<AlertState>
});

export default rootReducer;