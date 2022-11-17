import React, { Component, lazy } from 'react';
import { Route, Switch, withRouter } from 'react-router-dom-old';

import Home from './Home';
import { GetLoggedInUser } from 'utils/api';
import CUIHeader from 'components/CUIHeader/CUIHeader';
// Lazy load these dependencies (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
const InvalidPermissions = lazy(() => import('pages/InvalidPermissions/InvalidPermissions'));

class AdminWrapper extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoggedIn: false,
    };
  }

  componentDidMount() {
    GetLoggedInUser()
      .then(() => this.setState({ isLoggedIn: true }))
      .catch(() => this.setState({ isLoggedIn: false }));

    document.title = 'Transcom PPP: Admin';
  }

  render() {
    const defaultComponent = this.state.isLoggedIn ? Home : SignIn;
    return (
      <div id="app-root">
        <CUIHeader />
        <Switch>
          {/* no auth */}
          <Route path="/sign-in" component={SignIn} />
          <Route path="/invalid-permissions" component={InvalidPermissions} />
          <Route path="/" component={defaultComponent} />)
        </Switch>
      </div>
    );
  }
}

export default withRouter(AdminWrapper);
