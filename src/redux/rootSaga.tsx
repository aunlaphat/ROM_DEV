import { all } from "redux-saga/effects";
import authenSaga from "./auth/saga";

export default function* rootSaga(): Generator {
  yield all([authenSaga()]);
}
