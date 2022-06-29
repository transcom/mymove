import React, { useEffect, useRef, useState } from 'react';
import { useIdleTimer } from 'react-idle-timer';

import Alert from 'shared/Alert';

const maxIdleTimeInSeconds = 14 * 60;
const maxWarningTimeBeforeTimeoutInSeconds = 60;
const maxIdleTimeInMilliseconds = maxIdleTimeInSeconds * 1000;
const keepAliveEndpoint = '/internal/users/logged_in';

const LogoutOnInactivity = () => {
  const [isIdle, setIsIdle] = useState(false);
  const [timeLeftInSeconds, setTimeLeftInSeconds] = useState(maxWarningTimeBeforeTimeoutInSeconds);
  const timerRef = useRef(null);

  const handleOnActive = () => {
    setIsIdle(false);
    clearInterval(timerRef.current);
    setTimeLeftInSeconds(maxWarningTimeBeforeTimeoutInSeconds);
    fetch(keepAliveEndpoint);
  };

  const handleOnIdle = () => {
    setIsIdle(true);
    clearInterval(timerRef.current);
    timerRef.current = setInterval(() => {
      setTimeLeftInSeconds((current) => {
        return current - 1;
      });
    }, 1000);
  };

  useIdleTimer({
    element: document,
    timeout: maxIdleTimeInMilliseconds,
    onActive: handleOnActive,
    onIdle: handleOnIdle,
    events: ['blur', 'focus', 'mousedown', 'touchstart', 'MSPointerDown'],
  });

  useEffect(() => {
    if (isIdle) {
      if (timeLeftInSeconds <= 0) {
        clearInterval(timerRef.current);
        // log out
      }
    }
  }, [isIdle, timeLeftInSeconds]);

  return (
    <div data-testid="logoutOnInactivityWrapper">
      <h1>Time left in seconds: {timeLeftInSeconds}</h1>
      <h1>Idle: {isIdle.toString()}</h1>
      {isIdle && (
        <div data-testid="logoutAlert">
          <Alert data-testid="logoutAlert" type="warning" heading="Inactive user">
            You have been inactive and will be logged out in {timeLeftInSeconds} seconds unless you touch or click on
            the page.
          </Alert>
        </div>
      )}
    </div>
  );
};

export default LogoutOnInactivity;
