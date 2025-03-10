// src/redux/draftConfirm/reducer.ts
import { DraftConfirmActionTypes, DraftConfirmState, initialDraftConfirmState } from './types';
import { logger } from '../../utils/logger';

const draftConfirmReducer = (state: DraftConfirmState = initialDraftConfirmState, action: any): DraftConfirmState => {
  switch (action.type) {
    // Request actions - เริ่มการโหลด
    case DraftConfirmActionTypes.FETCH_ORDERS_REQUEST:
    case DraftConfirmActionTypes.FETCH_ORDER_DETAILS_REQUEST:
    case DraftConfirmActionTypes.FETCH_CODE_R_REQUEST:
    case DraftConfirmActionTypes.ADD_ITEM_REQUEST:
    case DraftConfirmActionTypes.REMOVE_ITEM_REQUEST:
    case DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_REQUEST:
    case DraftConfirmActionTypes.CANCEL_ORDER_REQUEST:
      logger.redux.action(action.type);
      return {
        ...state,
        loading: true,
        error: null
      };

    // Fetch Orders Success
    case DraftConfirmActionTypes.FETCH_ORDERS_SUCCESS:
      logger.redux.action(action.type, { orderCount: action.payload.length });
      return {
        ...state,
        orders: action.payload,
        loading: false,
        error: null
      };

    // Fetch Order Details Success
    case DraftConfirmActionTypes.FETCH_ORDER_DETAILS_SUCCESS:
      logger.redux.action(action.type, { orderNo: action.payload.orderNo });
      return {
        ...state,
        selectedOrder: action.payload,
        loading: false,
        error: null
      };

    // Fetch CodeR List Success
    case DraftConfirmActionTypes.FETCH_CODE_R_SUCCESS:
      logger.redux.action(action.type, { codeRCount: action.payload.length });
      return {
        ...state,
        codeRList: action.payload,
        loading: false,
        error: null
      };

    // Add Item Success
    case DraftConfirmActionTypes.ADD_ITEM_SUCCESS:
      logger.redux.action(action.type, { items: action.payload });
      return {
        ...state,
        loading: false,
        error: null
      };

    // Remove Item Success
    case DraftConfirmActionTypes.REMOVE_ITEM_SUCCESS:
      logger.redux.action(action.type);
      return {
        ...state,
        loading: false,
        error: null
      };

    // Confirm Draft Order Success
    case DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_SUCCESS:
      logger.redux.action(action.type, { 
        orderNo: action.payload.orderNo,
        statusReturnID: action.payload.statusReturnID,
        statusConfID: action.payload.statusConfID
      });
      return {
        ...state,
        loading: false,
        error: null
      };

    // Cancel Order Success
    case DraftConfirmActionTypes.CANCEL_ORDER_SUCCESS:
      logger.redux.action(action.type);
      return {
        ...state,
        loading: false,
        error: null
      };

    // Clear Selected Order
    case DraftConfirmActionTypes.CLEAR_SELECTED_ORDER:
      return {
        ...state,
        selectedOrder: null
      };

    // Error actions
    case DraftConfirmActionTypes.FETCH_ORDERS_FAILURE:
    case DraftConfirmActionTypes.FETCH_ORDER_DETAILS_FAILURE:
    case DraftConfirmActionTypes.FETCH_CODE_R_FAILURE:
    case DraftConfirmActionTypes.ADD_ITEM_FAILURE:
    case DraftConfirmActionTypes.REMOVE_ITEM_FAILURE:
    case DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_FAILURE:
    case DraftConfirmActionTypes.CANCEL_ORDER_FAILURE:
      logger.redux.action(action.type, { error: action.payload });
      return {
        ...state,
        loading: false,
        error: action.payload
      };

    default:
      return state;
  }
};

export default draftConfirmReducer;