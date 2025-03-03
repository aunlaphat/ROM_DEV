import { BeforeReturnOrderResponse, SearchOrderResponse, ReturnOrderState, SearchOrderItem } from './api';
import { ReturnOrderActionTypes } from './types';

const initialState: ReturnOrderState = {
  orderData: null,
  searchResult: null,
  returnOrder: null,
  loading: false,
  error: null,
  currentStep: 'search',
  isEdited: false,
  orderLines: [],
  srCreated: false
};

export default function returnOrderReducer(state = initialState, action: any): ReturnOrderState {
  switch (action.type) {
    // Request cases
    case ReturnOrderActionTypes.RETURN_ORDER_SEARCH_REQ:
    case ReturnOrderActionTypes.RETURN_ORDER_CREATE_REQ:
    case ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_REQ:
    case ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_REQ:
    case ReturnOrderActionTypes.RETURN_ORDER_CANCEL_REQ:
      return {
        ...state,
        loading: true,
        error: null
      };

    // Success cases  
    case ReturnOrderActionTypes.RETURN_ORDER_SEARCH_SUCCESS:
      return {
        ...state,
        loading: false,
        orderData: {
          head: {
            orderNo: action.payload.orderNo,  // แก้จาก soNo เป็น orderNo
            soNo: action.payload.soNo,
            srNo: null, // เริ่มต้นเป็น null เพราะยังไม่มีการสร้าง SR
            salesStatus: action.payload.salesStatus,
            mkpStatus: action.payload.statusMKP,
            locationTo: 'Return' // default value
          },
          lines: action.payload.items.map((item: SearchOrderItem) => ({
            ...item,
            price: Math.abs(item.price) // แปลงให้เป็นค่าบวกเสมอ
          }))
        },
        orderLines: action.payload.items,
        currentStep: 'create',
        error: null
      };

    case ReturnOrderActionTypes.RETURN_ORDER_CREATE_SUCCESS:
      return {
        ...state,
        loading: false,
        returnOrder: action.payload,
        currentStep: 'confirm'
      };

    case ReturnOrderActionTypes.RETURN_ORDER_MARK_EDITED_SUCCESS:
      return {
        ...state,
        loading: false,
        isEdited: true
      };

    case ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_SUCCESS:
      return {
        ...state,
        loading: false,
        srCreated: true,
        orderData: state.orderData ? {
          ...state.orderData,
          head: {
            ...state.orderData.head,
            srNo: action.payload.data.srNo,
          }
        } : null
      };

    // Failure cases
    case ReturnOrderActionTypes.RETURN_ORDER_SEARCH_FAIL:
    case ReturnOrderActionTypes.RETURN_ORDER_CREATE_FAIL:
    case ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_FAIL:
    case ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_FAIL:
    case ReturnOrderActionTypes.RETURN_ORDER_CANCEL_FAIL:
      return {
        ...state,
        loading: false,
        error: action.payload
      };

    case ReturnOrderActionTypes.RETURN_ORDER_RESET:
      return initialState; // รีเซ็ตทุก state กลับไปเป็นค่าเริ่มต้น

    default:
      return state;
  }
}
