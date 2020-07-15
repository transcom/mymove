import { all } from 'redux-saga/effects';

import watchFetchUser from './auth';

export default function* rootSaga() {
  yield all([watchFetchUser()]);
}
