/* This matches the width at which submit buttons go full-width for mobile devices */
export const mobileSize = 481;
/* REACT_APP_NODE_ENV allows you to override NODE_ENV when using `yarn build`
 * For more information see https://facebook.github.io/create-react-app/docs/adding-custom-environment-variables
 */
export const isProduction = process.env.NODE_ENV === 'production';
export const isDevelopment = process.env.NODE_ENV === 'development' || process.env.REACT_APP_NODE_ENV === 'development';
export const isTest = process.env.NODE_ENV === 'test';

export const gitBranch = process.env.REACT_APP_GIT_BRANCH || 'unknown';
export const gitSha = process.env.REACT_APP_GIT_COMMIT || 'unknown';

export const NULL_UUID = '00000000-0000-0000-0000-000000000000';

export const hostname = window && window.location && window.location.hostname;
export const isMilmoveSite = hostname.startsWith('my') || hostname.startsWith('mil') || '';
export const isOfficeSite = hostname.startsWith('office') || '';
export const isAdminSite = hostname.startsWith('admin') || '';
export const technicalHelpDeskURL =
  'https://www.militaryonesource.mil/resources/gov/customer-service-contacts-for-military-pcs/#technical-help-desk';

export function serviceName() {
  if (isAdminSite) {
    return 'admin';
  }
  if (isOfficeSite) {
    return 'office';
  }
  return 'my';
}

export const titleCase = (str) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

