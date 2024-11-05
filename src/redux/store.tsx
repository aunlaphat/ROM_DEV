import { createStore, applyMiddleware } from "redux";
import rootReducer from "./rootReducer";
import createSagaMiddleware from "redux-saga";
import rootSaga from "./rootSaga";
import { composeWithDevTools } from "redux-devtools-extension";
import { configureStore } from "@reduxjs/toolkit";

/**
 * Redux-saga --> https://redux-saga.js.org/
 * Tutorial with Example --> https://www.blog.duomly.com/implement-redux-saga-with-reactjs-and-redux/
 */

/** Redux dev tool -->  https://www.npmjs.com/package/@redux-devtools/extension
 * ถ้าจะใช้ต้องลง lib และ extension ใน chrome ด้วย [Redux DevTools]
 */
// const composeEnhancers = composeWithDevTools({
//   // Specify here name, actionsDenylist, actionsCreators and other options
//   trace: true,
// });

// const initialState = {
//   authen: {
//     loading: false,
//     users: [],
//   },
// };

const sagaMiddleware = createSagaMiddleware();

const store = configureStore({
  reducer: rootReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(sagaMiddleware),
});
sagaMiddleware.run(rootSaga);

export default store;
