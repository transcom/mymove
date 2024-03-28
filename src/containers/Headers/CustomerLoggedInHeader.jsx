import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { useLocation, useNavigate } from 'react-router-dom';

import MilMoveHeader from 'components/MilMoveHeader/index';
import CustomerUserInfo from 'components/MilMoveHeader/CustomerUserInfo';
import { LogoutUser } from 'utils/api';
import { logOut as logOutAction } from 'store/auth/actions';
import { selectIsProfileComplete } from 'store/entities/selectors';
import { selectCurrentMoveId } from 'store/general/selectors';

const CustomerLoggedInHeader = ({ state, isProfileComplete, logOut, moveId }) => {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const moveID = pathname.split('/')[2];

  let specialOrderType = '';
  if (Object.keys(state.entities.orders).length > 0) {
    const currentOrderType = Object.values(state.entities.orders).filter((order) => order.moves[0] === moveID)[0];
    specialOrderType = currentOrderType?.orders_type;
  }

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
    <MilMoveHeader specialOrderType={specialOrderType}>
      <CustomerUserInfo showProfileLink={isProfileComplete} handleLogout={handleLogout} moveId={moveId} />
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
  state,
  isProfileComplete: selectIsProfileComplete(state),
  // Grab moveId from state that was set from the most recent navigation to a move
  moveId: selectCurrentMoveId(state),
});

const mapDispatchToProps = {
  logOut: logOutAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerLoggedInHeader);
