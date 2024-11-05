import * as type from "./types";

const initialState = {
  open: false,
  loading: false,
  type: "success",
  model: "notification",
  title: "",
  message: "",
};

export default function alertReducer(state = initialState, action: any) {
  switch (action.type) {
    case type.ALERT_REQ:
      return {
        ...state,
        ...action.payload,
        open: true,
      };
    case type.ALERT_CLOSE:
      return {
        ...state,
        type: "",
        model: "notification",
        title: "",
        message: "",
        open: false,
      };
    case type.ALERT_LOADING:
      return {
        ...state,
        loading: true,
        open: false,
      };
    case type.ALERT_LOADING_CLOSE:
      return {
        ...state,
        loading: false,
      };
    default:
      return state;
  }
}
