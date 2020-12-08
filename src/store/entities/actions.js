export const UPDATE_SERVICE_MEMBER = 'UPDATE_SERVICE_MEMBER';
export const UPDATE_BACKUP_CONTACT = 'UPDATE_BACKUP_CONTACT';
export const UPDATE_MOVE = 'UPDATE_MOVE';
export const UPDATE_MTO_SHIPMENT = 'UPDATE_MTO_SHIPMENT';

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
