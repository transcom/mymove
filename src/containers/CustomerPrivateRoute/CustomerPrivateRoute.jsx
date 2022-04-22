import React from 'react';
import PropTypes from 'prop-types';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { selectGetCurrentUserIsLoading, selectIsLoggedIn } from 'store/auth/selectors';
import { LocationShape } from 'types/index';

const CustomerPrivateRoute = ({ loginIsLoading, userIsLoggedIn, location, ...routeProps }) => {
  if (loginIsLoading) return <LoadingPlaceholder />;

  const { hash, search } = location;

  if (!userIsLoggedIn) {
    return (
      <Redirect
        to={{
          pathname: '/sign-in',
          hash,
          search,
        }}
      />
    );
  }

  // eslint-disable-next-line react/jsx-props-no-spreading
  return <Route {...routeProps} />;
};

CustomerPrivateRoute.propTypes = {
  loginIsLoading: PropTypes.bool,
  userIsLoggedIn: PropTypes.bool,
  location: LocationShape,
};

CustomerPrivateRoute.defaultProps = {
  loginIsLoading: true,
  userIsLoggedIn: false,
  location: {},
};

const mapStateToProps = (state) => {
  return {
    loginIsLoading: selectGetCurrentUserIsLoading(state),
    userIsLoggedIn: selectIsLoggedIn(state),
  };
};

export default connect(mapStateToProps)(CustomerPrivateRoute);
