import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import IdleTimer from 'react-idle-timer';

import { isProduction } from 'shared/constants';
import Alert from 'shared/Alert';
import { selectCurrentUser } from 'shared/Data/users';
import { LogoutUser } from 'shared/User/api.js';

const fifteenMinutesInMilliseconds = 900000;
const tenMinutesInMilliseconds = 600000;
const oneMinuteInMilliseconds = 60000;
export class LogoutOnInactivity extends React.Component {
  state = {
    isIdle: false,
    showLoggedOutAlert: false,
  };
  componentDidMount() {
    this.interval = setInterval(() => fetch(this.props.keepAliveEndpoint), this.props.keepAliveInterval);
  }
  componentWillUnmount() {
    clearInterval(this.interval);
    if (this.timeout) clearTimeout(this.timeout);
  }
  componentDidUpdate(prevProps) {
    if (!this.props.isLoggedIn && prevProps.isLoggedIn) {
      this.setState({ showLoggedOutAlert: true });
    }
  }
  onActive = () => {
    this.setState({ isIdle: false });
  };
  onIdle = () => {
    this.setState({ isIdle: true });
    this.timeout = setTimeout(() => {
      if (this.state.isIdle) {
        LogoutUser();
      }
    }, this.props.logoutAfterWarningTimeout);
  };
  render() {
    const props = this.props;
    return (
      <React.Fragment>
        {isProduction &&
          props.isLoggedIn && (
            <IdleTimer
              ref="idleTimer"
              element={document}
              activeAction={this.onActive}
              idleAction={this.onIdle}
              timeout={this.props.idleTimeout}
            >
              {this.state.isIdle && (
                <Alert type="warning" heading="Inactive user">
                  You have been inactive and will be logged out shortly.
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
  idleTimeout: fifteenMinutesInMilliseconds,
  keepAliveInterval: tenMinutesInMilliseconds,
  logoutAfterWarningTimeout: oneMinuteInMilliseconds,
  keepAliveEndpoint: '/internal/swagger.yaml',
};
LogoutOnInactivity.propTypes = {
  idleTimeout: PropTypes.number.isRequired,
  keepAliveInterval: PropTypes.number.isRequired,
  logoutAfterWarningTimeout: PropTypes.number.isRequired,
  keepAliveEndpoint: PropTypes.string.isRequired,
};

const mapStateToProps = state => {
  const user = selectCurrentUser(state);
  return {
    isLoggedIn: user.isLoggedIn,
  };
};
export default connect(mapStateToProps)(LogoutOnInactivity);
