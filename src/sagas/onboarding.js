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
import { selectServiceMemberFromLoggedInUser, selectServiceMemberProfileState } from 'store/entities/selectors';
import { profileStates } from 'constants/customerStates';

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

const findNextServiceMemberStep = (serviceMemberId, profileState) => {
  const profilePathPrefix = `/service-member/${serviceMemberId}`;

  switch (profileState) {
    case profileStates.EMPTY_PROFILE:
      return `${profilePathPrefix}/conus-status`;
    case profileStates.DOD_INFO_COMPLETE:
      return `${profilePathPrefix}/name`;
    case profileStates.NAME_COMPLETE:
      return `${profilePathPrefix}/contact-info`;
    case profileStates.CONTACT_INFO_COMPLETE:
      return `${profilePathPrefix}/duty-station`;
    case profileStates.DUTY_STATION_COMPLETE:
      return `${profilePathPrefix}/residence-address`;
    case profileStates.ADDRESS_COMPLETE:
      return `${profilePathPrefix}/backup-mailing-address`;
    case profileStates.BACKUP_ADDRESS_COMPLETE:
      return `${profilePathPrefix}/backup-contacts`;
    default:
      return '/';
  }
};

export function* initializeOnboarding() {
  try {
    const user = yield call(fetchCustomerData);
    if (!user.serviceMembers) {
      yield call(createServiceMember);
    }

    // Determine where user should be directed
    const serviceMember = yield select(selectServiceMemberFromLoggedInUser);
    const serviceMemberProfileState = yield select(selectServiceMemberProfileState);
    const nextPagePath = findNextServiceMemberStep(serviceMember.id, serviceMemberProfileState);

    yield put(push(nextPagePath));

    yield put(initOnboardingComplete());
    yield all([call(watchFetchCustomerData), call(watchUpdateServiceMember)]);
  } catch (error) {
    yield put(initOnboardingFailed(error));
  }
}

export function* watchInitializeOnboarding() {
  yield takeLatest(INIT_ONBOARDING, initializeOnboarding);
}
