import axios from "axios";
//import { env } from "../utils/env/config";
import { getCookies } from "../store/useCookies";

// export const MODE_API: any = {
//   DEVELOPMENT: env.url_dev,
//   PRODUCTION: env.url_prd,
// };

export const CONNECT_API = process.env.REACT_APP_BACKEND_URL;

const api = axios.create({
  baseURL: CONNECT_API,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000,
});

const apiupload = axios.create({
  baseURL: CONNECT_API,
  withCredentials: true,
  headers: {
    "Content-Type": "multipart/form-data",
  },
  timeout: 10000,
});

export const GET = async (path: string, header?: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await api.get(`${path}`, header ? header : config);
    let data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const POST = async (path: string, data: any, header?: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await api.post(`${path}`, data, header ? header : config);
    data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const PATCH = async (path: string, data: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await api.patch(`${path}`, data, config);
    data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const PUT = async (path: string, data: any, header?: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await api.put(`${path}`, data, header ? header : config);
    data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const DELETE = async (path: string, data: any, header?: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await api.delete(`${path}`, {
      headers: header ? header : config,
      data,
    });
    data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const UPLODA = async (
  path: string,
  formdata: FormData,
  header?: any
) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await apiupload.post(
      `${path}`,
      formdata,
      header ? header : config
    );
    let data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};
