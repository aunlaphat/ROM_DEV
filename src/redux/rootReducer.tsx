import { combineReducers } from "redux";
import auth from "./auth/reducer";
import alert from "./alert/reducer";
import returnOrder from "./orders/reducer";

const rootReducer = combineReducers({
  auth,
  alert,
  returnOrder
});

export default rootReducer;
