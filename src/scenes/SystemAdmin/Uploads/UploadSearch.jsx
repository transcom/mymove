import React from 'react';

import { Component, Fragment } from 'react';
import { Navigate } from 'react-router-dom';

export class UploadSearch extends Component {
  state = { ...this.initialState };

  get initialState() {
    return {
      showUpload: false,
    };
  }

  setUploadIDinState = (e) => {
    this.setState({ uploadID: e.target.value });
  };

  redirectToShowUpload = () => {
    this.setState({ showUpload: true });
  };

  render() {
    if (!this.state.showUpload) {
      return (
        <Fragment>
          <span>Search by upload ID</span>
          <form onSubmit={this.redirectToShowUpload}>
            <input onChange={this.setUploadIDinState} name="uploadID" component="input" type="text" />
            <button type="submit">Search</button>
          </form>
        </Fragment>
      );
    } else {
      return <Navigate to={`/uploads/${this.state.uploadID}/show`} replace />;
    }
  }
}

export default UploadSearch;
