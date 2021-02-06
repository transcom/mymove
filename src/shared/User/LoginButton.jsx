import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { selectIsLoggedIn } from 'store/entities/selectors';
import { isDevelopment } from 'shared/constants';
import { LogoutUser } from 'utils/api';
import { logOut } from 'store/auth/actions';

const LoginButton = (props) => {
  if (!props.isLoggedIn) {
    return (
      <React.Fragment>
        {props.showDevlocalButton && (
          <li className="usa-nav__primary-item">
            <a
              className="usa-nav__link"
              data-hook="devlocal-signin"
              style={{ marginRight: '2em' }}
              href="/devlocal-auth/login"
            >
              Local Sign In
            </a>
          </li>
        )}
        <li className="usa-nav__primary-item">
          <a className="usa-nav__link" data-hook="signin" href="/auth/login-gov">
            Sign In
          </a>
        </li>
      </React.Fragment>
    );
  } else {
    const handleLogOut = () => {
      props.logOut();
      LogoutUser();
    };

    return (
      <li className="usa-nav__primary-item">
        <a className="usa-nav__link" href="#" onClick={handleLogOut}>
          Sign Out
        </a>
      </li>
    );
  }
};

function mapStateToProps(state) {
  return {
    isLoggedIn: selectIsLoggedIn(state),
    showDevlocalButton: get(state, 'isDevelopment', isDevelopment),
  };
}

const mapDispatchToProps = {
  logOut,
};

export default connect(mapStateToProps, mapDispatchToProps)(LoginButton);
