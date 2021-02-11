import React, { useState } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { selectCurrentUser } from 'shared/Data/users';
import { isDevelopment } from 'shared/constants';
import { LogoutUser } from 'utils/api';
import { logOut } from 'store/auth/actions';
import { Button, Overlay } from '@trussworks/react-uswds';
import EulaModal from 'components/EulaModal';

import styles from './LoginButton.module.scss';

const LoginButton = (props) => {
  const [showEula, setShowEula] = useState(false);

  if (!props.isLoggedIn) {
    return (
      <>
        {showEula ? <Overlay /> : ''}
        <EulaModal
          isOpen={showEula}
          acceptTerms={() => {
            window.location.href = '/auth/login-gov';
          }}
          closeModal={() => setShowEula(false)}
        />
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
          <Button aria-label="Sign In" className={styles.signIn} onClick={() => setShowEula(!showEula)} type="button">
            Sign In
          </Button>
        </li>
      </>
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
