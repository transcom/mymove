export const UPDATE_SERVICE_MEMBER = 'UPDATE_SERVICE_MEMBER';
export const UPDATE_BACKUP_CONTACT = 'UPDATE_BACKUP_CONTACT';
export const UPDATE_MOVE = 'UPDATE_MOVE';
export const UPDATE_MTO_SHIPMENT = 'UPDATE_MTO_SHIPMENT';
export const UPDATE_MTO_SHIPMENTS = 'UPDATE_MTO_SHIPMENTS';
export const UPDATE_ORDERS = 'UPDATE_ORDERS';
export const UPDATE_OKTA_USER_STATE = 'SET_OKTA_USER';
export const UPDATE_ACTIVE_ROLE = 'UPDATE_ACTIVE_ROLE';
export const UPDATE_ALL_MOVES = 'UPDATE_ALL_MOVES';

export const updateOktaUserState = (oktaUser) => ({
  type: UPDATE_OKTA_USER_STATE,
  oktaUser,
});

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

export const updateAllMoves = (payload) => ({
  type: UPDATE_ALL_MOVES,
  payload,
});
