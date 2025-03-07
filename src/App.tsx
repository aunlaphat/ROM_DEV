import "./App.css";
import { ConfigProvider } from "antd";
import thTH from "antd/lib/locale/th_TH";
import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";
import { Provider } from "react-redux";
import store from "./redux/store";
import { ROUTES, ROUTES_NO_AUTH } from "./resources/routes";
import Loading from "./components/loading";
import Alert from "./components/alert/alert";
import LayoutPage from "./components/Layouts";
import ProtectedRoute from "./resources/ProtectedRoute";
import { AuthProvider } from "./hooks/auth";
import { NotFound } from "./screens/NotFound";
import { RoleID } from "./constants/roles";
import { Login } from "./screens/auth";
import { lazy, Suspense } from "react";

// Lazy-load เฉพาะหน้าที่มีขนาดใหญ่เพื่อเพิ่มประสิทธิภาพ
const CreateReturnOrderMKP = lazy(() => import('./screens/Orders/Marketplace'));
// const DraftAndConfirm = lazy(() => import('./screens/Draft&Confirm'));
const ManageUser = lazy(() => import('./screens/ManageUser'));
const Home = lazy(() => import('./screens/Home'));

/**
 * โครงสร้างของ routes ที่มีความซับซ้อน โดยกำหนดสิทธิ์แยกตามหน้า
 */
const ROLE_BASED_ROUTES = [
  {
    path: "home",
    element: Home,
    roles: [RoleID.ADMIN, RoleID.TRADE_CONSIGN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.VIEWER],
  },
  {
    path: "orders/marketplace",
    element: CreateReturnOrderMKP,
    roles: [RoleID.ADMIN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.TRADE_CONSIGN, RoleID.VIEWER],
  },
  // {
  //   path: "draft-confirm",
  //   element: DraftAndConfirm,
  //   roles: [RoleID.ADMIN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.TRADE_CONSIGN, RoleID.VIEWER],
  // },
  {
    path: "manage-user",
    element: ManageUser,
    roles: [RoleID.ADMIN], // เฉพาะ admin เท่านั้น
  },
  // เพิ่ม routes อื่นๆ ตามต้องการ
];

const App = () => {
  return (
    <Provider store={store}>
      <ConfigProvider locale={thTH}>
        <Loading>
          <Alert />
          <AuthProvider>
            <Router>
              <Suspense fallback={<Loading />}>
                <Routes>
                  {/* Public Routes */}
                  <Route path="/login" element={<Login />} />
                  
                  {/* Redirect จากหน้าหลักไปที่ /home */}
                  <Route path="/" element={<Navigate to="/home" replace />} />
                  
                  {/* Protected Routes ใช้ Layout หลัก */}
                  <Route path="/" element={<LayoutPage />}>
                    {/* สร้าง routes จาก array ที่มีการกำหนด roles */}
                    {ROLE_BASED_ROUTES.map((route) => (
                      <Route
                        key={route.path}
                        path={route.path}
                        element={
                          <ProtectedRoute allowedRoles={route.roles}>
                            <route.element />
                          </ProtectedRoute>
                        }
                      />
                    ))}
                    
                    {/* หน้า 404 */}
                    <Route path="*" element={<NotFound />} />
                  </Route>
                </Routes>
              </Suspense>
            </Router>
          </AuthProvider>
        </Loading>
      </ConfigProvider>
    </Provider>
  );
};

export default App;