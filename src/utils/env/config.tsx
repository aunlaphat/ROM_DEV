export const env = {
  title: process.env.REACT_APP_TITLE,
  url: process.env.REACT_APP_URL_HOST,
  node: process.env.REACT_APP_NODE_ENV,
  node_env: process.env.REACT_APP_NODE_ENV,
  port: process.env.REACT_APP_PORT,
  version: process.env.REACT_APP_VERSION,
  url_dev: process.env.REACT_APP_URL_DEV,
  url_prd: process.env.REACT_APP_URL_PRD,
  key_storage: process.env.REACT_APP_KEY_STORAGE ?? "",
  key_hash: process.env.REACT_APP_KEY_HASH,
  api_base_url: process.env.REACT_APP_URL_DEV || "http://localhost:8000", // ใช้ค่า dev หรือ default 8000
};

