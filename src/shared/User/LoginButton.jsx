import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { Button } from '@trussworks/react-uswds';

import { selectCurrentUser } from 'shared/Data/users';
import { isDevelopment } from 'shared/constants';
import { LogoutUser } from 'shared/User/api.js';
import { logOut } from 'store/auth/actions';

const LoginButton = (props) => {
  if (!props.isLoggedIn) {
    return (
      <React.Fragment>
        {props.showDevlocalButton && (
          <a
            className="usa-nav__link"
            data-hook="devlocal-signin"
            style={{ marginRight: '2em' }}
            href="/devlocal-auth/login"
          >
            Local Sign In
          </a>
        )}
        <a className="usa-nav__link" data-hook="signin" href="/auth/login-gov">
          Sign In
        </a>
      </React.Fragment>
    );
  } else {
    const handleLogOut = () => {
      props.logOut();
      LogoutUser();
    };

    return (
      <Button unstyled className="usa-nav__link" onClick={handleLogOut} type="button">
        Sign out
      </Button>
    );
  }
};

function mapStateToProps(state) {
  const user = selectCurrentUser(state);
  return {
    isLoggedIn: user.isLoggedIn,
    showDevlocalButton: get(state, 'isDevelopment', isDevelopment),
  };
}

const mapDispatchToProps = {
  logOut,
};

export default connect(mapStateToProps, mapDispatchToProps)(LoginButton);
