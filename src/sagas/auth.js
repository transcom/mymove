import { takeLatest, put, call } from 'redux-saga/effects';
import { normalize } from 'normalizr';

import { LOAD_USER } from 'store/auth/actions';
import { GetIsLoggedIn, GetLoggedInUser } from 'shared/User/api';
import { ordersArray } from 'shared/Entities/schema';
import { addEntities } from 'shared/Entities/actions';
import { getLoggedInActions } from 'shared/Data/users';
import { showLoggedInUser } from 'shared/Entities/modules/user';

/**
 * This saga mirrors the getCurrentUserInfo thunk (shared/Data/users.js)
 * and is triggered by the 'LOAD_USER' action (currently only called by
 * the OfficeApp)
 */
export function* fetchUser() {
  yield put(getLoggedInActions.start());

  try {
    const isLoggedIn = yield call(GetIsLoggedIn);
    if (isLoggedIn) {
      try {
        // Fire API call to put user info in entities
        yield put(showLoggedInUser());

        // Legacy call to put user in user reducer
        const user = yield call(GetLoggedInUser);

        if (user.service_member) {
          const data = normalize(user.service_member.orders, ordersArray);
          const filtered = { ...data.entities.addresses };
          yield put(addEntities(filtered));
        }

        yield put(getLoggedInActions.success(user));
      } catch (e) {
        yield put(getLoggedInActions.error(e));
      }
    } else {
      yield put(getLoggedInActions.error('User is not logged in'));
    }
  } catch (e) {
    yield put(getLoggedInActions.error(e));
  }
}

export default function* watchFetchUser() {
  yield takeLatest(LOAD_USER, fetchUser);
}
