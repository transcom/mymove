import { all } from 'redux-saga/effects';

import watchFetchUser from './auth';
import { watchInitializeOnboarding } from './onboarding';
import { watchUpdateEntities } from './entities';

export default function* rootSaga() {
  yield all([watchFetchUser()]);
}

export function* rootCustomerSaga() {
  yield all([watchFetchUser(), watchUpdateEntities(), watchInitializeOnboarding()]);
}
