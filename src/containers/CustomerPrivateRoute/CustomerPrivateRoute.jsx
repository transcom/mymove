import React from 'react';
import PropTypes from 'prop-types';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';
import { get } from 'lodash';

import SignIn from 'pages/SignIn/SignIn';
import AccessCode from 'shared/User/AccessCode';
import { fetchAccessCode as fetchAccessCodeAction } from 'shared/Entities/modules/accessCodes';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { selectIsLoggedIn } from 'store/auth/selectors';

class CustomerPrivateRoute extends React.Component {
  componentDidMount() {
    const { fetchAccessCode } = this.props;
    fetchAccessCode();
  }

  render() {
    const { isLoggedIn, requiresAccessCode, accessCode, path, ...routeProps } = this.props;
    if (!isLoggedIn) return <Route path={path} component={SignIn} />;
    if (isLoggedIn && requiresAccessCode && !accessCode) return <Route path={path} component={AccessCode} />;

    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Route {...routeProps} />;
  }
}

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const accessCodes = get(state, 'entities.accessCodes');

  return {
    isLoggedIn: selectIsLoggedIn(state),
    requiresAccessCode: serviceMember?.requires_access_code,
    accessCode: accessCodes && Object.values(accessCodes).length > 0 ? Object.values(accessCodes)[0].code : null,
  };
};

CustomerPrivateRoute.propTypes = {
  fetchAccessCode: PropTypes.func.isRequired,
  isLoggedIn: PropTypes.bool,
  requiresAccessCode: PropTypes.bool,
  accessCode: PropTypes.string,
  path: PropTypes.string,
};

CustomerPrivateRoute.defaultProps = {
  isLoggedIn: false,
  requiresAccessCode: false,
  accessCode: undefined,
  path: undefined,
};

const mapDispatchToProps = { fetchAccessCode: fetchAccessCodeAction };

export default connect(mapStateToProps, mapDispatchToProps)(CustomerPrivateRoute);
