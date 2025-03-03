import Home from "../screens/home";
import LoginScene from "../screens/Authen/LoginScene";
import NotfoundScene from "../screens/notFound";
import { env } from '../utils/env/config';
import ImportOrder from "../screens/Return_import/Import_Return_Order";
import ReturnOrder from "../screens/Return/Returnorder";
import CreateReturn from "../screens/CreateReturn/CreateReturn";
import SRPage from "../screens/CreateReturn/SR";
import IJPage from "../screens/CreateReturn/IJ";
import CreateTradeReturn from "../screens/CreateTradeReturn/CreateTradeReturn";
import ConfirmReturnTrade from "../screens/ConfirmReturnTrade/ConfirmReturnTrade";
import OtherReturn from "../screens/SaleReturn/SR/OtherReturn";
import SaleReturn from "../screens/SaleReturn/Sale_Return";
import CreateBlind from "../screens/Create_Blind/CreateBlindReturn";
import Takepicture from "../screens/Create_Blind/Takepicture";
import Report from "../screens/Report/Report";
import DraftandConfirm from "../screens/Draft&Confirm/Draft&Confirm";
import ManageUser from "../screens/ManageUser/Manageuser";

export const ROUTE_LOGIN = process.env.REACT_APP_FRONTEND_URL+"/";

type Route = {
  KEY: string;                // ชื่อ Key ของ Route (ใช้เป็นตัวระบุ)
  PATH: string;               // URL Path เช่น "/home"
  LABEL: string;              // ชื่อ Route ที่แสดงผลใน UI
  COMPONENT: React.ComponentType; // Component ที่จะใช้ใน Route
  ELEMENT: () => JSX.Element; // ฟังก์ชันที่ Return JSX
};

export type Routes = {
  ROUTE_MAIN: Route;
  ROUTE_NOTFOUND: Route;
  ROUTE_IMPORTORDER: Route;
  ROUTE_RETURNORDER: Route;
  ROUTE_CREATERETURN: Route;
  ROUTE_SR: Route;
  ROUTE_IJ: Route;
  ROUTE_CREATETRADERETURN: Route;
  ROUTE_CONFIRMRETURNTRADE: Route;
  ROUTE_SALERETURN: Route;
  ROUTE_OTHER: Route;
  ROUTE_CREATEBLIND: Route;
  ROUTE_TAKEPICTURE: Route;
  ROUTE_REPORT: Route;
  ROUTE_DRAFTANDCONFIRM: Route;
  ROUTE_MANAGEUSER: Route;
};

export type Routes_Not_Auth = {
  ROUTE_LOGIN: Route;
};

export const ROUTE_NOT_AUTHEN: Routes_Not_Auth = {
  ROUTE_LOGIN: {
    KEY: "login",
    PATH: "/",
    LABEL: "ล็อคอิน",
    COMPONENT: LoginScene,
    ELEMENT: () => <LoginScene />,
  },
};

