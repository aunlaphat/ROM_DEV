import { ReturnOrderActionTypes, ReturnOrderState } from './types';
import { logger } from '../../utils/logger';

const initialState: ReturnOrderState = {
  step: 'search',
  orderDetails: null,
  returnDetails: null,
  srDetails: null,
  status: 'idle',
  error: null,
};

export default function returnOrderReducer(state = initialState, action: any): ReturnOrderState {
  switch (action.type) {
    case ReturnOrderActionTypes.SEARCH_ORDER_REQUEST:
    case ReturnOrderActionTypes.CREATE_RETURN_REQUEST:
    case ReturnOrderActionTypes.GENERATE_SR_REQUEST:
    case ReturnOrderActionTypes.CONFIRM_RETURN_REQUEST:
      logger.log('info', `[ReturnOrder] Processing ${action.type}`);
      return { ...state, status: 'loading', error: null };

    case ReturnOrderActionTypes.SEARCH_ORDER_SUCCESS:
      logger.log('info', `[ReturnOrder] Search Success`, action.payload);
      return { ...state, orderDetails: action.payload, status: 'success', step: 'createReturn' };

    case ReturnOrderActionTypes.CREATE_RETURN_SUCCESS:
      logger.log('info', `[ReturnOrder] Create Success`, action.payload);
      return { ...state, returnDetails: action.payload, status: 'success', step: 'generateSR' };

    case ReturnOrderActionTypes.GENERATE_SR_SUCCESS:
      logger.log('info', `[ReturnOrder] Generate SR Success`, action.payload);
      return { ...state, srDetails: action.payload, status: 'success', step: 'preview' };

    case ReturnOrderActionTypes.CONFIRM_RETURN_SUCCESS:
      logger.log('info', `[ReturnOrder] Confirm Success`);
      return { ...state, status: 'success', step: 'confirm' };

    case ReturnOrderActionTypes.SEARCH_ORDER_FAILURE:
    case ReturnOrderActionTypes.CREATE_RETURN_FAILURE:
    case ReturnOrderActionTypes.GENERATE_SR_FAILURE:
    case ReturnOrderActionTypes.CONFIRM_RETURN_FAILURE:
      logger.log('error', `[ReturnOrder] Failure: ${action.type}`, action.payload);
      return { ...state, status: 'error', error: action.payload };

    case ReturnOrderActionTypes.SET_STEP:
      return { ...state, step: action.payload };

    case ReturnOrderActionTypes.RESET:
      return initialState;

    default:
      return state;
  }
}
