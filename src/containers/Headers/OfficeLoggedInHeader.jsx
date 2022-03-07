import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link, useHistory } from 'react-router-dom';

import MilMoveHeader from 'components/MilMoveHeader/index';
import OfficeUserInfo from 'components/MilMoveHeader/OfficeUserInfo';
import { LogoutUser } from 'utils/api';
import { logOut as logOutAction } from 'store/auth/actions';
import { OfficeUserInfoShape } from 'types/index';
import { selectLoggedInUser } from 'store/entities/selectors';
import { roleTypes } from 'constants/userRoles';

const OfficeLoggedInHeader = ({ officeUser, activeRole, logOut }) => {
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

  let queueText = '';
  if (activeRole === roleTypes.TOO) {
    queueText = 'moves';
  } else if (activeRole === roleTypes.TIO) {
    queueText = 'payment requests';
  }

  return (
    <MilMoveHeader>
      {officeUser?.transportation_office && (
        <ul className="usa-nav__primary">
          <li className="usa-nav__primary-item">
            <Link to="/">
              {officeUser.transportation_office.gbloc} {queueText}
            </Link>
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
