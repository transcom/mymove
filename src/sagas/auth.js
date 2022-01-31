import { takeLatest, put, call } from 'redux-saga/effects';
import { normalize } from 'normalizr';

import { LOAD_USER, getLoggedInUserStart, getLoggedInUserSuccess, getLoggedInUserFailure } from 'store/auth/actions';
import { setFlashMessage } from 'store/flash/actions';
import { GetIsLoggedIn, GetLoggedInUser } from 'utils/api';
import { loggedInUser } from 'shared/Entities/schema';
import { addEntities } from 'shared/Entities/actions';

/**
 * This saga mirrors the getCurrentUserInfo thunk (shared/Data/users.js)
 * and is triggered by the 'LOAD_USER' action
 */
export function* fetchUser() {
  yield put(getLoggedInUserStart());

  try {
    // The `GetIsLoggedIn` call returns a object with a parameter isLoggedIn
    const { isLoggedIn } = yield call(GetIsLoggedIn);
    if (isLoggedIn) {
      try {
        const user = yield call(GetLoggedInUser); // make user API call

        const userEntities = normalize(user, loggedInUser);

        yield put(addEntities(userEntities.entities)); // populate entities
        yield put(getLoggedInUserSuccess(user));
      } catch (e) {
        yield put(
          setFlashMessage(
            'USER_GET_ERROR',
            'error',
            'There was an error loading your user information.',
            'An error occurred',
          ),
        );
        yield put(getLoggedInUserFailure(e));
      }
    } else {
      // No flash message here - in this case the user should be shown the Log In screen
      yield put(getLoggedInUserFailure('User is not logged in'));
    }
  } catch (e) {
    yield put(
      setFlashMessage(
        'LOGGED_IN_GET_ERROR',
        'error',
        'There was an error loading your user information.',
        'An error occurred',
      ),
    );
    yield put(getLoggedInUserFailure(e));
  }
}

export default function* watchFetchUser() {
  yield takeLatest(LOAD_USER, fetchUser);
}
