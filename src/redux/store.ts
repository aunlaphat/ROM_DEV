// src/redux/store.ts
import { configureStore } from '@reduxjs/toolkit';
import createSagaMiddleware from 'redux-saga';
import { persistStore, persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import rootReducer from './rootReducer';
import rootSaga from './rootSaga';
import { RootState } from './types';

// Saga middleware configuration
const sagaMiddleware = createSagaMiddleware();

// Redux Persist configuration
const persistConfig = {
  key: 'root',
  storage,
  whitelist: ['auth'],  // เฉพาะ auth ที่จะถูกเก็บใน localStorage
  blacklist: ['order', 'alert'] // อันนี้จะไม่ถูกเก็บ
};

// Create persisted reducer
const persistedReducer = persistReducer<RootState>(persistConfig, rootReducer);

// Create store
export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      thunk: false,
      serializableCheck: {
        // Ignore these action types
        ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
        // Ignore these paths in state
        ignoredPaths: ['some.path.to.ignore'],
      },
    }).concat(sagaMiddleware),
  devTools: process.env.NODE_ENV !== 'production',
});

// Create persistor
export const persistor = persistStore(store);

// Run saga
sagaMiddleware.run(rootSaga);

// Export types
export type AppStore = typeof store;
export type AppDispatch = typeof store.dispatch;

export default store;