import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { get } from 'lodash';

import MilMoveHeader from 'components/MilMoveHeader/index';
import CustomerUserInfo from 'components/MilMoveHeader/CustomerUserInfo';
import { LogoutUser } from 'utils/api';
import { isDevelopment } from 'shared/constants';
import { logOut as logOutAction } from 'store/auth/actions';
import { selectIsProfileComplete } from 'store/entities/selectors';

const CustomerLoggedInHeader = ({ isProfileComplete, logOut, isLocalSignIn }) => {
  const navigate = useNavigate();
  const handleLogout = () => {
    logOut();
    LogoutUser().then((r) => {
      const redirectURL = r.body;
      if (redirectURL && !isLocalSignIn) {
        window.location.href = redirectURL;
      } else {
        navigate('/sign-in', { state: { hasLoggedOut: true } });
      }
    });
  };

  return (
    <MilMoveHeader>
      <CustomerUserInfo showProfileLink={isProfileComplete} handleLogout={handleLogout} />
    </MilMoveHeader>
  );
};

CustomerLoggedInHeader.propTypes = {
  isProfileComplete: PropTypes.bool,
  logOut: PropTypes.func.isRequired,
};

CustomerLoggedInHeader.defaultProps = {
  isProfileComplete: false,
};

const mapStateToProps = (state) => ({
  isProfileComplete: selectIsProfileComplete(state),
  isLocalSignIn: get(state, 'isDevelopment', isDevelopment),
});

const mapDispatchToProps = {
  logOut: logOutAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerLoggedInHeader);
