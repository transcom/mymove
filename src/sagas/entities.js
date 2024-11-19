import { all, takeLatest, put, call } from 'redux-saga/effects';

import {
  UPDATE_SERVICE_MEMBER,
  UPDATE_BACKUP_CONTACT,
  UPDATE_MOVE,
  UPDATE_MTO_SHIPMENT,
  UPDATE_MTO_SHIPMENTS,
  UPDATE_ORDERS,
  UPDATE_ALL_MOVES,
} from 'store/entities/actions';
import { normalizeResponse } from 'services/swaggerRequest';
import { addEntities, updateMTOShipmentsEntity, setOktaUser } from 'shared/Entities/actions';

export function* updateOktaUserState(action) {
  const { payload } = action;

  yield put(setOktaUser(payload));
}

export function* updateServiceMember(action) {
  const { payload } = action;

  const normalizedData = yield call(normalizeResponse, payload, 'serviceMember');
  yield put(addEntities(normalizedData));
}

export function* updateBackupContact(action) {
  const { payload } = action;

  const normalizedData = yield call(normalizeResponse, payload, 'backupContact');
  yield put(addEntities(normalizedData));
}

export function* updateOrders(action) {
  const { payload } = action;
  const normalizedData = yield call(normalizeResponse, payload, 'orders');
  yield put(addEntities(normalizedData));
}

export function* updateMove(action) {
  const { payload } = action;

  const normalizedData = yield call(normalizeResponse, payload, 'move');
  yield put(addEntities(normalizedData));
}

export function* updateMTOShipment(action) {
  const { payload } = action;
  const normalizedData = yield call(normalizeResponse, payload, 'mtoShipment');
  yield put(addEntities(normalizedData));
}

export function* updateMTOShipments(action) {
  const { payload } = action;

  yield put(updateMTOShipmentsEntity(payload));
}

export function* updateAllMoves(action) {
  const { payload } = action;

  yield put(addEntities({ serviceMemberMoves: payload }));
}

export function* unlockedMoves(action) {
  const { payload } = action;

  const normalizedData = yield call(normalizeResponse, payload, 'move');
  yield put(addEntities(normalizedData));
}

export function* watchUpdateEntities() {
  yield all([
    takeLatest(UPDATE_SERVICE_MEMBER, updateServiceMember),
    takeLatest(UPDATE_BACKUP_CONTACT, updateBackupContact),
    takeLatest(UPDATE_ORDERS, updateOrders),
    takeLatest(UPDATE_MOVE, updateMove),
    takeLatest(UPDATE_MTO_SHIPMENT, updateMTOShipment),
    takeLatest(UPDATE_MTO_SHIPMENTS, updateMTOShipments),
    takeLatest(UPDATE_ALL_MOVES, updateAllMoves),
  ]);
}
