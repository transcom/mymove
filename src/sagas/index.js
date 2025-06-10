import { all } from 'redux-saga/effects';

import { watchFetchUser, watchHandleSetActiveRole } from './auth';
import { watchInitializeOnboarding } from './onboarding';
import { watchUpdateEntities } from './entities';

export default function* rootSaga() {
  yield all([watchFetchUser(), watchHandleSetActiveRole()]);
}

export function* rootCustomerSaga() {
  yield all([watchFetchUser(), watchUpdateEntities(), watchInitializeOnboarding()]);
}
