import { all } from "redux-saga/effects";
import authSaga from "./auth/sagas";
import orderSaga from "./orders/sagas";

export default function* rootSaga(): Generator {
  yield all([
    authSaga(),
    orderSaga()
  ]);
}
