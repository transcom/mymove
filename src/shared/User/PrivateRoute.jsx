import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';

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
              <Redirect
                to={{
                  pathname: '/',
                  state: { from: props.location },
                }}
              />
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
