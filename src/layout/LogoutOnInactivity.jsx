import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { useHistory } from 'react-router-dom';
import { useIdleTimer } from 'react-idle-timer';
import { useSelector, useDispatch } from 'react-redux';
import { push } from 'connected-react-router';

import { logOut } from 'store/auth/actions';
import { selectIsLoggedIn } from 'store/auth/selectors';
import Alert from 'shared/Alert';
import { LogoutUser } from 'utils/api';

const defaultMaxIdleTimeInSeconds = 14 * 60;
const defaultMaxWarningTimeInSeconds = 60;
const keepAliveEndpoint = '/internal/users/logged_in';

const LogoutOnInactivity = ({ maxIdleTimeInSeconds, maxWarningTimeInSeconds }) => {
  const [isIdle, setIsIdle] = useState(false);
  const [timeLeftInSeconds, setTimeLeftInSeconds] = useState(maxWarningTimeInSeconds);
  const history = useHistory();
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const dispatch = useDispatch();

  const handleOnActive = () => {
    setIsIdle(false);
    setTimeLeftInSeconds(maxWarningTimeInSeconds);
    if (isLoggedIn) {
      fetch(keepAliveEndpoint);
    }
  };

  const handleOnIdle = () => {
    setIsIdle(true);
  };

  useIdleTimer({
    element: document,
    timeout: maxIdleTimeInSeconds * 1000,
    onActive: handleOnActive,
    onIdle: handleOnIdle,
    events: ['blur', 'focus', 'mousedown', 'touchstart', 'MSPointerDown'],
    startOnMount: true,
  });

  useEffect(() => {
    let warningTimer;
    if (isIdle && isLoggedIn) {
      let timeLeft = maxWarningTimeInSeconds;
      warningTimer = setInterval(() => {
        setTimeLeftInSeconds((current) => current - 1);
        if (timeLeft > 1) {
          timeLeft -= 1;
        } else {
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
      }, 1000);
    }
    return () => clearInterval(warningTimer);
  }, [isIdle, isLoggedIn, history, maxWarningTimeInSeconds, dispatch]);

  return (
    isLoggedIn && (
      <div data-testid="logoutOnInactivityWrapper">
        {isIdle && (
          <div data-testid="logoutAlert">
            <Alert type="warning" heading="Inactive user">
              You have been inactive and will be logged out in {timeLeftInSeconds} seconds unless you touch or click on
              the page.
            </Alert>
          </div>
        )}
      </div>
    )
  );
};

LogoutOnInactivity.defaultProps = {
  maxIdleTimeInSeconds: defaultMaxIdleTimeInSeconds,
  maxWarningTimeInSeconds: defaultMaxWarningTimeInSeconds,
};

LogoutOnInactivity.propTypes = {
  maxIdleTimeInSeconds: PropTypes.number,
  maxWarningTimeInSeconds: PropTypes.number,
};

export default LogoutOnInactivity;
