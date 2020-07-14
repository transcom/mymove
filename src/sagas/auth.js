import { takeLatest, put, call } from 'redux-saga/effects';
import { normalize } from 'normalizr';

import { generateAsyncActions } from 'shared/ReduxHelpers';
import { GetIsLoggedIn, GetLoggedInUser } from 'shared/User/api';
import { ordersArray } from 'shared/Entities/schema';
import { addEntities } from 'shared/Entities/actions';

const getLoggedInUserType = 'GET_LOGGED_IN_USER';

export const getLoggedInActions = generateAsyncActions(getLoggedInUserType);

export function* fetchUser() {
  try {
    const isLoggedIn = yield call(GetIsLoggedIn);
    if (isLoggedIn) {
      try {
        const user = yield call(GetLoggedInUser);
        if (user.service_member) {
          const data = normalize(user.service_member.orders, ordersArray);
          const filtered = { ...data.entities.addresses };
          yield put(addEntities(filtered));
        }

        yield put(getLoggedInActions.success(user));
      } catch (e) {
        put(getLoggedInActions.error(e));
      }
    } else {
      put(getLoggedInActions.error('User is not logged in'));
    }
  } catch (e) {
    put(getLoggedInActions.error(e));
  }
}

export default function* watchFetchUser() {
  yield takeLatest('GET_LOGGED_IN_USER_START', fetchUser);
}
