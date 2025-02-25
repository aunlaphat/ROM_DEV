import { configureStore } from "@reduxjs/toolkit";
import createSagaMiddleware from "redux-saga";
import rootReducer from "./rootReducer";
import rootSaga from "./rootSaga";

const sagaMiddleware = createSagaMiddleware();

// ✅ สร้าง Redux Store พร้อม Redux-Saga Middleware
export const store = configureStore({
  reducer: rootReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({ thunk: false }).concat(sagaMiddleware), // ✅ ปิด Redux-Thunk ที่ไม่ใช้
  devTools: process.env.NODE_ENV !== "production", // ✅ เปิด DevTools เฉพาะใน Development Mode
});

// ✅ รัน Redux-Saga
sagaMiddleware.run(rootSaga);

// ✅ กำหนด Type สำหรับ Redux Store
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
