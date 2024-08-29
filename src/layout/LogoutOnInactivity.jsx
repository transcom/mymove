import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { useNavigate } from 'react-router-dom';
import { useIdleTimer } from 'react-idle-timer';
import { useSelector } from 'react-redux';

import { selectIsLoggedIn } from 'store/auth/selectors';
import Alert from 'shared/Alert';
import { LogoutUser } from 'utils/api';

const defaultIdleTimeout = 15_000 * 60;
const defaultWarningTime = 1_000 * 60;
const keepAliveEndpoint = '/internal/users/logged_in';

/**
 * @description The component that handles logging out inactive users.
 * @param {int} idleTimeout the amount of time in milliseconds that the user can be idle before being logged out
 * @param {int} warningTime the amount of time in milliseconds that the user will be shown a warning before being logged out
 * @return {JSX.Element}
 * */
const LogoutOnInactivity = ({ idleTimeout, warningTime }) => {
  const navigate = useNavigate();
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const [showLogoutWarning, setShowLogoutWarning] = useState(false);
  const [remaining, setRemaining] = useState(0);

  const onPrompt = () => {
    setShowLogoutWarning(true);
  };

  const onActive = () => {
    setShowLogoutWarning(false);
    if (isLoggedIn) {
      fetch(keepAliveEndpoint);
    }
  };

  const onAction = (_event, idleTimer) => {
    idleTimer.activate();
  };

  const onIdle = () => {
    if (isLoggedIn) {
      // explicit clear session storage
      window.sessionStorage.clear();
      LogoutUser().then((r) => {
        const redirectURL = r.body;
        // checking to see if "Local Sign In" button was used to sign in
        const urlParams = new URLSearchParams(redirectURL.split('?')[1]);
        const idTokenHint = urlParams.get('id_token_hint');
        if (redirectURL && idTokenHint !== 'devlocal') {
          window.location.href = redirectURL;
        } else {
          navigate('/sign-in', { state: { hasLoggedOut: true } });
        }
      });
    }
  };

  const { getRemainingTime } = useIdleTimer({
    element: document,
    events: ['blur', 'focus', 'mousedown', 'touchstart', 'MSPointerDown'],
    onAction,
    onActive,
    onIdle,
    onPrompt,
    promptBeforeIdle: warningTime,
    startOnMount: true,
    timeout: idleTimeout,
  });

  useEffect(() => {
    const interval = setInterval(() => {
      setRemaining(Math.floor(getRemainingTime() / 1000));
    }, 500);

    return () => {
      clearInterval(interval);
    };
  });

  return (
    isLoggedIn && (
      <div data-testid="logoutOnInactivityWrapper">
        {showLogoutWarning && (
          <div data-testid="logoutAlert">
            <Alert type="warning" heading="Inactive user">
              You have been inactive and will be logged out in {remaining} seconds unless you touch or click on the
              page.
            </Alert>
          </div>
        )}
      </div>
    )
  );
};

LogoutOnInactivity.defaultProps = {
  idleTimeout: defaultIdleTimeout,
  warningTime: defaultWarningTime,
};

LogoutOnInactivity.propTypes = {
  idleTimeout: PropTypes.number,
  warningTime: PropTypes.number,
};

export default LogoutOnInactivity;
