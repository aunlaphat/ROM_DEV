import "./App.css";
import { ConfigProvider } from "antd";
import thTH from "antd/lib/locale/th_TH";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import LayoutPage from "./components/Layouts";
import { Provider } from "react-redux";
import store from "./redux/store";
import { ROUTES_PATH, ROUTE_NOT_AUTHEN } from "./resources/routes-name";
import Loading from "./components/loading";
import Alert from "./components/alert/alert";

const App = () => {
  console.log(Object.values(ROUTES_PATH).map((item) => {console.log(item)}))
  return (
    <Provider store={store}>
      <ConfigProvider locale={thTH}>
        <Loading>
          <Alert />
          <Router>
            <Routes>
              {/* {Object.values(ROUTE_NOT_AUTHEN).map((item) => (
                <Route
                  path={item.PATH}
                  key={item.KEY}
                  element={item.ELEMENT()}
                />
              ))} */}
              <Route element={<LayoutPage />}>
                {Object.values(ROUTES_PATH).map((item) => (
                  <Route
                    path={item.PATH}
                    key={item.KEY}
                    element={item.ELEMENT()}
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
