import { combineReducers } from "redux";
import authen from "./authen/reducer";
import alert from "./alert/reducer";

const rootReducer = combineReducers({
  authen: authen,
  alert: alert,
});

export default rootReducer;
