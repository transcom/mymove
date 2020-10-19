import { all, takeLatest, put } from 'redux-saga/effects';

import { UPDATE_SERVICE_MEMBER } from 'store/entities/actions';
import { normalizeResponse } from 'services/swaggerRequest';
import { addEntities } from 'shared/Entities/actions';

export function* updateServiceMember(action) {
  const { payload } = action;

  const normalizedData = normalizeResponse(payload, 'serviceMember');
  yield put(addEntities(normalizedData));
  yield put({
    type: 'UPDATE_SERVICE_MEMBER_SUCCESS',
    payload,
  });
}

export function* watchUpdateEntities() {
  yield all([takeLatest(UPDATE_SERVICE_MEMBER, updateServiceMember)]);
}
