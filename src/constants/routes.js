export const generalRoutes = {
  HOME_PATH: '/',
  SIGN_IN_PATH: '/sign-in',
  PRIVACY_SECURITY_POLICY_PATH: '/privacy-and-security-policy',
  ACCESSIBILITY_PATH: '/accessibility',
};

export const customerRoutes = {
  MOVE_HOME_PAGE: '/move',
  MOVE_HOME_PATH: '/move/:moveId',
  CONUS_OCONUS_PATH: '/service-member/conus-oconus',
  DOD_INFO_PATH: '/service-member/dod-info',
  NAME_PATH: '/service-member/name',
  CONTACT_INFO_PATH: '/service-member/contact-info',
  CURRENT_ADDRESS_PATH: '/service-member/current-address',
  BACKUP_ADDRESS_PATH: '/service-member/backup-address',
  BACKUP_CONTACTS_PATH: '/service-member/backup-contact',
  ORDERS_ADD_PATH: '/orders/add',
  ORDERS_INFO_PATH: '/orders/info/:orderId',
  ORDERS_UPLOAD_PATH: '/orders/upload/:orderId',
  ORDERS_AMEND_PATH: '/orders/amend/:orderId',
  ORDERS_EDIT_PATH: '/move/:moveId/review/edit-orders/:orderId',
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
  SHIPMENT_PPM_REVIEW_PATH: '/moves/:moveId/shipments/:mtoShipmentId/review',
  SHIPMENT_PPM_PRO_GEAR_PATH: '/moves/:moveId/shipments/:mtoShipmentId/pro-gear',
  SHIPMENT_PPM_PRO_GEAR_EDIT_PATH: '/moves/:moveId/shipments/:mtoShipmentId/pro-gear/:proGearId',
  SHIPMENT_PPM_EXPENSES_PATH: '/moves/:moveId/shipments/:mtoShipmentId/expenses',
  SHIPMENT_PPM_EXPENSES_EDIT_PATH: '/moves/:moveId/shipments/:mtoShipmentId/expenses/:expenseId',
  SHIPMENT_PPM_COMPLETE_PATH: '/moves/:moveId/shipments/:mtoShipmentId/complete',
  MOVE_REVIEW_PATH: '/moves/:moveId/review',
  MOVE_AGREEMENT_PATH: '/moves/:moveId/agreement',
  EDIT_PROFILE_PATH: '/moves/review/edit-profile',
  EDIT_OKTA_PROFILE_PATH: '/moves/review/edit-okta-profile',
  PROFILE_PATH: '/service-member/profile',
  SERVICE_INFO_EDIT_PATH: '/moves/review/edit-service-info',
  CONTACT_INFO_EDIT_PATH: '/moves/review/edit-contact-info',
};

const BASE_COUNSELING_MOVE_PATH = '/counseling/moves/:moveCode';

export const servicesCounselingRoutes = {
  BASE_QUEUE_VIEW_PATH: '/counseling/queue',
  QUEUE_VIEW_PATH: 'queue',
  DEFAULT_QUEUE_PATH: '/',
  QUEUE_COUNSELING_PATH: 'counseling',
  BASE_QUEUE_COUNSELING_PATH: '/counseling',
  QUEUE_CLOSEOUT_PATH: 'PPM-closeout',
  BASE_QUEUE_CLOSEOUT_PATH: '/PPM-closeout',
  QUEUE_SEARCH_PATH: 'Search',
  BASE_QUEUE_SEARCH_PATH: '/Search',
  BASE_COUNSELING_MOVE_PATH,
  BASE_ALLOWANCES_EDIT_PATH: `${BASE_COUNSELING_MOVE_PATH}/allowances`,
  ALLOWANCES_EDIT_PATH: 'allowances',
  BASE_CUSTOMER_INFO_EDIT_PATH: `${BASE_COUNSELING_MOVE_PATH}/customer`,
  CUSTOMER_INFO_EDIT_PATH: 'customer',
  BASE_MOVE_VIEW_PATH: `${BASE_COUNSELING_MOVE_PATH}/details`,
  MOVE_VIEW_PATH: 'details',
  BASE_ORDERS_EDIT_PATH: `${BASE_COUNSELING_MOVE_PATH}/orders`,
  ORDERS_EDIT_PATH: 'orders',
  BASE_SHIPMENT_ADD_PATH: `${BASE_COUNSELING_MOVE_PATH}/new-shipment/:shipmentType`,
  SHIPMENT_ADD_PATH: 'new-shipment/:shipmentType',
  BASE_SHIPMENT_EDIT_PATH: `${BASE_COUNSELING_MOVE_PATH}/shipments/:shipmentId`,
  SHIPMENT_EDIT_PATH: 'shipments/:shipmentId',
  BASE_SHIPMENT_ADVANCE_PATH: `${BASE_COUNSELING_MOVE_PATH}/shipments/:shipmentId/advance`,
  SHIPMENT_ADVANCE_PATH: 'shipments/:shipmentId/advance',
  BASE_MOVE_HISTORY_PATH: `${BASE_COUNSELING_MOVE_PATH}/history`,
  MOVE_HISTORY_PATH: 'history',
  BASE_MTO_PATH: `${BASE_COUNSELING_MOVE_PATH}/mto`,
  MTO_PATH: 'mto',
  BASE_CUSTOMER_SUPPORT_REMARKS_PATH: `${BASE_COUNSELING_MOVE_PATH}/customer-support-remarks`,
  CUSTOMER_SUPPORT_REMARKS_PATH: '/customer-support-remarks',
  BASE_SHIPMENT_REVIEW_PATH: `${BASE_COUNSELING_MOVE_PATH}/shipments/:shipmentId/document-review`,
  SHIPMENT_REVIEW_PATH: 'shipments/:shipmentId/document-review',
  BASE_REVIEW_SHIPMENT_WEIGHTS_PATH: `${BASE_COUNSELING_MOVE_PATH}/review-shipment-weights`,
  REVIEW_SHIPMENT_WEIGHTS_PATH: 'review-shipment-weights',
  CREATE_CUSTOMER_PATH: '/onboarding/create-customer',
};

