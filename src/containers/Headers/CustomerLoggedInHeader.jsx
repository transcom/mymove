import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import MilMoveHeader from 'components/MilMoveHeader/index';
import CustomerUserInfo from 'components/MilMoveHeader/CustomerUserInfo';
import { LogoutUser } from 'utils/api';
import { logOut as logOutAction } from 'store/auth/actions';
import { selectCurrentOrders, selectIsProfileComplete } from 'store/entities/selectors';

const CustomerLoggedInHeader = ({ orderType, isProfileComplete, logOut }) => {
  const navigate = useNavigate();
  const isSpecialMove = ['BLUEBARK'].includes(orderType);

  const handleLogout = () => {
    logOut();
    LogoutUser().then((r) => {
      const redirectURL = r.body;
      const urlParams = new URLSearchParams(redirectURL.split('?')[1]);
      const idTokenHint = urlParams.get('id_token_hint');
      if (redirectURL && idTokenHint !== 'devlocal') {
        window.location.href = redirectURL;
      } else {
        navigate('/sign-in', { state: { hasLoggedOut: true } });
      }
    });
  };

  return (
    <MilMoveHeader isSpecialMove={isSpecialMove}>
      <CustomerUserInfo showProfileLink={isProfileComplete} handleLogout={handleLogout} />
    </MilMoveHeader>
  );
};

CustomerLoggedInHeader.propTypes = {
  orderType: PropTypes.string,
  isProfileComplete: PropTypes.bool,
  logOut: PropTypes.func.isRequired,
};

CustomerLoggedInHeader.defaultProps = {
  orderType: '',
  isProfileComplete: false,
};

const mapStateToProps = (state) => ({
  orderType: selectCurrentOrders(state).orders_type,
  isProfileComplete: selectIsProfileComplete(state),
});

const mapDispatchToProps = {
  logOut: logOutAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerLoggedInHeader);
