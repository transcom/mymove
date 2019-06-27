import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectGetCurrentUserIsSuccess, selectGetCurrentUserIsError } from 'shared/Data/users';
import SignIn from './SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class PrivateRouteContainer extends React.Component {
  render() {
    const { loginHasSucceeded, loginHasErrored, path, ...props } = this.props;
    if (loginHasSucceeded) return <Route {...props} />;
    else if (loginHasErrored) return <Route path={path} component={SignIn} />;
    else return <LoadingPlaceholder />;
  }
}
const mapStateToProps = state => ({
  loginHasErrored: selectGetCurrentUserIsError(state),
  loginHasSucceeded: selectGetCurrentUserIsSuccess(state),
});

const PrivateRoute = connect(mapStateToProps)(PrivateRouteContainer);

export default PrivateRoute;
