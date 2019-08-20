import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser, selectGetCurrentUserIsLoading } from 'shared/Data/users';
import SignIn from './SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class PrivateRouteContainer extends React.Component {
  render() {
    const { loginIsLoading, userIsLoggedIn, path, ...props } = this.props;
    if (userIsLoggedIn) return <Route {...props} />;
    else if (loginIsLoading) return <LoadingPlaceholder />;
    else return <Route path={path} component={SignIn} />;
  }
}
const mapStateToProps = state => ({
  loginIsLoading: selectGetCurrentUserIsLoading(state),
  userIsLoggedIn: selectCurrentUser(state).isLoggedIn,
});

const PrivateRoute = connect(mapStateToProps)(PrivateRouteContainer);

export default PrivateRoute;
