import React, { useEffect, useRef, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { useIdleTimer } from 'react-idle-timer';
import { useSelector, useDispatch } from 'react-redux';
import { push } from 'connected-react-router';

import { logOut } from 'store/auth/actions';
import { selectIsLoggedIn } from 'store/auth/selectors';
import Alert from 'shared/Alert';
import { LogoutUser } from 'utils/api';

const maxIdleTimeInSeconds = 14 * 60;
const maxWarningTimeBeforeTimeoutInSeconds = 60;
const maxIdleTimeInMilliseconds = maxIdleTimeInSeconds * 1000;
const keepAliveEndpoint = '/internal/users/logged_in';

const LogoutOnInactivity = () => {
  const [isIdle, setIsIdle] = useState(false);
  const [timeLeftInSeconds, setTimeLeftInSeconds] = useState(maxWarningTimeBeforeTimeoutInSeconds);
  const timerRef = useRef(null);
  const history = useHistory();
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const dispatch = useDispatch();

  const handleOnActive = () => {
    if (isLoggedIn) {
      setIsIdle(false);
      clearInterval(timerRef.current);
      setTimeLeftInSeconds(maxWarningTimeBeforeTimeoutInSeconds);
      fetch(keepAliveEndpoint);
    }
  };

  const handleOnIdle = () => {
    if (isLoggedIn) {
      setIsIdle(true);
      clearInterval(timerRef.current);
      timerRef.current = setInterval(() => {
        setTimeLeftInSeconds((current) => {
          return current - 1;
        });
      }, 1000);
    }
  };

  useIdleTimer({
    element: document,
    timeout: maxIdleTimeInMilliseconds,
    onActive: handleOnActive,
    onIdle: handleOnIdle,
    events: ['blur', 'focus', 'mousedown', 'touchstart', 'MSPointerDown'],
    startOnMount: true,
  });

  useEffect(() => {
    if (isIdle) {
      if (isLoggedIn && timeLeftInSeconds <= 0) {
        clearInterval(timerRef.current);
        dispatch(push(logOut));
        LogoutUser().then((r) => {
          const redirectURL = r.body;
          if (redirectURL) {
            window.location.href = redirectURL;
          } else {
            history.push({
              pathname: '/sign-in',
              state: { hasLoggedOut: true },
            });
          }
        });
      }
    }
  }, [isIdle, isLoggedIn, timeLeftInSeconds, history, dispatch]);

  return (
    <div data-testid="logoutOnInactivityWrapper">
      {isLoggedIn && (
        <>
          {/* testing - to remove */}
          <h1>Time left in seconds: {timeLeftInSeconds}</h1>
          <h1>Idle: {isIdle.toString()}</h1>
          <h1>Logged in: {isLoggedIn.toString()}</h1>
          {isIdle && (
            <div data-testid="logoutAlert">
              <Alert data-testid="logoutAlert" type="warning" heading="Inactive user">
                You have been inactive and will be logged out in {timeLeftInSeconds} seconds unless you touch or click
                on the page.
              </Alert>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default LogoutOnInactivity;
