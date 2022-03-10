import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';

import MilMoveHeader from 'components/MilMoveHeader/index';
import CustomerUserInfo from 'components/MilMoveHeader/CustomerUserInfo';
import { LogoutUser } from 'utils/api';
import { logOut as logOutAction } from 'store/auth/actions';
import { selectIsProfileComplete } from 'store/entities/selectors';

const CustomerLoggedInHeader = ({ isProfileComplete, logOut }) => {
  const history = useHistory();
  const handleLogout = () => {
    logOut();
    LogoutUser().then((r) => {
      const redirectURL = r.body;
      if (redirectURL) {
        window.location.href = redirectURL;
      } else {
        history.push({
          pathname: '/sign-in',
          state: { hasLoggedOut: true },
        });
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
});

const mapDispatchToProps = {
  logOut: logOutAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerLoggedInHeader);
