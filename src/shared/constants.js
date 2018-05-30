/* This matches the width at which submit buttons go full-width for mobile devices */
export const mobileSize = 481;
export const isProduction = process.env.NODE_ENV === 'production';
export const isDevelopment = process.env.NODE_ENV === 'development';
export const isTest = process.env.NODE_ENV === 'test';
export const NULL_UUID = '00000000-0000-0000-0000-000000000000';
export const ppmInfoPacket =
  'https://s3.amazonaws.com/prod.tracker2/resource/89609252/PPM_Info_Sheet.pdf?response-content-disposition=inline&response-content-type=application%2Fpdf&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAJJBSFJ4TCVKKGAIA%2F20180530%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20180530T230724Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=c66580dd5fe891c4fde94a4a28faceac307ac41d559c90d99719cfbb871c788f';
