import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import classnames from 'classnames';

import GblocSwitcher from 'components/Office/GblocSwitcher/GblocSwitcher';
import MilMoveHeader from 'components/MilMoveHeader/index';
import OfficeUserInfo from 'components/MilMoveHeader/OfficeUserInfo';
import { LogoutUser } from 'utils/api';
import { logOut as logOutAction } from 'store/auth/actions';
import { OfficeUserInfoShape } from 'types/index';
import { selectLoggedInUser } from 'store/entities/selectors';
import { roleTypes } from 'constants/userRoles';
import { unlockOfficerUserMoves } from 'services/ghcApi';

const OfficeLoggedInHeader = ({ officeUser, activeRole, logOut }) => {
  const navigate = useNavigate();
  const handleLogout = () => {
    // explicit clear session storage
    window.sessionStorage.clear();
    logOut();
    LogoutUser().then((r) => {
      const redirectURL = r.body;
      // checking to see if "Local Sign In" button was used to sign in
      const urlParams = new URLSearchParams(redirectURL.split('?')[1]);
      const idTokenHint = urlParams.get('id_token_hint');
      if (redirectURL && idTokenHint !== 'devlocal') {
        window.location.href = redirectURL;
      } else {
        navigate('/sign-in', { state: { hasLoggedOut: true } });
      }
    });
  };

  let queueText = '';
  const location = useLocation();
  if (activeRole === roleTypes.TOO) {
    queueText = 'moves';
  } else if (activeRole === roleTypes.TIO) {
    queueText = 'payment requests';
  } else if (activeRole === roleTypes.QAE && location.pathname === '/') {
    unlockOfficerUserMoves('move', officeUser.id);
  }

  return (
    <MilMoveHeader>
      {officeUser?.transportation_office && (
        <ul className="usa-nav__primary">
          <li className={classnames('usa-nav__primary-item')}>
            {activeRole === roleTypes.HQ ? (
              <GblocSwitcher officeUsersDefaultGbloc={officeUser.transportation_office.gbloc} />
            ) : (
              <Link to="/">
                {officeUser.transportation_office.gbloc} {queueText}
              </Link>
            )}
          </li>
        </ul>
      )}
      <OfficeUserInfo lastName={officeUser.last_name} firstName={officeUser.first_name} handleLogout={handleLogout} />
    </MilMoveHeader>
  );
};

OfficeLoggedInHeader.propTypes = {
  officeUser: OfficeUserInfoShape,
  activeRole: PropTypes.string,
  logOut: PropTypes.func.isRequired,
};

OfficeLoggedInHeader.defaultProps = {
  officeUser: {},
  activeRole: null,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
    activeRole: state.auth.activeRole,
  };
};

const mapDispatchToProps = {
  logOut: logOutAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(OfficeLoggedInHeader);
