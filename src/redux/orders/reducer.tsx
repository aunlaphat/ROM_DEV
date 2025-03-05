import { message } from 'antd';
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
        currentStep: 'create', // เปลี่ยนเป็น create หลังค้นหาสำเร็จ
        error: null
      };

    case ReturnOrderActionTypes.RETURN_ORDER_CREATE_SUCCESS:
      return {
        ...state,
        loading: false,
        returnOrder: action.payload,
        currentStep: 'sr' // เปลี่ยนเป็น sr หลังสร้าง return order สำเร็จ
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
        currentStep: 'preview', // เปลี่ยนเป็น preview หลังสร้าง SR สำเร็จ
        orderData: state.orderData ? {
          ...state.orderData,
          head: {
            ...state.orderData.head,
            srNo: action.payload.srNo,
          },
          // ใช้ข้อมูล items จาก returnOrder เสมอหลังจากสร้าง SR
          lines: state.returnOrder?.items.map(item => ({
            ...item,
            qty: item.qty,
            returnQty : item.returnQty,
            price: Math.abs(item.price)
          })) || []
        } : null,
        // อัพเดท returnOrder ถ้ามี
        returnOrder: state.returnOrder ? {
          ...state.returnOrder,
          srNo: action.payload.srNo
        } : null
      };

    case ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_SUCCESS:
      // ปิด loading message
      message.success({
        content: 'อัพเดตสถานะสำเร็จ',
        key: 'confirmStatus',
        duration: 2
      });
      
      return {
        ...state,
        loading: false,
        currentStep: 'confirm', // เพิ่มการเปลี่ยน step เป็น confirm
        orderData: state.orderData ? {
          ...state.orderData,
          head: {
            ...state.orderData.head,
            statusReturnID: action.payload.statusReturnID,
            statusConfID: action.payload.statusConfID,
            confirmBy: action.payload.confirmBy,
            confirmDate: action.payload.confirmDate
          }
        } : null
      };

    case ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_FAIL:
      // ปิด loading message พร้อมแสดง error
      message.error({
        content: 'อัพเดตสถานะไม่สำเร็จ',
        key: 'confirmStatus',
        duration: 2
      });
      return {
        ...state,
        loading: false,
        error: action.payload
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

    // เพิ่ม case สำหรับ SET_STEP
    case ReturnOrderActionTypes.RETURN_ORDER_SET_STEP:
      return {
        ...state,
        currentStep: action.payload
      };

    default:
      return state;
  }
}
