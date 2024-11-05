import * as type from "./types";

export function alert(action: any) {
  return {
    type: type.ALERT_REQ,
    payload: action.payload,
  };
}

export function removeAlert(action: any) {
  return {
    type: type.ALERT_CLOSE,
    payload: action.payload,
  };
}

export function loading() {
  return {
    type: type.ALERT_LOADING,
  };
}

export function loadingClose() {
  return {
    type: type.ALERT_LOADING_CLOSE,
  };
}
