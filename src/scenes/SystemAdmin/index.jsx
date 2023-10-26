import React, { Component, lazy } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';

import Home from './Home';

import { GetLoggedInUser } from 'utils/api';
// Logger
import { milmoveLogger } from 'utils/milmoveLog';
import { retryPageLoading } from 'utils/retryPageLoading';
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
  }

  componentDidCatch(error, info) {
    const { message } = error;
    milmoveLogger.error({ message, info });
    retryPageLoading(error);
  }

  render() {
    return (
      <>
        <div id="app-root">
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
