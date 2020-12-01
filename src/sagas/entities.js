import { all, takeLatest, put, call } from 'redux-saga/effects';

import { UPDATE_SERVICE_MEMBER, UPDATE_BACKUP_CONTACT } from 'store/entities/actions';
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

export function* watchUpdateEntities() {
  yield all([
    takeLatest(UPDATE_SERVICE_MEMBER, updateServiceMember),
    takeLatest(UPDATE_BACKUP_CONTACT, updateBackupContact),
  ]);
}
