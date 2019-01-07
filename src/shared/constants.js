/* This matches the width at which submit buttons go full-width for mobile devices */
export const mobileSize = 481;
/* REACT_APP_NODE_ENV allows you to override NODE_ENV when using `yarn build`
 * For more information see https://facebook.github.io/create-react-app/docs/adding-custom-environment-variables
 */
export const isProduction = process.env.NODE_ENV === 'production';
export const isDevelopment = process.env.NODE_ENV === 'development' || process.env.REACT_APP_NODE_ENV === 'development';
export const isTest = process.env.NODE_ENV === 'test';
export const NULL_UUID = '00000000-0000-0000-0000-000000000000';
export const ppmInfoPacket = '/downloads/ppm_info_sheet.pdf';
export const hhgInfoPacket = '/downloads/hhg_info_sheet.pdf';

export const hostname = window && window.location && window.location.hostname;
export const isOfficeSite = hostname.startsWith('office') || '';
export const isTspSite = hostname.startsWith('tsp') || '';

export const titleCase = str => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

// These constants are used to track network requests using component state
export const isError = 'REQUEST_ERROR';
export const isLoading = 'REQUEST_LOADING';
export const isSuccess = 'REQUEST_SUCCESS';
