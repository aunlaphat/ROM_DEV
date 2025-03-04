import { all } from "redux-saga/effects";
import authenSaga from "./auth/saga";
import returnOrderSaga from "./orders/saga";

export default function* rootSaga(): Generator {
  yield all([
    authenSaga(),
    returnOrderSaga()
  ]);
}
