import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { useNavigate } from 'react-router-dom';
import { useIdleTimer } from 'react-idle-timer';
import { useSelector } from 'react-redux';

import { selectIsLoggedIn } from 'store/auth/selectors';
import Alert from 'shared/Alert';
import { LogoutUser } from 'utils/api';

const defaultMaxIdleTime = 15_000 * 60;
const defaultWarningTime = 1_000 * 60;
const keepAliveEndpoint = '/internal/users/logged_in';

/**
 * @description The component that handles logging out inactive users.
 * @param {int} maxIdleTime the amount of time in milliseconds that the user can be idle before being logged out
 * @param {int} warningTime the amount of time in milliseconds that the user will be shown a warning before being logged out
 * @return {JSX.Element}
 * */
const LogoutOnInactivity = ({ maxIdleTime, warningTime }) => {
  const navigate = useNavigate();
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const [showLogoutWarning, setShowLogoutWarning] = useState(false);
  const [remaining, setRemaining] = useState(0);

  const handleOnPrompt = () => {
    setShowLogoutWarning(true);
  };

  const handleOnActive = () => {
    setShowLogoutWarning(false);
    if (isLoggedIn) {
      fetch(keepAliveEndpoint);
    }
  };

  const handleOnIdle = () => {
    if (isLoggedIn) {
      LogoutUser().then((r) => {
        const redirectURL = r.body;
        if (redirectURL) {
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
    onActive: handleOnActive,
    onIdle: handleOnIdle,
    onPrompt: handleOnPrompt,
    promptBeforeIdle: warningTime,
    startOnMount: true,
    timeout: maxIdleTime,
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
        {isLoggedIn && showLogoutWarning && (
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
  maxIdleTime: defaultMaxIdleTime,
  warningTime: defaultWarningTime,
};

LogoutOnInactivity.propTypes = {
  maxIdleTime: PropTypes.number,
  warningTime: PropTypes.number,
};

export default LogoutOnInactivity;
