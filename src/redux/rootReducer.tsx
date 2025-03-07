import { combineReducers } from "redux";
import auth from "./auth/reducer";
import alert from "./alert/reducer";
import order from "./orders/reducer";

const rootReducer = combineReducers({
  auth,
  alert,
  order
});

export default rootReducer;
