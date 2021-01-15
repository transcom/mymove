import { all, takeLatest, put, call } from 'redux-saga/effects';

import {
  UPDATE_SERVICE_MEMBER,
  UPDATE_BACKUP_CONTACT,
  UPDATE_MOVE,
  UPDATE_MTO_SHIPMENT,
  UPDATE_ORDERS,
  UPDATE_PPMS,
  UPDATE_PPM,
  UPDATE_PPM_ESTIMATE,
  UPDATE_PPM_SIT_ESTIMATE,
} from 'store/entities/actions';
import { normalizeResponse } from 'services/swaggerRequest';
import { addEntities } from 'shared/Entities/actions';

export function* updateServiceMember(action) {
  const { payload } = action;

  const normalizedData = yield call(normalizeResponse, payload, 'serviceMember');
  yield put(addEntities(normalizedData));
  yield put({
    type: 'UPDATE_SERVICE_MEMBER_SUCCESS',
    payload,
  });
}

export function* updateBackupContact(action) {
  const { payload } = action;

  const normalizedData = yield call(normalizeResponse, payload, 'backupContact');
  yield put(addEntities(normalizedData));
  yield put({
    type: 'UPDATE_BACKUP_CONTACT_SUCCESS',
    payload,
  });
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
  yield put({
    type: 'CREATE_OR_UPDATE_MOVE_SUCCESS',
    payload,
  });
}

export function* updateMTOShipment(action) {
  const { payload } = action;
  const normalizedData = yield call(normalizeResponse, payload, 'mtoShipment');
  yield put(addEntities(normalizedData));
}

export function* updatePPMs(action) {
  const { payload } = action;
  const normalizedData = yield call(normalizeResponse, payload, 'personallyProcuredMoves');
  yield put(addEntities(normalizedData));
}

export function* updatePPM(action) {
  const { payload } = action;
  const normalizedData = yield call(normalizeResponse, payload, 'personallyProcuredMove');
  yield put(addEntities(normalizedData));
  yield put({
    type: 'CREATE_OR_UPDATE_PPM_SUCCESS',
    payload,
  });
}

export function* updatePPMEstimate(action) {
  const { payload } = action;
  const normalizedData = yield call(normalizeResponse, payload, 'ppmEstimateRange');
  yield put(addEntities(normalizedData));
}

export function* updatePPMSitEstimate(action) {
  const { payload } = action;
  const normalizedData = yield call(normalizeResponse, payload, 'ppmSitEstimate');
  yield put(addEntities(normalizedData));
}

export function* watchUpdateEntities() {
  yield all([
    takeLatest(UPDATE_SERVICE_MEMBER, updateServiceMember),
    takeLatest(UPDATE_BACKUP_CONTACT, updateBackupContact),
    takeLatest(UPDATE_ORDERS, updateOrders),
    takeLatest(UPDATE_MOVE, updateMove),
    takeLatest(UPDATE_MTO_SHIPMENT, updateMTOShipment),
    takeLatest(UPDATE_PPMS, updatePPMs),
    takeLatest(UPDATE_PPM, updatePPM),
    takeLatest(UPDATE_PPM_ESTIMATE, updatePPMEstimate),
    takeLatest(UPDATE_PPM_SIT_ESTIMATE, updatePPMSitEstimate),
  ]);
}
