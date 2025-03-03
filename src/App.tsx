import "./App.css";
import { ConfigProvider } from "antd";
import thTH from "antd/lib/locale/th_TH";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { Provider } from "react-redux";
import store from "./redux/store";
import { ROUTES, ROUTES_NO_AUTH } from "./resources/routes";
import Loading from "./components/loading";
import Alert from "./components/alert/alert";
import LayoutPage from "./components/Layouts";

const App = () => {
  return (
    <Provider store={store}>
      <ConfigProvider locale={thTH}>
        <Loading>
          <Alert />
          <Router>
            <Routes>
              {/* Public Routes */}
              {Object.values(ROUTES_NO_AUTH).map((route) => (
                <Route
                  key={route.KEY}
                  path={route.PATH}
                  element={<route.ELEMENT />}
                />
              ))}

              {/* Protected Routes */}
              <Route path="/*" element={<LayoutPage />}>
                <Route path="home" element={<ROUTES.ROUTE_MAIN.COMPONENT />} />
                {Object.values(ROUTES)
                  .filter(route => route.KEY !== "home") // ไม่รวม home route เพราะกำหนดแล้ว
                  .map((route) => (
                    <Route
                      key={route.KEY}
                      path={route.PATH.replace('/', '')} // ตัด / ออกเพราะอยู่ใน nested route
                      element={<route.COMPONENT />}
                    />
                  ))}
              </Route>
            </Routes>
          </Router>
        </Loading>
      </ConfigProvider>
    </Provider>
  );
};

export default App;
