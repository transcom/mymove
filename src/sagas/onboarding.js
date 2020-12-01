import { takeLatest, put, call, all } from 'redux-saga/effects';

import {
  INIT_ONBOARDING,
  FETCH_CUSTOMER_DATA,
  initOnboardingFailed,
  initOnboardingComplete,
} from 'store/onboarding/actions';
import { setFlashMessage } from 'store/flash/actions';
import {
  getLoggedInUser,
  getMTOShipmentsForMove,
  createServiceMember as createServiceMemberApi,
} from 'services/internalApi';
import { addEntities } from 'shared/Entities/actions';
import { CREATE_SERVICE_MEMBER } from 'scenes/ServiceMembers/ducks';
import { normalizeResponse } from 'services/swaggerRequest';

export function* fetchCustomerData() {
  // First load the user & store in entities
  const user = yield call(getLoggedInUser);
  yield put(addEntities(user));

  // TODO - fork/spawn additional API calls
  // Load MTO shipments if there is a move
  const { moves } = user;
  if (moves && Object.keys(moves).length > 0) {
    const [moveId] = Object.keys(moves);
    // User has a move, load MTO shipments & store in entities
    const mtoShipments = yield call(getMTOShipmentsForMove, moveId);
    yield put(addEntities(mtoShipments));
  }

  return user;
}

export function* watchFetchCustomerData() {
  yield takeLatest(FETCH_CUSTOMER_DATA, fetchCustomerData);
}

export function* updateServiceMember(action) {
  const { payload } = action;
  const normalizedData = normalizeResponse(payload, 'serviceMember');
  yield put(addEntities(normalizedData));
}

// legacy action - delete after SM entity refactor is complete
export function* watchUpdateServiceMember() {
  yield takeLatest('UPDATE_SERVICE_MEMBER_SUCCESS', updateServiceMember);
}

export function* createServiceMember() {
  try {
    yield put({ type: CREATE_SERVICE_MEMBER.start });
    const serviceMember = yield call(createServiceMemberApi);
    yield put({ type: CREATE_SERVICE_MEMBER.success, payload: serviceMember });
    yield call(fetchCustomerData);
  } catch (e) {
    yield put({ type: CREATE_SERVICE_MEMBER.failure, error: e });
    yield put(
      setFlashMessage(
        'error',
        'There was an error creating your profile information.',
        'An error occurred',
        'SERVICE_MEMBER_CREATE_ERROR',
      ),
    );
  }
}

export function* initializeOnboarding() {
  try {
    const user = yield call(fetchCustomerData);
    if (!user.serviceMembers) {
      yield call(createServiceMember);
    }
    yield put(initOnboardingComplete());
    yield all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]);
  } catch (error) {
    yield put(initOnboardingFailed(error));
  }
}

export function* watchInitializeOnboarding() {
  yield takeLatest(INIT_ONBOARDING, initializeOnboarding);
}
