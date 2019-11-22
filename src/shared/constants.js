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

export const hostname = window && window.location && window.location.hostname;
export const isMilmoveSite = hostname.startsWith('my') || hostname.startsWith('mil') || '';
export const isOfficeSite = hostname.startsWith('office') || '';
export const isTspSite = hostname.startsWith('tsp') || '';
export const isAdminSite = hostname.startsWith('admin') || '';
export const isSystemAdminSite = isAdminSite; // once we start building program admin, we can flesh this out

export const titleCase = str => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

export const MOVE_DOC_TYPE = {
  WEIGHT_TICKET_SET: 'WEIGHT_TICKET_SET',
  EXPENSE: 'EXPENSE',
  GBL: 'GOV_BILL_OF_LADING',
};

export const MOVE_DOC_STATUS = {
  OK: 'OK',
  AWAITING_REVIEW: 'AWAITING_REVIEW',
  HAS_ISSUE: 'HAS_ISSUE',
  EXCLUDE: 'EXCLUDE_FROM_CALCULATION',
};

export const UPLOAD_SCAN_STATUS = {
  CLEAN: 'CLEAN',
  INFECTED: 'INFECTED',
  PROCESSING: 'PROCESSING',
};

// These constants are used to track network requests using component state
export const isError = 'REQUEST_ERROR';
export const isLoading = 'REQUEST_LOADING';
export const isSuccess = 'REQUEST_SUCCESS';

// documentSizeLimitMsg is used in several files around document upload
export const documentSizeLimitMsg = 'Please keep each file under 25MB.';

// new window dimensions in relation to current window
export const defaultRelativeWindowSize = 2 / 3;
