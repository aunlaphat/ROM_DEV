export const getStored = (filedName: string) => {
  const value = localStorage.getItem(filedName);
  return value && localStorage.getItem(filedName) !== "undefined"
    ? JSON.parse(value)
    : "";
};

export const setStored = (filedName: string, value: any) => {
  return localStorage.setItem(filedName, JSON.stringify(value));
};

export const removeStored = (filedName: string) => {
  return localStorage.removeItem(filedName);
};

export const clearStored = () => {
  return localStorage.clear();
};
