export const createImgUrl = (imgName: string) =>
  import.meta.env.VITE_APP_API_URL + imgName;

export const createPdfUrl = (fileName: string) =>
  import.meta.env.VITE_APP_API_URL + '/pdf/' + fileName;
