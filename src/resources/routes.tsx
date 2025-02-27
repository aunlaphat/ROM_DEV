import Home from "../screens/home";
import { Login } from "../screens/auth";
import { NotFound } from "../screens/notFound";
import { Report } from "../screens/report"; // Ensure that the file '../screens/report.tsx' exists and is correctly named
import { ManageUser } from "../screens/manageUser";

export const ROUTE_LOGIN = process.env.REACT_APP_FRONTEND_URL + "/";

type RouteType = {
  KEY: string;
  PATH: string;
  LABEL: string;
  COMPONENT: React.ComponentType;
  ELEMENT: () => JSX.Element;
};

export type RoutesType = {
  ROUTE_MAIN: RouteType;
  ROUTE_NOTFOUND: RouteType;
  // ROUTE_IMPORTORDER: RouteType;
  // ROUTE_RETURNORDER: RouteType;
  // ROUTE_CREATERETURN: RouteType;
  // ROUTE_SR: RouteType;
  // ROUTE_IJ: RouteType;
  // ROUTE_CREATETRADERETURN: RouteType;
  // ROUTE_CONFIRMRETURNTRADE: RouteType;
  // ROUTE_SALERETURN: RouteType;
  // ROUTE_OTHERRETURN: RouteType;
  // ROUTE_CREATEBLINDRETURN: RouteType;
  // ROUTE_TAKEPICTURE: RouteType;
  ROUTE_REPORT: RouteType;
  // ROUTE_DRAFTANDCONFIRM: RouteType;
  ROUTE_MANAGEUSER: RouteType;
};

export const ROUTES: RoutesType = {
  ROUTE_MAIN: {
    KEY: "home",
    PATH: "/home",
    LABEL: "หน้าแรก",
    COMPONENT: Home,
    ELEMENT: () => <Home />,
  },
  ROUTE_NOTFOUND: {
    KEY: "notFound",
    PATH: "*",
    LABEL: "ไม่พบหน้านี้",
    COMPONENT: NotFound,
    ELEMENT: () => <NotFound />,
  },
  // ROUTE_IMPORTORDER: {
  //   KEY: "importOrder",
  //   PATH: "/import-order",
  //   LABEL: "นำเข้าข้อมูลการคืนสินค้า",
  //   COMPONENT: ImportOrder,
  //   ELEMENT: () => <ImportOrder />,
  // },
  // ROUTE_RETURNORDER: {
  //   KEY: "returnOrder",
  //   PATH: "/return-order",
  //   LABEL: "การคืนสินค้า",
  //   COMPONENT: ReturnOrder,
  //   ELEMENT: () => <ReturnOrder />,
  // },
  // ROUTE_CREATERETURN: {
  //   KEY: "createReturn",
  //   PATH: "/create-return",
  //   LABEL: "สร้างรายการคืนสินค้า",
  //   COMPONENT: CreateReturn,
  //   ELEMENT: () => <CreateReturn />,
  // },
  // ROUTE_SR: {
  //   KEY: "sr",
  //   PATH: "/sr",
  //   LABEL: "SR",
  //   COMPONENT: SRPage,
  //   ELEMENT: () => <SRPage />,
  // },
  // ROUTE_IJ: {
  //   KEY: "ij",
  //   PATH: "/ij",
  //   LABEL: "IJ",
  //   COMPONENT: IJPage,
  //   ELEMENT: () => <IJPage />,
  // },
  // ROUTE_CREATETRADERETURN: {
  //   KEY: "createTradeReturn",
  //   PATH: "/create-trade-return",
  //   LABEL: "สร้างรายการคืนสินค้าสำหรับฝ่ายค้าขาย",
  //   COMPONENT: CreateTradeReturn,
  //   ELEMENT: () => <CreateTradeReturn />,
  // },
  // ROUTE_CONFIRMRETURNTRADE: {
  //   KEY: "confirmReturnTrade",
  //   PATH: "/confirm-return-trade",
  //   LABEL: "ยืนยันการคืนสินค้า",
  //   COMPONENT: ConfirmReturnTrade,
  //   ELEMENT: () => <ConfirmReturnTrade />,
  // },
  // ROUTE_OTHERRETURN: {
  //   KEY: "otherReturn",
  //   PATH: "/other-return",
  //   LABEL: "การคืนสินค้าอื่นๆ",
  //   COMPONENT: OtherReturn,
  //   ELEMENT: () => <OtherReturn />,
  // },
  // ROUTE_SALERETURN: {
  //   KEY: "saleReturn",
  //   PATH: "/sale-return",
  //   LABEL: "การคืนสินค้าฝ่ายขาย",
  //   COMPONENT: SaleReturn,
  //   ELEMENT: () => <SaleReturn />,
  // },
  // ROUTE_CREATEBLINDRETURN: {
  //   KEY: "createBlindReturn",
  //   PATH: "/create-blind-return",
  //   LABEL: "สร้างรายการคืนสินค้าสำหรับการตรวจสอบ",
  //   COMPONENT: CreateBlindReturn,
  //   ELEMENT: () => <CreateBlindReturn />,
  // },
  // ROUTE_TAKEPICTURE: {
  //   KEY: "takePicture",
  //   PATH: "/take-picture",
  //   LABEL: "ถ่ายรูปสินค้า",
  //   COMPONENT: TakePicture,
  //   ELEMENT: () => <TakePicture />,
  // },
  ROUTE_REPORT: {
    KEY: "report",
    PATH: "/report",
    LABEL: "รายงานผลการคืนสินค้า",
    COMPONENT: Report,
    ELEMENT: () => <Report />,
  },
  // ROUTE_DRAFTANDCONFIRM: {
  //   KEY: "draftAndConfirm",
  //   PATH: "/draft-and-confirm",
  //   LABEL: "ร่างและยืนยันการคืนสินค้า",
  //   COMPONENT: DraftAndConfirm,
  //   ELEMENT: () => <DraftAndConfirm />,
  // },
  ROUTE_MANAGEUSER: {
    KEY: "manageUser",
    PATH: "/manage-user",
    LABEL: "หน้าจัดการผู้ใช้งาน",
    COMPONENT: ManageUser,
    ELEMENT: () => <ManageUser />,
  },
};

export type RoutesNoAuthType = {
  ROUTE_LOGIN: RouteType;
};

export const ROUTES_NO_AUTH: RoutesNoAuthType = {
  ROUTE_LOGIN: {
    KEY: "login",
    PATH: "/",
    LABEL: "เข้าสู่ระบบ",
    COMPONENT: Login,
    ELEMENT: () => <Login />,
  },
};
