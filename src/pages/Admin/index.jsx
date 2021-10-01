import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { get } from 'lodash';
import * as Cookies from 'js-cookie';
import { fetchUtils } from 'react-admin';

import Home from './Home';

import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import SignIn from 'scenes/SystemAdmin/shared/SignIn';
import { isDevelopment } from 'shared/constants';
import { LoginButton } from 'scenes/SystemAdmin/shared/LoginButton';
import { GetLoggedInUser } from 'utils/api';
import FOUOHeader from 'components/FOUOHeader';
import restProvider from 'scenes/SystemAdmin/shared/rest_provider';

const httpClient = (url, options = {}) => {
  const headers = options.headers || new Headers({ Accept: 'application/json' });

  const token = Cookies.get('masked_gorilla_csrf');
  if (!token) {
    milmoveLog(MILMOVE_LOG_LEVEL.WARN, 'Unable to retrieve CSRF Token from cookie');
  }

  headers.set('X-CSRF-TOKEN', token);
  return fetchUtils.fetchJson(url, {
    ...options,
    headers,
    credentials: 'same-origin',
  });
};

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
    const { isLoggedIn } = this.state;
    const { basename, dataProvider } = this.props;

    if (!isLoggedIn) {
      return (
        <>
          <div id="app-root">
            <FOUOHeader />
            <LoginButton showDevlocalButton={get(this.state, 'isDevelopment', isDevelopment)} isLoggedIn={isLoggedIn} />
            <SignIn location={window.location} />
          </div>
          <div id="modal-root" />
        </>
      );
    }

    return <Home basename={basename} dataProvider={dataProvider} />;
  }
}

AdminWrapper.propTypes = {
  basename: PropTypes.string,
  dataProvider: PropTypes.oneOfType([PropTypes.func, PropTypes.shape({})]),
};

AdminWrapper.defaultProps = {
  basename: '/system',
  dataProvider: restProvider('/admin/v1', httpClient),
};

export default AdminWrapper;