export const MOVE_STATUSES = {
  DRAFT: 'DRAFT',
  SUBMITTED: 'SUBMITTED',
  APPROVED: 'APPROVED',
  CANCELED: 'CANCELED',
  NEEDS_SERVICE_COUNSELING: 'NEEDS SERVICE COUNSELING',
  SERVICE_COUNSELING_COMPLETED: 'SERVICE COUNSELING COMPLETED',
  APPROVALS_REQUESTED: 'APPROVALS REQUESTED',
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

export const WEIGHT_TICKET_SET_TYPE = {
  CAR: 'CAR',
  CAR_TRAILER: 'CAR_TRAILER',
  BOX_TRUCK: 'BOX_TRUCK',
  PRO_GEAR: 'PRO_GEAR',
};

export const UPLOAD_SCAN_STATUS = {
  CLEAN: 'CLEAN',
  INFECTED: 'INFECTED',
  PROCESSING: 'PROCESSING',
};

export const CONUS_STATUS = {
  CONUS: 'CONUS',
  OCONUS: 'OCONUS',
};

export const SHIPMENT_OPTIONS = {
  HHG: 'HHG',
  PPM: 'PPM',
  NTS: 'HHG_INTO_NTS_DOMESTIC',
  NTSR: 'HHG_OUTOF_NTS_DOMESTIC',
  BOAT: 'BOAT',
  BOAT_HAUL_AWAY: 'BOAT',
  BOAT_TOW_AWAY: 'BOAT',
  MOBILE_HOME: 'MOBILE_HOME',
  UNACCOMPANIED_BAGGAGE: 'UNACCOMPANIED_BAGGAGE',
};

export const MARKET_CODES = {
  INTERNATIONAL: 'i',
  DOMESTIC: 'd',
};

export const SHIPMENT_TYPES = {
  HHG: 'HHG',
  PPM: 'PPM',
  NTS: 'HHG_INTO_NTS_DOMESTIC',
  NTSR: 'HHG_OUTOF_NTS_DOMESTIC',
  BOAT_HAUL_AWAY: 'BOAT_HAUL_AWAY',
  BOAT_TOW_AWAY: 'BOAT_TOW_AWAY',
  MOBILE_HOME: 'MOBILE_HOME',
  UNACCOMPANIED_BAGGAGE: 'UNACCOMPANIED_BAGGAGE',
};

// These constants are used for forming URLs that have the shipment type in
// them so that they are human readable.
export const SHIPMENT_OPTIONS_URL = {
  HHG: 'HHG',
  PPM: 'PPM',
  NTS: 'NTS',
  NTSrelease: 'NTSrelease',
  BOAT: 'Boat',
  MOBILE_HOME: 'Mobilehome',
  UNACCOMPANIED_BAGGAGE: 'UB',
};

export const LOA_TYPE = {
  HHG: 'HHG',
  NTS: 'NTS',
};

export const SIGNED_CERT_OPTIONS = {
  SHIPMENT: 'SHIPMENT',
  PPM_PAYMENT: 'PPM_PAYMENT',
};

export const shipmentOptionLabels = [
  { key: SHIPMENT_OPTIONS.NTSR, label: 'NTS-release' },
  { key: SHIPMENT_OPTIONS.NTS, label: 'NTS' },
  { key: SHIPMENT_OPTIONS.HHG, label: 'HHG' },
  { key: SHIPMENT_OPTIONS.PPM, label: 'PPM' },
  { key: SHIPMENT_OPTIONS.BOAT, label: 'Boat' },
  { key: SHIPMENT_OPTIONS.MOBILE_HOME, label: 'Mobile Home' },
  { key: SHIPMENT_TYPES.BOAT_HAUL_AWAY, label: 'Boat' },
  { key: SHIPMENT_TYPES.BOAT_TOW_AWAY, label: 'Boat' },
  { key: SHIPMENT_TYPES.UNACCOMPANIED_BAGGAGE, label: 'UB' },
];

export const SERVICE_ITEM_STATUS = {
  SUBMITTED: 'SUBMITTED',
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
};

export const MTO_SERVICE_ITEM_STATUS = {
  SUBMITTED: 'Move Task Order Requested',
  APPROVED: 'Move Task Order Approved',
  REJECTED: 'Move Task Order Rejected',
};

export const PAYMENT_REQUEST_STATUS = {
  PENDING: 'PENDING',
  REVIEWED: 'REVIEWED',
  SENT_TO_GEX: 'SENT_TO_GEX',
  TPPS_RECEIVED: 'TPPS_RECEIVED',
  PAID: 'PAID',
  REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED: 'REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED',
  EDI_ERROR: 'EDI_ERROR',
  DEPRECATED: 'DEPRECATED',
  PAYMENT_REQUESTED: 'PAYMENT_REQUESTED',
};

export const PAYMENT_SERVICE_ITEM_STATUS = {
  DENIED: 'DENIED',
  REQUESTED: 'REQUESTED',
  APPROVED: 'APPROVED',
};

export const WEBHOOK_SUBSCRIPTION_STATUS = {
  ACTIVE: 'ACTIVE',
  DISABLED: 'DISABLED',
  FAILING: 'FAILING',
};

export const MTOAgentType = {
  RELEASING: 'RELEASING_AGENT',
  RECEIVING: 'RECEIVING_AGENT',
};

export const MoveOrderDocumentType = {
  ALL: 'ALL',
  ORDERS: 'ORDERS',
  AMENDMENTS: 'AMENDMENTS',
};

// These constants are used to track network requests using component state
export const isError = 'REQUEST_ERROR';
export const isLoading = 'REQUEST_LOADING';
export const isSuccess = 'REQUEST_SUCCESS';

// documentSizeLimitMsg is used in several files around document upload
export const documentSizeLimitMsg = 'Please keep each file under 25MB.';

// new window dimensions in relation to current window
export const defaultRelativeWindowSize = 2 / 3;

// The date only used to format date strings
export const DATE_FORMAT_STRING = 'DD MMM YYYY';
// The date and time used to format date strings
export const DATE_TIME_FORMAT_STRING = 'DD MMM YYYY, hh:mm a';

export const DEFAULT_EMPTY_VALUE = '—'; // emdash

export const FEATURE_FLAG_KEYS = {
  PPM: 'ppm',
  NTS: 'nts',
  NTSR: 'ntsr',
  BOAT: 'boat',
  MOBILE_HOME: 'mobile_home',
  UNACCOMPANIED_BAGGAGE: 'unaccompanied_baggage',
  ENABLE_ALASKA: 'enable_alaska',
};

export const MOVE_DOCUMENT_TYPE = {
  ORDERS: 'ORDERS',
  AMENDMENTS: 'AMENDMENTS',
  SUPPORTING: 'SUPPORTING',
};

export const ADDRESS_TYPES = {
  PICKUP: 'pickupAddress',
  SECOND_PICKUP: 'secondaryPickupAddress',
  THIRD_PICKUP: 'tertiaryPickupAddress',
  DESTINATION: 'destinationAddress',
  SECOND_DESTINATION: 'secondaryDeliveryAddress',
  THIRD_DESTINATION: 'tertiaryDeliveryAddress',
};

const ADDRESS_LABELS_MAP = {
  [ADDRESS_TYPES.PICKUP]: 'Pickup Address',
  [ADDRESS_TYPES.SECOND_PICKUP]: 'Second Pickup Address',
  [ADDRESS_TYPES.THIRD_PICKUP]: 'Third Pickup Address',
  [ADDRESS_TYPES.DESTINATION]: 'Delivery Address',
  [ADDRESS_TYPES.SECOND_DESTINATION]: 'Second Delivery Address',
  [ADDRESS_TYPES.THIRD_DESTINATION]: 'Third Delivery Address',
};

export const getAddressLabel = (type) => ADDRESS_LABELS_MAP[type];