export const ROUTES_PATH: Routes = {
  ROUTE_MAIN: {
    KEY: "home",
    PATH: "home",
    LABEL: "ใบงานทั้งหมด",
    COMPONENT: Home,
    ELEMENT: () => <Home />,
  },
  ROUTE_NOTFOUND: {
    KEY: "not_found",
    PATH: "*",
    LABEL: "NOT FOUND",
    COMPONENT: NotfoundScene,
    ELEMENT: () => <NotfoundScene />,
  },
  ROUTE_IMPORTORDER: {
    KEY: "import_order",
    PATH: "import_order",
    LABEL: "Import Return order ",
    COMPONENT: ImportOrder,
    ELEMENT: () => <ImportOrder />,
  },
  ROUTE_RETURNORDER: {
    KEY: "return_order",
    PATH: "return_order",
    LABEL: "Home",
    COMPONENT: ReturnOrder,
    ELEMENT: () => < ReturnOrder />,
  },
  ROUTE_CREATERETURN: {
    KEY: "CreateReturn",
    PATH: "CreateReturn",
    LABEL: "Create Return",
    COMPONENT: CreateReturn,
    ELEMENT: () => < CreateReturn />,
  },
  ROUTE_SR: {
    KEY: "SR",
    PATH: "SR",
    LABEL: "SR",
    COMPONENT: SRPage,
    ELEMENT: () => < SRPage />,
  },
  ROUTE_IJ: {
    KEY: "IJ",
    PATH: "IJ",
    LABEL: "IJ",
    COMPONENT: IJPage,
    ELEMENT: () => < IJPage />,
  },
  ROUTE_CREATETRADERETURN: {
    KEY: "CreateTrandReturn",
    PATH: "CreateTrandReturn",
    LABEL: "Create Trand Return",
    COMPONENT: CreateTradeReturn,
    ELEMENT: () => < CreateTradeReturn />,
  },
  ROUTE_CONFIRMRETURNTRADE: {
    KEY: "ConfirmReturnTrade",
    PATH: "ConfirmReturnTrade",
    LABEL: "Confirm Return Trade",
    COMPONENT: ConfirmReturnTrade,
    ELEMENT: () => < ConfirmReturnTrade />,
  },
  ROUTE_OTHER: {
    KEY: "Other",
    PATH: "Other",
    LABEL: "SR/IJ และ อื่นๆ Return ",
    COMPONENT: OtherReturn,
    ELEMENT: () => <OtherReturn />,
  },
  ROUTE_SALERETURN: {
    KEY: "SaleReturn",
    PATH: "SaleReturn",
    LABEL: "Sale Return",
    COMPONENT: SaleReturn,
    ELEMENT: () => <SaleReturn />,
  },
  ROUTE_CREATEBLIND: {
    KEY: "CreateBlindReturn",
    PATH: "CreateBlindReturn",
    LABEL: "Create Build Return",
    COMPONENT: CreateBlind,
    ELEMENT: () => <CreateBlind />,
  },
  ROUTE_TAKEPICTURE: {
    KEY: "Takepicture",
    PATH: "Takepicture",
    LABEL: "Takepicture",
    COMPONENT: Takepicture,
    ELEMENT: () => <Takepicture />,
  },
  ROUTE_REPORT: {
    KEY: "Report",
    PATH: "Report",
    LABEL: "Report",
    COMPONENT: Report,
    ELEMENT: () => <Report />,
  },
  ROUTE_DRAFTANDCONFIRM: {
    KEY: "DraftandConfirm",
    PATH: "DraftandConfirm",
    LABEL: "Draft&Confirm",
    COMPONENT: DraftandConfirm,
    ELEMENT: () => <DraftandConfirm />,
  },
  ROUTE_MANAGEUSER: {
    KEY: "Manageuser",
    PATH: "Manageuser",
    LABEL: "Manage User",
    COMPONENT: ManageUser,
    ELEMENT: () => <ManageUser />,
  },

};

export type RoutesWorker = {
  ROUTE_MAIN: Route;
  ROUTE_NOTFOUND: Route;
  ROUTE_IMPORTORDER: Route;
  ROUTE_RETURNORDER: Route;
  ROUTE_CREATERETURN: Route;
  ROUTE_SR: Route;
  ROUTE_IJ: Route;
  ROUTE_CREATETRADERETURN: Route;
  ROUTE_CONFIRMRETURNTRADE: Route;
  ROUTE_OTHER: Route;
  ROUTE_SALERETURN: Route;
  ROUTE_CREATEBLIND: Route;
  ROUTE_TAKEPICTURE: Route;
  ROUTE_REPORT: Route;
  ROUTE_DRAFTANDCONFIRM: Route;
  ROUTE_MANAGEUSER: Route;
};

