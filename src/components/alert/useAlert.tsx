import {
  alert,
  loading,
  loadingClose,
  removeAlert,
} from "../../redux/alert/action";
import store from "../../redux/store";
import { ALERT_REQ, ALERT_CLOSE } from "../../redux/alert/types";

export const openLoading = () => {
  store.dispatch(loading());
};

export const closeLoading = () => {
  store.dispatch(loadingClose());
};

export const openAlert = (obj: any) => {
  store.dispatch(alert({ type: ALERT_REQ, payload: obj }));
  store.dispatch(loadingClose());
};

export const closeAlert = () => {
  store.dispatch(removeAlert({ type: ALERT_CLOSE }));
};
