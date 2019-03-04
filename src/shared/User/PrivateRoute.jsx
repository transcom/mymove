import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser } from 'shared/Data/users';
import SignIn from './SignIn';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class PrivateRouteContainer extends React.Component {
  render() {
    const { isLoggedIn, path, ...props } = this.props;
    if (isLoggedIn) return <Route {...props} />;
    else return <Route path={path} component={SignIn} />;
  }
}
const mapStateToProps = state => {
  const user = selectCurrentUser(state);
  return {
    isLoggedIn: user.isLoggedIn,
  };
};
const PrivateRoute = connect(mapStateToProps)(PrivateRouteContainer);

export default PrivateRoute;
