import "./App.css";
import { ConfigProvider } from "antd";
import thTH from "antd/lib/locale/th_TH";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import LayoutPage from "./components/Layouts";
import { Provider } from "react-redux";
import store from "./redux/store";
import { ROUTES_PATH, ROUTE_NOT_AUTHEN } from "./resources/routes";
import Loading from "./components/loading";
import Alert from "./components/alert/alert";

const App = () => {
    console.log("🔍 Debugging Available Routes...");
    console.log("✅ ROUTES_PATH:", ROUTES_PATH);
    console.log("✅ ROUTE_NOT_AUTHEN:", ROUTE_NOT_AUTHEN);
  
    return (
      <Provider store={store}>
        <ConfigProvider locale={thTH}>
          <Loading>
            <Alert />
            <Router>
              <Routes>
                {/* 🔓 Debug Log: ตรวจสอบ Routes ที่ไม่ต้องล็อกอิน */}
                {Object.values(ROUTE_NOT_AUTHEN)?.map((item) => {
                  return (
                    <Route
                      path={item?.PATH}
                      key={item?.KEY}
                      element={item?.ELEMENT ? item.ELEMENT() : <div>❌ Error: No Component</div>}
                    />
                  );
                })}
  
                {/* 🔐 Debug Log: ตรวจสอบ Routes ที่ต้องล็อกอิน */}
                <Route path="/*" element={<LayoutPage />}>
                  {Object.values(ROUTES_PATH)?.map((item) => {
                    return (
                      <Route
                        path={item?.PATH}
                        key={item?.KEY}
                        element={item?.ELEMENT ? item.ELEMENT() : <div>❌ Error: No Component</div>}
                      />
                    );
                  })}
                </Route>
  
                {/* ❌ ถ้าไม่มีเส้นทางนี้ให้แสดง 404 */}
                <Route path="*" element={<div>❌ 404 Not Found</div>} />
              </Routes>
            </Router>
          </Loading>
        </ConfigProvider>
      </Provider>
    );
  };
  
export default App;