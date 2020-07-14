import { takeEvery } from 'redux-saga/effects';

export function* startRequest(action) {
  // eslint-disable-next-line no-console
  console.log('start request', action);
  yield;
}

export default function* watchSwaggerRequests() {
  yield takeEvery((action) => action.type.indexOf('@@swagger') === 0, startRequest);
}
