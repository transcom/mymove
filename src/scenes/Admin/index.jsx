import React, { Component } from 'react';
import { Redirect, Switch, Route } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import AdminHeader from 'shared/Header/Admin';
import { getCurrentUserInfo } from 'shared/Data/users';
import { no_op } from 'shared/utils';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';
import { loadPublicSchema } from 'shared/Swagger/ducks';
import restProvider from 'ra-data-simple-rest';
import { Admin } from 'react-admin';

const dataProvider = restProvider('http://admin/v1/...');

const AdminHome = () => (
  <div className="admin-system-wrapper">
    <Admin dataProvider={dataProvider} history={history}>
      {/*<Resource />*/}
    </Admin>
  </div>
);

class AdminWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Admin';
    this.props.loadPublicSchema();
    this.props.getCurrentUserInfo();
  }

  render() {
    return (
      <ConnectedRouter history={history}>
        <div className="Admin site">
          <AdminHeader />
          <main className="site__content">
            <div>
              <LogoutOnInactivity />
              <Switch>
                <Route
                  exact
                  path="/"
                  component={({ location }) => (
                    <Redirect
                      from="/"
                      to={{
                        ...location,
                        pathname: '/system',
                      }}
                    />
                  )}
                />
                <PrivateRoute path="/system" component={AdminHome} />
              </Switch>
            </div>
          </main>
        </div>
      </ConnectedRouter>
    );
  }
}

AdminWrapper.defaultProps = {
  loadPublicSchema: no_op,
  getCurrentUserInfo: no_op,
};

const mapStateToProps = state => ({
  swaggerError: state.swaggerPublic.hasErrored,
});

const mapDispatchToProps = dispatch => bindActionCreators({ loadPublicSchema, getCurrentUserInfo }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AdminWrapper);
