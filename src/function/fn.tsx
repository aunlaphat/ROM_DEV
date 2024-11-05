export const delay = (ms = 500) => {
  return new Promise((resolve: any) => setTimeout(resolve, ms));
};
