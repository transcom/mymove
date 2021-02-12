import React, { useState } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { bool } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { isDevelopment } from 'shared/constants';
import { LogoutUser } from 'utils/api';
import { logOut } from 'store/auth/actions';
import EulaModal from 'components/EulaModal';

import styles from './LoginButton.module.scss';
import { selectIsLoggedIn } from '../../store/auth/selectors';

const LoginButton = (props) => {
  const [showEula, setShowEula] = useState(false);

  if (!props.isLoggedIn) {
    return (
      <>
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
          {props.useEula ? (
            <Button
              aria-label="Sign In"
              className={styles.signIn}
              data-hook="signin"
              onClick={() => setShowEula(!showEula)}
              type="button"
            >
              Sign In
            </Button>
          ) : (
            <a className="usa-nav__link" data-hook="signin" href="/auth/login-gov">
              Sign In
            </a>
          )}
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

LoginButton.propTypes = {
  useEula: bool,
};

LoginButton.defaultProps = {
  useEula: false,
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
