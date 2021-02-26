import { takeLatest, put, call, all, select } from 'redux-saga/effects';
import { push } from 'connected-react-router';

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
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { NULL_UUID } from 'shared/constants';

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
  // TODO - delete legacy actions after service member reducer is deleted
  try {
    yield put({ type: CREATE_SERVICE_MEMBER.start });
    const serviceMember = yield call(createServiceMemberApi);
    yield put({ type: CREATE_SERVICE_MEMBER.success, payload: serviceMember });
    yield call(fetchCustomerData);
  } catch (e) {
    yield put({ type: CREATE_SERVICE_MEMBER.failure, error: e });
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

const findNextServiceMemberStep = (serviceMember) => {
  if (!serviceMember.rank || !serviceMember.edipi || !serviceMember.affiliation) return '/conus-status';

  if (!serviceMember.first_name || !serviceMember.last_name) return '/name';

  if (
    !serviceMember.telephone ||
    !serviceMember.personal_email ||
    !(serviceMember.phone_is_preferred || serviceMember.email_is_preferred)
  )
    return '/contact-info';

  if (!serviceMember.current_station || serviceMember.current_station.id === NULL_UUID) return '/duty-station';

  if (!serviceMember.residential_address) return '/residence-address';

  if (!serviceMember.backup_mailing_address) return '/backup-mailing-address';

  if (!serviceMember.backup_contacts || !serviceMember.backup_contacts.length) return '/backup-contacts';

  return '/';
};

export function* initializeOnboarding() {
  try {
    const user = yield call(fetchCustomerData);
    if (!user.serviceMembers) {
      yield call(createServiceMember);
    }

    // Determine where the user should be directed
    const serviceMember = yield select(selectServiceMemberFromLoggedInUser);

    // console.log('check SM state', serviceMember);

    const nextPageRoute = findNextServiceMemberStep(serviceMember);
    // console.log('redirect to', nextPageRoute);

    if (nextPageRoute === '/') {
      yield put(push('/'));
    } else {
      yield put(push(`/service-member/${serviceMember.id}${nextPageRoute}`));
    }

    yield put(initOnboardingComplete());

    // Watch for update actions
    yield all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]);
  } catch (error) {
    yield put(initOnboardingFailed(error));
  }
}

export function* watchInitializeOnboarding() {
  yield takeLatest(INIT_ONBOARDING, initializeOnboarding);
}
