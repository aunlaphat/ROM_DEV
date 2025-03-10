import { Home } from "../screens/Home";
import { Login } from "../screens/auth";
import { NotFound } from "../screens/NotFound";
import CreateReturnOrderMKP from "../screens/Orders/Marketplace";
import { DraftAndConfirm } from "../screens/Draft&Confirm";
import { ConfirmReturnTrade } from "../screens/ConfirmReturnTrade";
import { Report } from "../screens/Report";
import { ManageUser } from "../screens/ManageUser";

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
  ROUTE_CREATERETURNORDERMKP: RouteType;
  // ROUTE_IJ: RouteType;
  // ROUTE_CREATETRADERETURN: RouteType;
  ROUTE_CONFIRMRETURNTRADE: RouteType;
  // ROUTE_SALERETURN: RouteType;
  // ROUTE_OTHERRETURN: RouteType;
  // ROUTE_CREATEBLINDRETURN: RouteType;
  // ROUTE_TAKEPICTURE: RouteType;
  ROUTE_REPORT: RouteType;
  ROUTE_DRAFTANDCONFIRM: RouteType;
  ROUTE_MANAGEUSER: RouteType;
};

export const ROUTES: RoutesType = {
  ROUTE_MAIN: {
    KEY: "home",
    PATH: "/home",
    LABEL: "Home",
    COMPONENT: Home,
    ELEMENT: () => <Home />,
  },
  ROUTE_NOTFOUND: {
    KEY: "notFound",
    PATH: "*",
    LABEL: "Page Not Found",
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
  ROUTE_CREATERETURNORDERMKP: {
    KEY: "createReturnOrderMKP",
    PATH: "/create-return-order-mkp",
    LABEL: "Create Return Order MKP",
    COMPONENT: CreateReturnOrderMKP,
    ELEMENT: () => <CreateReturnOrderMKP />,
  },
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
  ROUTE_CONFIRMRETURNTRADE: {
    KEY: "confirmReturnTrade",
    PATH: "/confirm-return-trade",
    LABEL: "Confirm Return Trade",
    COMPONENT: ConfirmReturnTrade,
    ELEMENT: () => <ConfirmReturnTrade />,
  },
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
    LABEL: "Report",
    COMPONENT: Report,
    ELEMENT: () => <Report />,
  },
  ROUTE_DRAFTANDCONFIRM: {
    KEY: "draftAndConfirm",
    PATH: "/draft-and-confirm",
    LABEL: "Draft and Confirm",
    COMPONENT: DraftAndConfirm,
    ELEMENT: () => <DraftAndConfirm />,
  },
  ROUTE_MANAGEUSER: {
    KEY: "manageUser",
    PATH: "/manage-user",
    LABEL: "Manage User",
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
    LABEL: "Login",
    COMPONENT: Login,
    ELEMENT: () => <Login />,
  },
};
