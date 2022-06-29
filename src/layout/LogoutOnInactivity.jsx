import React, { useState } from 'react';
import { useIdleTimer } from 'react-idle-timer';

import Alert from 'shared/Alert';

const maxIdleTimeInSeconds = 14 * 60;
// const maxWarningTimeBeforeTimeoutInSeconds = 60;
const maxIdleTimeInMilliseconds = maxIdleTimeInSeconds * 1000;
const keepAliveEndpoint = '/internal/users/logged_in';

const LogoutOnInactivity = () => {
  const [isIdle, setIsIdle] = useState(false);

  const handleOnActive = () => {
    setIsIdle(false);
    fetch(keepAliveEndpoint);
  };

  const handleOnIdle = () => {
    setIsIdle(true);
  };

  useIdleTimer({
    element: document,
    timeout: maxIdleTimeInMilliseconds,
    onActive: handleOnActive,
    onIdle: handleOnIdle,
    events: ['blur', 'focus', 'mousedown', 'touchstart', 'MSPointerDown'],
  });

  return (
    <div data-testid="logoutOnInactivityWrapper">
      <h1>Time left in seconds: {}</h1>
      <h1>Idle: {isIdle.toString()}</h1>
      {isIdle && (
        <div data-testid="logoutAlert">
          <Alert data-testid="logoutAlert" type="warning" heading="Inactive user">
            You have been inactive and will be logged out in {} seconds unless you touch or click on the page.
          </Alert>
        </div>
      )}
    </div>
  );
};

export default LogoutOnInactivity;
