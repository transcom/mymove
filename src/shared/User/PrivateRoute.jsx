import React from 'react';
import { Route } from 'react-router-dom';

import { connect } from 'react-redux';

import LoginButton from 'shared/User/LoginButton';

class PrivateRouteContainer extends React.Component {
  render() {
    const { isLoggedIn, component: Component, ...props } = this.props;

    return (
      <Route
        {...props}
        render={props =>
          isLoggedIn ? (
            <Component {...props} />
          ) : (
            <div className="usa-grid">
              <h1>Please login to access this page </h1>
              <LoginButton />
            </div>
          )
        }
      />
    );
  }
}

const PrivateRoute = connect(state => ({
  isLoggedIn: state.user.isLoggedIn,
}))(PrivateRouteContainer);

export default PrivateRoute;
