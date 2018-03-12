import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import LoginButton from './LoginButton';

const NotAuthenticated = ({ location }) => (
  <div>
    <h3>
      Please login to access <code>{location.pathname}</code>
    </h3>
  </div>
);

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class PrivateRouteContainer extends React.Component {
  render() {
    const { isLoggedIn, component: Component, ...props } = this.props;
    return (
      <Route
        {...props}
        render={props => {
          if (isLoggedIn) {
            return <Component {...props} />;
          } else {
            return (
              <div className="usa-grid">
                <NotAuthenticated location={props.location} />
                <LoginButton />
              </div>
            );
          }
        }}
      />
    );
  }
}
const mapStateToProps = state => ({
  isLoggedIn: state.user.isLoggedIn,
});
const PrivateRoute = connect(mapStateToProps)(PrivateRouteContainer);

export default PrivateRoute;
