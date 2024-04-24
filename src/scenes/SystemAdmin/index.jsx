import React, { Component, lazy } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';

import Home from './Home';

import { GetLoggedInUser } from 'utils/api';
// Logger
import { milmoveLogger } from 'utils/milmoveLog';
import { retryPageLoading } from 'utils/retryPageLoading';
import { OktaLoggedOutBanner, OktaNeedsLoggedOutBanner } from 'components/OktaLogoutBanner';
import CUIHeader from 'components/CUIHeader/CUIHeader';

// Lazy load these dependencies (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
const InvalidPermissions = lazy(() => import('pages/InvalidPermissions/InvalidPermissions'));

class AdminWrapper extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoggedIn: false,
      oktaLoggedOut: undefined,
      oktaNeedsLoggedOut: undefined,
    };
  }

  componentDidMount() {
    GetLoggedInUser()
      .then(() => this.setState({ isLoggedIn: true }))
      .catch(() => this.setState({ isLoggedIn: false }));
    // We need to check if the user was redirected back from Okta after logging out
    // This can occur when they click "sign out" or if they try to access MM
    // while still logged into Okta which will force a redirect to logout
    const currentUrl = new URL(window.location.href);
    const oktaLoggedOutParam = currentUrl.searchParams.get('okta_logged_out');

    // If the params "okta_logged_out=true" are in the url, we will change some state
    // so a banner will display
    if (oktaLoggedOutParam === 'true') {
      this.setState({
        oktaLoggedOut: true,
      });
    } else if (oktaLoggedOutParam === 'false') {
      this.setState({
        oktaNeedsLoggedOut: true,
      });
    }
  }

  componentDidCatch(error, info) {
    const { message } = error;
    milmoveLogger.error({ message, info });
    retryPageLoading(error);
  }

  render() {
    const { oktaLoggedOut, oktaNeedsLoggedOut } = this.state;
    const script = document.createElement('script');

    script.src = '//rum-static.pingdom.net/pa-6567b05deff3250012000426.js';
    script.async = true;

    document.body.appendChild(script);
    return (
      <>
        <div id="app-root">
          <CUIHeader />
          {oktaLoggedOut && <OktaLoggedOutBanner />}
          {oktaNeedsLoggedOut && <OktaNeedsLoggedOutBanner />}
          <Routes>
            {/* no auth */}
            <Route path="/sign-in" element={<SignIn />} />
            <Route path="/invalid-permissions" element={<InvalidPermissions />} />
            {/* system is basename of admin app, see https://marmelab.com/react-admin/Routing.html#using-react-admin-inside-a-route */}
            <Route path="/system/*" element={this.state.isLoggedIn ? <Home /> : <SignIn />} />)
            <Route path="*" element={<Navigate to="/system" />} />
          </Routes>
        </div>
        <div id="modal-root" />
      </>
    );
  }
}

export default AdminWrapper;
