import { takeLatest, put, call } from 'redux-saga/effects';

import { INIT_ONBOARDING, initOnboardingFailed, initOnboardingComplete } from 'store/onboarding/actions';
import { getLoggedInUser, getMTOShipmentsForMove } from 'services/internalApi';
import { addEntities } from 'shared/Entities/actions';

export function* initializeOnboarding() {
  try {
    // First load the user & store in entities
    const user = yield call(getLoggedInUser);
    yield put(addEntities(user));

    // TODO - create service member if doesn't exist

    // Load MTO shipments if there is a move
    const { moves } = user;
    if (moves && Object.keys(moves).length > 0) {
      const [moveId] = Object.keys(moves);
      // User has a move, load MTO shipments & store in entities
      const mtoShipments = yield call(getMTOShipmentsForMove, moveId);
      yield put(addEntities(mtoShipments));
    }

    yield put(initOnboardingComplete());
  } catch (error) {
    yield put(initOnboardingFailed(error));
  }
}

export function* watchInitializeOnboarding() {
  yield takeLatest(INIT_ONBOARDING, initializeOnboarding);
}
