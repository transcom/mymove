import { all } from 'redux-saga/effects';

import watchSwaggerRequests from './request';
import watchFetchUser from './auth';

export default function* rootSaga() {
  yield all([watchSwaggerRequests(), watchFetchUser()]);
}
