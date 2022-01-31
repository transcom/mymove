import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import IdleTimer from 'react-idle-timer';
import { withRouter } from 'react-router-dom';

import { logOut as logOutAction } from 'store/auth/actions';
import { selectIsLoggedIn } from 'store/auth/selectors';
import Alert from 'shared/Alert';
import { LogoutUser } from 'utils/api';
import { HistoryShape } from 'types/customerShapes';

const maxIdleTimeInSeconds = 15 * 60;
const maxWarningTimeBeforeTimeoutInSeconds = 60;
const maxIdleTimeInMilliseconds = maxIdleTimeInSeconds * 1000;
const maxWarningTimeBeforeTimeoutInMilliseconds = maxWarningTimeBeforeTimeoutInSeconds * 1000;
const timeToDisplayWarningInMilliseconds = maxIdleTimeInMilliseconds - maxWarningTimeBeforeTimeoutInMilliseconds;
const keepAliveEndpoint = '/internal/users/logged_in';

export class LogoutOnInactivity extends React.Component {
  constructor(props) {
    super(props);

    this.onActive = this.onActive.bind(this);
    this.onIdle = this.onIdle.bind(this);

    this.state = {
      isIdle: false,
      timeLeftInSeconds: maxWarningTimeBeforeTimeoutInSeconds,
      loggedIn: true,
    };
  }

  onActive = () => {
    clearInterval(this.timer);
    this.setState({ isIdle: false });
    this.setState({ timeLeftInSeconds: maxWarningTimeBeforeTimeoutInSeconds });
    fetch(keepAliveEndpoint);
  };

  onIdle = () => {
    this.setState({ isIdle: true });
    clearInterval(this.timer);
    this.timer = setInterval(this.countdown, 1000);
  };

  countdown = () => {
    const { timeLeftInSeconds, loggedIn } = this.state;
    const { history, logOut } = this.props;

    if (timeLeftInSeconds === 0 && loggedIn) {
      LogoutUser().then(() => {
        this.setState({ loggedIn: false });
        logOut();
        history.push({
          pathname: '/sign-in',
          state: { timedout: true },
        });
      });
    } else {
      this.setState({ timeLeftInSeconds: timeLeftInSeconds - 1 });
    }
  };

  render() {
    const { isLoggedIn } = this.props;

    const { isIdle, timeLeftInSeconds } = this.state;

    return (
      isLoggedIn && (
        <IdleTimer
          element={document}
          onActive={this.onActive}
          onIdle={this.onIdle}
          timeout={timeToDisplayWarningInMilliseconds}
          events={['blur', 'focus', 'mousedown', 'touchstart', 'MSPointerDown']}
        >
          {isIdle && (
            <Alert type="warning" heading="Inactive user">
              You have been inactive and will be logged out in {timeLeftInSeconds} seconds unless you touch or click on
              the page.
            </Alert>
          )}
        </IdleTimer>
      )
    );
  }
}

LogoutOnInactivity.propTypes = {
  isLoggedIn: PropTypes.bool,
  history: HistoryShape.isRequired,
  logOut: PropTypes.func.isRequired,
};

LogoutOnInactivity.defaultProps = {
  isLoggedIn: false,
};

const mapStateToProps = (state) => {
  return {
    isLoggedIn: selectIsLoggedIn(state),
  };
};

const mapDispatchToProps = {
  logOut: logOutAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(LogoutOnInactivity));
