export const generalRoutes = {
  HOME: '/',
  SIGN_IN: '/sign-in',
  PRIVACY_SECURITY_POLICY: '/privacy-security',
  ACCESSIBILITY: '/accessibility',
};

export const customerRoutes = {
  ACCESS_CODE: '/access-code',
  CONUS_OCONUS: '/service-member/conus-oconus',
  DOD_INFO: '/service-member/dod-info',
  NAME: '/service-member/name',
  CONTACT_INFO: '/service-member/contact-info',
  CURRENT_DUTY_STATION: '/service-member/current-duty',
  CURRENT_ADDRESS: '/service-member/current-address',
  BACKUP_ADDRESS: '/service-member/backup-address',
  BACKUP_CONTACTS: '/service-member/backup-contact',
  ORDERS_INFO: '/orders/info',
  ORDERS_UPLOAD: '/orders/upload',
  SHIPMENT_MOVING_INFO: '/moves/:moveId/moving-info',
  SHIPMENT_SELECT_TYPE: '/moves/:moveId/shipment-type',
  SHIPMENT_EDIT: '/moves/:moveId/shipments/:mtoShipmentId/edit',
};
