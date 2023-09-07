export const UPDATE_SERVICE_MEMBER = 'UPDATE_SERVICE_MEMBER';
export const UPDATE_BACKUP_CONTACT = 'UPDATE_BACKUP_CONTACT';
export const UPDATE_MOVE = 'UPDATE_MOVE';
export const UPDATE_MTO_SHIPMENT = 'UPDATE_MTO_SHIPMENT';
export const UPDATE_MTO_SHIPMENTS = 'UPDATE_MTO_SHIPMENTS';
export const UPDATE_ORDERS = 'UPDATE_ORDERS';
export const UPDATE_PPMS = 'UPDATE_PPMS';
export const UPDATE_PPM = 'UPDATE_PPM';
export const UPDATE_PPM_ESTIMATE = 'UPDATE_PPM_ESTIMATE';
export const UPDATE_PPM_SIT_ESTIMATE = 'UPDATE_PPM_SIT_ESTIMATE';

export const updateServiceMember = (payload) => ({
  type: UPDATE_SERVICE_MEMBER,
  payload,
});

export const updateBackupContact = (payload) => ({
  type: UPDATE_BACKUP_CONTACT,
  payload,
});

export const updateMove = (payload) => ({
  type: UPDATE_MOVE,
  payload,
});

export const updateMTOShipment = (payload) => ({
  type: UPDATE_MTO_SHIPMENT,
  payload,
});

export const updateMTOShipments = (payload) => ({
  type: UPDATE_MTO_SHIPMENTS,
  payload,
});

export const updateOrders = (payload) => ({
  type: UPDATE_ORDERS,
  payload,
});

export const updatePPMs = (payload) => ({
  type: UPDATE_PPMS,
  payload,
});

export const updatePPM = (payload) => ({
  type: UPDATE_PPM,
  payload,
});

export const updatePPMEstimate = (payload) => ({
  type: UPDATE_PPM_ESTIMATE,
  payload,
});

export const updatePPMSitEstimate = (payload) => ({
  type: UPDATE_PPM_SIT_ESTIMATE,
  payload,
});
