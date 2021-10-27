export const generalRoutes = {
  HOME_PATH: '/',
  SIGN_IN_PATH: '/sign-in',
  PRIVACY_SECURITY_POLICY_PATH: '/privacy-security',
  ACCESSIBILITY_PATH: '/accessibility',
};

export const customerRoutes = {
  ACCESS_CODE_PATH: '/access-code',
  CONUS_OCONUS_PATH: '/service-member/conus-oconus',
  DOD_INFO_PATH: '/service-member/dod-info',
  NAME_PATH: '/service-member/name',
  CONTACT_INFO_PATH: '/service-member/contact-info',
  CURRENT_DUTY_STATION_PATH: '/service-member/current-duty',
  CURRENT_ADDRESS_PATH: '/service-member/current-address',
  BACKUP_ADDRESS_PATH: '/service-member/backup-address',
  BACKUP_CONTACTS_PATH: '/service-member/backup-contact',
  ORDERS_INFO_PATH: '/orders/info',
  ORDERS_UPLOAD_PATH: '/orders/upload',
  ORDERS_AMEND_PATH: '/orders/amend',
  ORDERS_EDIT_PATH: '/moves/:moveId/review/edit-orders',
  SHIPMENT_MOVING_INFO_PATH: '/moves/:moveId/moving-info',
  SHIPMENT_SELECT_TYPE_PATH: '/moves/:moveId/shipment-type',
  SHIPMENT_CREATE_PATH: '/moves/:moveId/new-shipment',
  SHIPMENT_EDIT_PATH: '/moves/:moveId/shipments/:mtoShipmentId/edit',
  MOVE_REVIEW_PATH: '/moves/:moveId/review',
  MOVE_AGREEMENT_PATH: '/moves/:moveId/agreement',
  EDIT_PROFILE_PATH: '/moves/review/edit-profile',
  PROFILE_PATH: '/service-member/profile',
  SERVICE_INFO_EDIT_PATH: '/moves/review/edit-service-info',
  CONTACT_INFO_EDIT_PATH: '/moves/review/edit-contact-info',
};

const BASE_MOVE_PATH = '/counseling/moves/:moveCode';

export const servicesCounselingRoutes = {
  ALLOWANCES_EDIT_PATH: `${BASE_MOVE_PATH}/allowances`,
  BASE_MOVE_PATH,
  CUSTOMER_INFO_EDIT_PATH: `${BASE_MOVE_PATH}/customer`,
  MOVE_VIEW_PATH: `${BASE_MOVE_PATH}/details`,
  ORDERS_EDIT_PATH: `${BASE_MOVE_PATH}/orders`,
  QUEUE_VIEW_PATH: '/counseling/queue',
  SHIPMENT_ADD_PATH: `${BASE_MOVE_PATH}/new-HHG`,
  SHIPMENT_EDIT_PATH: `${BASE_MOVE_PATH}/shipments/:shipmentId`,
};

export const tioRoutes = {
  PAYMENT_REQUESTS_PATH: '/moves/:moveCode/payment-requests',
  BILLABLE_WEIGHT_PATH: `/moves/:moveCode/billable-weight`,
};

// README: Prime API endpoints accept either a Move code or ID.
// The base path doesn't have prime-simulator due to matching issues with /prime.
const BASE_PRIME_SIMULATOR_PATH = '/simulator/moves/:moveCodeOrID';

export const primeSimulatorRoutes = {
  VIEW_MOVE_PATH: `${BASE_PRIME_SIMULATOR_PATH}/details`,
  UPDATE_SHIPMENT_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId`,
  CREATE_PAYMENT_REQUEST_PATH: `${BASE_PRIME_SIMULATOR_PATH}/payment-requests/new`,
  UPLOAD_DOCUMENTS_PATH: `${BASE_PRIME_SIMULATOR_PATH}/payment-requests/:paymentRequestId/upload`,
};
