export const generalRoutes = {
  HOME_PATH: '/',
  SIGN_IN_PATH: '/sign-in',
  PRIVACY_SECURITY_POLICY_PATH: '/privacy-security',
  ACCESSIBILITY_PATH: '/accessibility',
};

export const customerRoutes = {
  CONUS_OCONUS_PATH: '/service-member/conus-oconus',
  DOD_INFO_PATH: '/service-member/dod-info',
  NAME_PATH: '/service-member/name',
  CONTACT_INFO_PATH: '/service-member/contact-info',
  CURRENT_DUTY_LOCATION_PATH: '/service-member/current-duty',
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
  SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH: '/moves/:moveId/shipments/:mtoShipmentId/estimated-weight',
  SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH: '/moves/:moveId/shipments/:mtoShipmentId/estimated-incentive',
  SHIPMENT_PPM_ADVANCES_PATH: '/moves/:moveId/shipments/:mtoShipmentId/advances',
  SHIPMENT_PPM_ABOUT_PATH: '/moves/:moveId/shipments/:mtoShipmentId/about',
  SHIPMENT_PPM_WEIGHT_TICKETS_PATH: '/moves/:moveId/shipments/:mtoShipmentId/weight-tickets',
  SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH: '/moves/:moveId/shipments/:mtoShipmentId/weight-tickets/:weightTicketId',
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
  SHIPMENT_ADD_PATH: `${BASE_MOVE_PATH}/new-:shipmentType`,
  SHIPMENT_EDIT_PATH: `${BASE_MOVE_PATH}/shipments/:shipmentId`,
  MOVE_HISTORY_PATH: `${BASE_MOVE_PATH}/history`,
  CUSTOMER_SUPPORT_REMARKS_PATH: `/counseling/moves/:moveCode/customer-support-remarks`,
};

export const tioRoutes = {
  PAYMENT_REQUESTS_PATH: '/moves/:moveCode/payment-requests',
  BILLABLE_WEIGHT_PATH: `/moves/:moveCode/billable-weight`,
};

export const tooRoutes = {
  SHIPMENT_EDIT_PATH: '/moves/:moveCode/shipments/:shipmentId',
  MOVE_VIEW_PATH: '/moves/:moveCode/details',
  ORDERS_EDIT_PATH: '/moves/:moveCode/orders',
};

export const qaeCSRRoutes = {
  MOVE_SEARCH_PATH: '/qaecsr/search',
};

// README: Prime API endpoints accept either a Move code or ID.
// The base path doesn't have prime-simulator due to matching issues with /prime.
const BASE_PRIME_SIMULATOR_PATH = '/simulator/moves/:moveCodeOrID';

export const primeSimulatorRoutes = {
  VIEW_MOVE_PATH: `${BASE_PRIME_SIMULATOR_PATH}/details`,
  UPDATE_SHIPMENT_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId`,
  CREATE_PAYMENT_REQUEST_PATH: `${BASE_PRIME_SIMULATOR_PATH}/payment-requests/new`,
  CREATE_SERVICE_ITEM_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/service-items/new`,
  UPLOAD_DOCUMENTS_PATH: `${BASE_PRIME_SIMULATOR_PATH}/payment-requests/:paymentRequestId/upload`,
  SHIPMENT_UPDATE_ADDRESS_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/addresses/update`,
  SHIPMENT_UPDATE_REWEIGH_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/reweigh/:reweighId/update`,
};
