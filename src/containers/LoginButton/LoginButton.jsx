import React, { useState } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { bool, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { selectIsLoggedIn } from '../../store/auth/selectors';

import styles from './LoginButton.module.scss';

import { isDevelopment } from 'shared/constants';
import { LogoutUser } from 'utils/api';
import { logOut as logOutFunction } from 'store/auth/actions';
import ConnectedEulaModal from 'components/EulaModal';

const LoginButton = ({ isLoggedIn, logOut, showDevlocalButton }) => {
  const [showEula, setShowEula] = useState(false);

  if (!isLoggedIn) {
    return (
      <>
        <ConnectedEulaModal
          isOpen={showEula}
          acceptTerms={() => {
            window.location.href = '/auth/login-gov';
          }}
          closeModal={() => setShowEula(false)}
        />
        {showDevlocalButton && (
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
          <Button
            aria-label="Sign In"
            className={styles.signIn}
            data-hook="signin"
            onClick={() => setShowEula(!showEula)}
            type="button"
          >
            Sign In
          </Button>
        </li>
      </>
    );
  }
  const handleLogOut = () => {
    logOut();
    LogoutUser();
  };

  return (
    <li className="usa-nav__primary-item">
      {/* eslint-disable-next-line jsx-a11y/anchor-is-valid */}
      <a className="usa-nav__link" href="#" onClick={handleLogOut}>
        Sign Out
      </a>
    </li>
  );
};

LoginButton.propTypes = {
  isLoggedIn: bool.isRequired,
  logOut: func.isRequired,
  showDevlocalButton: bool.isRequired,
};

function mapStateToProps(state) {
  return {
    isLoggedIn: selectIsLoggedIn(state),
    showDevlocalButton: get(state, 'isDevelopment', isDevelopment),
  };
}

const mapDispatchToProps = {
  logOut: logOutFunction,
};

export default connect(mapStateToProps, mapDispatchToProps)(LoginButton);
