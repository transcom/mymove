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
  getAllMoves,
} from 'services/internalApi';
import { addEntities } from 'shared/Entities/actions';

export function* fetchCustomerData() {
  // First load the user & store in entities
  const user = yield call(getLoggedInUser);
  yield put(addEntities(user));

  // Load MTO shipments if there is a move
  const { moves } = user;
  if (moves && Object.keys(moves).length > 0) {
    const [moveId] = Object.keys(moves);
    // User has a move, load MTO shipments & store in entities
    const mtoShipments = yield call(getMTOShipmentsForMove, moveId);
    yield put(addEntities(mtoShipments));
  }

  // loading serviceMemberMoves for the user
  const { serviceMembers } = user;
  const key = Object.keys(serviceMembers)[0];
  const serviceMemberId = serviceMembers[key].id;
  const allMoves = yield call(getAllMoves, serviceMemberId);
  yield put(addEntities({ serviceMemberMoves: allMoves }));

  return user;
}

export function* watchFetchCustomerData() {
  yield takeLatest(FETCH_CUSTOMER_DATA, fetchCustomerData);
}

export function* createServiceMember() {
  try {
    yield call(createServiceMemberApi);
    yield call(fetchCustomerData);
  } catch (e) {
    yield put(
      setFlashMessage(
        'SERVICE_MEMBER_CREATE_ERROR',
        'error',
        'There was an error creating your profile information.',
        'An error occurred',
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
    yield all([call(watchFetchCustomerData)]);
  } catch (error) {
    yield put(initOnboardingFailed(error));
  }
}

export function* watchInitializeOnboarding() {
  yield takeLatest(INIT_ONBOARDING, initializeOnboarding);
}
