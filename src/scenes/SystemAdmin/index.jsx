import React, { Component } from 'react';

import { no_op } from 'shared/utils.js';
import { get } from 'lodash';
import Home from './Home';
import SignIn from 'scenes/SystemAdmin/shared/SignIn';
import { isDevelopment } from 'shared/constants';
import { LoginButton } from 'scenes/SystemAdmin/shared/LoginButton';
import { GetLoggedInUser } from 'shared/User/api.js';
import FOUOHeader from 'components/FOUOHeader';

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
    this.props.loadPublicSchema();
    this.props.getCurrentUserInfo();
  }

  render() {
    if (!this.state.isLoggedIn) {
      return (
        <React.Fragment>
          <FOUOHeader />
          <LoginButton
            showDevlocalButton={get(this.state, 'isDevelopment', isDevelopment)}
            isLoggedIn={this.state.isLoggedIn}
          />
          <SignIn location={window.location} />
        </React.Fragment>
      );
    } else {
      return <Home />;
    }
  }
}

AdminWrapper.defaultProps = {
  loadPublicSchema: no_op,
  getCurrentUserInfo: no_op,
};

export default AdminWrapper;
