import { all } from "redux-saga/effects";
import authenSaga from "./authen/saga";

export default function* rootSaga() {
  yield all([authenSaga()]);
}
