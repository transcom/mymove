import React from 'react';
import { connect } from 'react-redux';
import { compose } from 'redux';
// this is from https://hackernoon.com/role-based-authorization-in-react-c70bb7641db4
// as of 3/8 we are not using this yet, but it seems like it would work once we have a need to include roles in user state
const AuthorizationContainer = (WrappedComponent, feature) =>
  function WithAuthorization(props) {
    const { features } = props;
    if (!features || !features.includes(feature)) {
      return (
        <div className="usa-grid">
          <h1>You are not authorized to view this page</h1>
        </div>
      );
    }

    return <WrappedComponent {...props} />;
  };

const mapStateToProps = state => ({
  features: state.user.features,
});

// see https://medium.com/practo-engineering/connected-higher-order-component-hoc-93ee63c91526
const Authorization = compose(connect(mapStateToProps, null), AuthorizationContainer);
export default Authorization;
