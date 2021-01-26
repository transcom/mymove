import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import IdleTimer from 'react-idle-timer';

import Alert from 'shared/Alert';
import { selectCurrentUser } from 'shared/Data/users';
import { LogoutUser } from 'shared/User/api.js';

const maxIdleTimeInSeconds = 15 * 60;
const maxWarningTimeBeforeTimeoutInSeconds = 60;
const maxIdleTimeInMilliseconds = maxIdleTimeInSeconds * 1000;
const maxWarningTimeBeforeTimeoutInMilliseconds = maxWarningTimeBeforeTimeoutInSeconds * 1000;
const timeToDisplayWarningInMilliseconds = maxIdleTimeInMilliseconds - maxWarningTimeBeforeTimeoutInMilliseconds;

export class LogoutOnInactivity extends React.Component {
  state = {
    isIdle: false,
    showLoggedOutAlert: false,
    timeLeftInSeconds: maxWarningTimeBeforeTimeoutInSeconds,
  };

  componentDidUpdate(prevProps) {
    if (!this.props.isLoggedIn && prevProps.isLoggedIn) {
      this.setState({ showLoggedOutAlert: true });
    }
  }

  onActive = () => {
    clearInterval(this.timer);
    this.setState({ isIdle: false });
    this.setState({ timeLeftInSeconds: maxWarningTimeBeforeTimeoutInSeconds });
  };

  onIdle = () => {
    this.setState({ isIdle: true });
    clearInterval(this.timer);
    this.timer = setInterval(this.countdown, 1000);
  };

  onAction = () => {
    fetch(this.props.keepAliveEndpoint);
  };

  countdown = () => {
    if (this.state.timeLeftInSeconds === 0) {
      LogoutUser();
    } else {
      this.setState({ timeLeftInSeconds: this.state.timeLeftInSeconds - 1 });
    }
  };

  render() {
    const props = this.props;
    return (
      <React.Fragment>
        {props.isLoggedIn && (
          <IdleTimer
            ref="idleTimer"
            element={document}
            onActive={this.onActive}
            onAction={this.onAction}
            onIdle={this.onIdle}
            timeout={this.props.warningTimeout}
            events={['keydown', 'mousedown', 'touchstart', 'MSPointerDown']}
          >
            {this.state.isIdle && (
              <Alert type="warning" heading="Inactive user">
                You have been inactive and will be logged out in {this.state.timeLeftInSeconds} seconds unless you touch
                or click on the page.
              </Alert>
            )}
          </IdleTimer>
        )}

        {this.state.showLoggedOutAlert && (
          <Alert type="error" heading="Logged out">
            You have been logged out due to inactivity.
          </Alert>
        )}
      </React.Fragment>
    );
  }
}
LogoutOnInactivity.defaultProps = {
  warningTimeout: timeToDisplayWarningInMilliseconds,
  timeRemaining: maxWarningTimeBeforeTimeoutInMilliseconds,
  keepAliveEndpoint: '/internal/swagger.yaml',
};
LogoutOnInactivity.propTypes = {
  warningTimeout: PropTypes.number.isRequired,
  timeRemaining: PropTypes.number.isRequired,
  keepAliveEndpoint: PropTypes.string.isRequired,
};

const mapStateToProps = (state) => {
  const user = selectCurrentUser(state);
  return {
    isLoggedIn: user.isLoggedIn,
  };
};

export default connect(mapStateToProps)(LogoutOnInactivity);
