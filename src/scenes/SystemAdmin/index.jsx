import React, { Component } from 'react';

import { get } from 'lodash';
import Home from './Home';
import SignIn from 'scenes/SystemAdmin/shared/SignIn';
import { isDevelopment } from 'shared/constants';
import { LoginButton } from 'scenes/SystemAdmin/shared/LoginButton';
import { GetLoggedInUser } from 'utils/api';
import CUIHeader from 'components/CUIHeader/CUIHeader';

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
    if (!this.state.isLoggedIn) {
      return (
        <>
          <div id="app-root">
            <CUIHeader />
            <LoginButton
              showDevlocalButton={get(this.state, 'isDevelopment', isDevelopment)}
              isLoggedIn={this.state.isLoggedIn}
            />
            <SignIn location={window.location} />
          </div>
          <div id="modal-root" />
        </>
      );
    } else {
      return <Home />;
    }
  }
}

export default AdminWrapper;