export const ROUTES_PATH_WORKER: RoutesWorker = {
  ROUTE_MAIN: {
    KEY: "home",
    PATH: "/home",
    LABEL: "ใบงานทั้งหมด",
    COMPONENT: Home,
    ELEMENT: () => <Home />,
  },
  ROUTE_NOTFOUND: {
    KEY: "not_found",
    PATH: "/*",
    LABEL: "NOT FOUND",
    COMPONENT: NotfoundScene,
    ELEMENT: () => <NotfoundScene />,
  },
  ROUTE_IMPORTORDER: {
    KEY: "import_order",
    PATH: "import_order",
    LABEL: "Import Return order ",
    COMPONENT: NotfoundScene,
    ELEMENT: () => <ImportOrder />,
  },
  ROUTE_RETURNORDER: {
    KEY: "return_order",
    PATH: "/return_order",
    LABEL: "return order",
    COMPONENT: ReturnOrder,
    ELEMENT: () => < ReturnOrder />,
  },
  ROUTE_CREATERETURN: {
    KEY: "CreateReturn",
    PATH: "/CreateReturn",
    LABEL: "Create Return",
    COMPONENT: CreateReturn,
    ELEMENT: () => < CreateReturn />,
  },
  ROUTE_SR: {
    KEY: "SR",
    PATH: "/SR",
    LABEL: "SR",
    COMPONENT: SRPage,
    ELEMENT: () => < SRPage />,
  },
  ROUTE_IJ: {
    KEY: "IJ",
    PATH: "/IJ",
    LABEL: "IJ",
    COMPONENT: IJPage,
    ELEMENT: () => < IJPage />,
  },
  ROUTE_CREATETRADERETURN: {
    KEY: "CreateTrandReturn",
    PATH: "/CreateTrandReturn",
    LABEL: "Create Trand Return",
    COMPONENT: CreateTradeReturn,
    ELEMENT: () => < CreateTradeReturn />,
  },
  ROUTE_CONFIRMRETURNTRADE: {
    KEY: "ConfirmReturnTrade",
    PATH: "/ConfirmReturnTrade",
    LABEL: "Confirm Return Trade",
    COMPONENT: ConfirmReturnTrade,
    ELEMENT: () => < ConfirmReturnTrade />,
  },
  ROUTE_OTHER: {
    KEY: "Other",
    PATH: "/Other",
    LABEL: "SR/IJ และ อื่นๆ Return ",
    COMPONENT: OtherReturn,
    ELEMENT: () => < OtherReturn />,
  },
  ROUTE_SALERETURN: {
    KEY: "SaleReturn",
    PATH: "/SaleReturn",
    LABEL: "Sale Return",
    COMPONENT: SaleReturn,
    ELEMENT: () => <SaleReturn />,
  },
  ROUTE_CREATEBLIND: {
    KEY: "CreateBlindReturn",
    PATH: "/CreateBlindReturn",
    LABEL: "Create Build Return",
    COMPONENT: CreateBlind,
    ELEMENT: () => <CreateBlind />,
  },
  ROUTE_TAKEPICTURE: {
    KEY: "Takepicture",
    PATH: "/Takepicture",
    LABEL: "Takepicture",
    COMPONENT: Takepicture,
    ELEMENT: () => <Takepicture />,
  },
  ROUTE_REPORT: {
    KEY: "Report",
    PATH: "/Report",
    LABEL: "Report",
    COMPONENT: Report,
    ELEMENT: () => <Report />,
  },
  ROUTE_DRAFTANDCONFIRM: {
    KEY: "DraftandConfirm",
    PATH: "/DraftandConfirm",
    LABEL: "Draft&Confirm",
    COMPONENT: DraftandConfirm,
    ELEMENT: () => <DraftandConfirm />,
  },
  ROUTE_MANAGEUSER: {
    KEY: "Manageuser",
    PATH: "/Manageuser",
    LABEL: "Manage User",
    COMPONENT: ManageUser,
    ELEMENT: () => <ManageUser />,
  },

};

export type RoutesNoPermission = {
  ROUTE_NOTFOUND: Route;
};

export const ROUTES_PATH_NOPERMISSION: RoutesNoPermission = {
  ROUTE_NOTFOUND: {
    KEY: "not_found",
    PATH: "/*",
    LABEL: "NOT FOUND",
    COMPONENT: NotfoundScene,
    ELEMENT: () => <NotfoundScene />,
  },
};