const BASE_MOVES_PATH = '/moves/:moveCode';
export const moveRoutes = {
  BASE_MOVE_HISTORY_PATH: `${BASE_MOVES_PATH}/history`,
  MOVE_HISTORY_PATH: 'history',
};

export const tioRoutes = {
  BASE_PAYMENT_REQUESTS_PATH: `${BASE_MOVES_PATH}/payment-requests`,
  PAYMENT_REQUESTS_PATH: 'payment-requests',
  BILLABLE_WEIGHT_PATH: 'billable-weight',
};

export const tooRoutes = {
  BASE_SHIPMENT_EDIT_PATH: `${BASE_MOVES_PATH}/shipments/:shipmentId`,
  SHIPMENT_EDIT_PATH: 'shipments/:shipmentId',
  BASE_MOVE_VIEW_PATH: `${BASE_MOVES_PATH}/details`,
  MOVE_VIEW_PATH: 'details',
  BASE_ORDERS_EDIT_PATH: `${BASE_MOVES_PATH}/orders`,
  ORDERS_EDIT_PATH: 'orders',
  BASE_SHIPMENT_ADVANCE_PATH_TOO: `${BASE_MOVES_PATH}/shipments/:shipmentId/advance`,
};

export const qaeCSRRoutes = {
  MOVE_SEARCH_PATH: '/qaecsr/search',
  BASE_EVALUATION_REPORTS_PATH: `${BASE_MOVES_PATH}/evaluation-reports`,
  EVALUATION_REPORTS_PATH: '/evaluation-reports',
  BASE_EVALUATION_REPORT_PATH: `${BASE_MOVES_PATH}/evaluation-reports/:reportId`,
  EVALUATION_REPORT_PATH: '/evaluation-reports/:reportId',
  BASE_EVALUATION_VIOLATIONS_PATH: `${BASE_MOVES_PATH}/evaluation-reports/:reportId/violations`,
  EVALUATION_VIOLATIONS_PATH: '/evaluation-reports/:reportId/violations',
  BASE_CUSTOMER_SUPPORT_REMARKS_PATH: `${BASE_MOVES_PATH}/customer-support-remarks`,
  CUSTOMER_SUPPORT_REMARKS_PATH: 'customer-support-remarks',
};

// README: Prime API endpoints accept either a Move code or ID.
// The base path doesn't have prime-simulator due to matching issues with /prime.
const BASE_PRIME_SIMULATOR_PATH = '/simulator/moves/:moveCodeOrID';

export const primeSimulatorRoutes = {
  VIEW_MOVE_PATH: `${BASE_PRIME_SIMULATOR_PATH}/details`,
  CREATE_SHIPMENT_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/new`,
  UPDATE_SHIPMENT_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId`,
  CREATE_PAYMENT_REQUEST_PATH: `${BASE_PRIME_SIMULATOR_PATH}/payment-requests/new`,
  CREATE_SERVICE_ITEM_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/service-items/new`,
  UPDATE_SIT_SERVICE_ITEM_PATH: `${BASE_PRIME_SIMULATOR_PATH}/mto-service-items/:mtoServiceItemId/update`,
  UPLOAD_DOCUMENTS_PATH: `${BASE_PRIME_SIMULATOR_PATH}/payment-requests/:paymentRequestId/upload`,
  UPLOAD_SERVICE_REQUEST_DOCUMENTS_PATH: `${BASE_PRIME_SIMULATOR_PATH}/mto-service-items/:mtoServiceItemId/upload`,
  SHIPMENT_UPDATE_ADDRESS_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/addresses/update`,
  SHIPMENT_UPDATE_REWEIGH_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/reweigh/:reweighId/update`,
  CREATE_SIT_EXTENSION_REQUEST_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/sit-extension-requests/new`,
  SHIPMENT_UPDATE_DESTINATION_ADDRESS_PATH: `${BASE_PRIME_SIMULATOR_PATH}/shipments/:shipmentId/updateDestinationAddress`,
};

export const adminRoutes = {
  HOME_PATH: '/',
};
