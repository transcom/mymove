import React from 'react';
import PropTypes from 'prop-types';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { fetchAccessCode as fetchAccessCodeAction } from 'shared/Entities/modules/accessCodes';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { selectGetCurrentUserIsLoading, selectIsLoggedIn } from 'store/auth/selectors';
import { LocationShape } from 'types/index';

class CustomerPrivateRoute extends React.Component {
  componentDidMount() {
    const { fetchAccessCode } = this.props;
    fetchAccessCode();
  }

  render() {
    const { loginIsLoading, userIsLoggedIn, requiresAccessCode, accessCode, location, ...routeProps } = this.props;
    if (loginIsLoading) return <LoadingPlaceholder />;

    const { hash, search } = location;

    if (!userIsLoggedIn)
      return (
        <Redirect
          to={{
            pathname: '/sign-in',
            hash,
            search,
          }}
        />
      );

    if (userIsLoggedIn && requiresAccessCode && !accessCode) return <Redirect to="/access-code" />;

    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Route {...routeProps} />;
  }
}

CustomerPrivateRoute.propTypes = {
  fetchAccessCode: PropTypes.func.isRequired,
  loginIsLoading: PropTypes.bool,
  userIsLoggedIn: PropTypes.bool,
  requiresAccessCode: PropTypes.bool,
  accessCode: PropTypes.string,
  location: LocationShape,
};

CustomerPrivateRoute.defaultProps = {
  loginIsLoading: true,
  userIsLoggedIn: false,
  requiresAccessCode: false,
  accessCode: undefined,
  location: {},
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const { accessCodes = {} } = state.entities;

  return {
    loginIsLoading: selectGetCurrentUserIsLoading(state),
    userIsLoggedIn: selectIsLoggedIn(state),
    requiresAccessCode: serviceMember?.requires_access_code,
    accessCode: accessCodes && Object.values(accessCodes).length > 0 ? Object.values(accessCodes)[0].code : null,
  };
};

const mapDispatchToProps = { fetchAccessCode: fetchAccessCodeAction };

export default connect(mapStateToProps, mapDispatchToProps)(CustomerPrivateRoute);
