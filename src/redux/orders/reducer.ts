// src/redux/orders/reducer.ts
import { OrderActionTypes, OrderState, OrderStep, SearchOrderItem } from './types';
import { logger } from '../../utils/logger';

// ขั้นตอนที่ 7: กำหนด Initial State และ Reducer

// Initial state
const initialState: OrderState = {
  orderData: null,
  searchResult: null,
  returnOrder: null,
  currentStep: 'search',
  loading: false,
  error: null,
  srCreated: false,
  isEdited: false
};

// Reducer
const orderReducer = (state: OrderState = initialState, action: any): OrderState => {
  switch (action.type) {
    // Request actions - เริ่มการโหลด
    case OrderActionTypes.SEARCH_ORDER_REQUEST:
    case OrderActionTypes.CREATE_RETURN_ORDER_REQUEST:
    case OrderActionTypes.GENERATE_SR_REQUEST:
    case OrderActionTypes.UPDATE_SR_REQUEST:
    case OrderActionTypes.UPDATE_STATUS_REQUEST:
    case OrderActionTypes.CANCEL_ORDER_REQUEST:
    case OrderActionTypes.MARK_EDITED_REQUEST:
      logger.redux.action(action.type);
      return {
        ...state,
        loading: true,
        error: null
      };

    // Search order success
    case OrderActionTypes.SEARCH_ORDER_SUCCESS:
  logger.redux.action(action.type, {
    orderNo: action.payload.orderNo,
    soNo: action.payload.soNo,
    itemCount: action.payload.items.length
  });

  // สร้าง orderData จาก searchResult
  const orderData = {
    orderNo: action.payload.orderNo,
    soNo: action.payload.soNo,
    srNo: null,
    isCNCreated: false,
    isEdited: false,
    head: {
      orderNo: action.payload.orderNo,
      soNo: action.payload.soNo,
      srNo: null,
      salesStatus: action.payload.salesStatus,
      mkpStatus: action.payload.statusMKP,
      locationTo: 'Return',
      // จากการเพิ่มฟิลด์ตาม suggestion ก่อนหน้า
      statusReturnID: undefined,
      statusConfID: undefined, 
      confirmBy: undefined,
      confirmDate: undefined
    },
    lines: action.payload.items.map((item: SearchOrderItem) => ({
      ...item,
      price: Math.abs(item.price),
      returnQty: 0
    }))
  };

  return {
    ...state,
    searchResult: action.payload,
    orderData: orderData,
    loading: false,
    error: null,
    currentStep: 'create'
  };

    // Create return order success
    case OrderActionTypes.CREATE_RETURN_ORDER_SUCCESS:
      logger.redux.action(action.type, {
        orderNo: action.payload.orderNo,
        soNo: action.payload.soNo,
        itemCount: action.payload.items.length
      });

      return {
        ...state,
        returnOrder: action.payload,
        loading: false,
        error: null,
        currentStep: 'sr' // เปลี่ยนขั้นตอนเป็น sr หลังสร้างสำเร็จ
      };

    // Generate SR success
    case OrderActionTypes.GENERATE_SR_SUCCESS:
      logger.redux.action(action.type, { srNo: action.payload });

      return {
        ...state,
        loading: false,
        error: null,
        // เพิ่มการอัพเดทค่า srNo ใน orderData
        orderData: state.orderData ? {
          ...state.orderData,
          srNo: action.payload,
          head: {
            ...state.orderData.head,
            srNo: action.payload
          }
        } : null
      };

    // แก้ไขกรณี UPDATE_SR_SUCCESS
    case OrderActionTypes.UPDATE_SR_SUCCESS:
      logger.redux.action(action.type, {
        orderNo: action.payload.orderNo,
        srNo: action.payload.srNo
      });

      return {
        ...state,
        returnOrder: state.returnOrder ? {
          ...state.returnOrder,
          srNo: action.payload.srNo
        } : null,
        // เพิ่มการอัพเดทค่า srNo ใน orderData
        orderData: state.orderData ? {
          ...state.orderData,
          srNo: action.payload.srNo,
          head: {
            ...state.orderData.head,
            srNo: action.payload.srNo
          }
        } : null,
        srCreated: true,
        loading: false,
        error: null,
        currentStep: 'preview'
      };

    // แก้ไขกรณี UPDATE_STATUS_SUCCESS
    case OrderActionTypes.UPDATE_STATUS_SUCCESS:
      logger.redux.action(action.type, {
        orderNo: action.payload.orderNo,
        statusReturnID: action.payload.statusReturnID,
        statusConfID: action.payload.statusConfID
      });

      return {
        ...state,
        returnOrder: state.returnOrder ? {
          ...state.returnOrder,
          statusReturnId: action.payload.statusReturnID,
          statusConfId: action.payload.statusConfID,
          confirmBy: action.payload.confirmBy,
          confirmDate: action.payload.confirmDate
        } : null,
        // เพิ่มการอัพเดทค่าสถานะใน orderData
        orderData: state.orderData ? {
          ...state.orderData,
          head: {
            ...state.orderData.head,
            statusReturnID: action.payload.statusReturnID,
            statusConfID: action.payload.statusConfID,
            confirmBy: action.payload.confirmBy,
            confirmDate: action.payload.confirmDate
          }
        } : null,
        loading: false,
        error: null,
        currentStep: 'confirm'
      };

    // Cancel order success
    case OrderActionTypes.CANCEL_ORDER_SUCCESS:
      logger.redux.action(action.type, {
        refID: action.payload.refID,
        sourceTable: action.payload.sourceTable
      });

      return {
        ...state,
        loading: false,
        error: null,
        // ไม่เปลี่ยนขั้นตอนเพราะอาจยกเลิกจากขั้นตอนใดก็ได้
      };

    // Mark as edited success
    case OrderActionTypes.MARK_EDITED_SUCCESS:
      logger.redux.action(action.type);

      return {
        ...state,
        isEdited: true,
        returnOrder: state.returnOrder ? {
          ...state.returnOrder,
          isEdited: true
        } : null,
        loading: false,
        error: null
      };

    // Error actions
    case OrderActionTypes.SEARCH_ORDER_FAILURE:
    case OrderActionTypes.CREATE_RETURN_ORDER_FAILURE:
    case OrderActionTypes.GENERATE_SR_FAILURE:
    case OrderActionTypes.UPDATE_SR_FAILURE:
    case OrderActionTypes.UPDATE_STATUS_FAILURE:
    case OrderActionTypes.CANCEL_ORDER_FAILURE:
    case OrderActionTypes.MARK_EDITED_FAILURE:
      logger.redux.action(action.type, { error: action.payload });

      return {
        ...state,
        loading: false,
        error: action.payload
      };

    // Set current step
    case OrderActionTypes.SET_CURRENT_STEP:
      return {
        ...state,
        currentStep: action.payload as OrderStep
      };

    // Reset order state
    case OrderActionTypes.RESET_ORDER:
      return initialState;

    // Clear error
    case OrderActionTypes.CLEAR_ERROR:
      return {
        ...state,
        error: null
      };

    default:
      return state;
  }
};

export default orderReducer;