// src/redux/rootSaga.ts
import { all, fork } from 'redux-saga/effects';
import authSaga from './auth/saga';
import orderSaga from './orders/saga';

// Root saga that combines all other sagas
export default function* rootSaga() {
  yield all([
    fork(authSaga),
    fork(orderSaga)
  ]);
}