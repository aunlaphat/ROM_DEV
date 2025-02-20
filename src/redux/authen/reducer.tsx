import * as type from "./types";

type User = {
  userID: string;
  userFullName: string;
  userRoleID: string;
};

type State = {
  users: User[];
  loading: boolean;
};

const initialState: State = {
  users: [],
  loading: false,
};

export default function authenReducer(state = initialState, action: any) {
  switch (action.type) {
    case type.AUTHEN_LOGIN_REQ:
    case type.AUTHEN_LOGIN_LARK_REQ:
      return {
        ...state,
        loading: true,
      };
    case type.AUTHEN_LOGIN_SUCCESS:
    case type.AUTHEN_LOGIN_LARK_SUCCESS:
      return {
        ...state,
        loading: false,
        users: action.users,
      };
    case type.AUTHEN_LOGIN_FAIL:
    case type.AUTHEN_LOGIN_LARK_FAIL:
      return {
        ...state,
        loading: false,
        error: action.message,
      };
    case type.AUTHEN_LOGOUT_SUCCESS:
      return {
        ...state,
        loading: false,
        users: [],
      };
    case type.AUTHEN_LOGOUT_FAIL:
      return {
        ...state,
        loading: false,
        error: action.message,
      };
    case type.AUTHEN_CHECK_REQ:
      return {
        ...state,
        loading: true,
      };
    case type.AUTHEN_CHECK_SUCCESS:
      return {
        ...state,
        loading: false,
        users: action.users,
      };
    case type.AUTHEN_CHECK_FAIL:
      return {
        ...state,
        loading: false,
        error: action.message,
      };
    default:
      return state;
  }
}
