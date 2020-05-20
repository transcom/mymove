import React from 'react';

import { Component, Fragment } from 'react';
import { Redirect } from 'react-router-dom';

export class UserSearch extends Component {
  state = { ...this.initialState };

  get initialState() {
    return {
      showUser: false,
    };
  }

  setUserIDinState = (e) => {
    this.setState({ userID: e.target.value });
  };

  redirectToShowUser = () => {
    this.setState({ showUser: true });
  };

  render() {
    if (!this.state.showUser) {
      return (
        <Fragment>
          <span>Search by user ID</span>
          <form onSubmit={this.redirectToShowUser}>
            <input onChange={this.setUserIDinState} name="userID" component="input" type="text" />
            <button type="submit">Search</button>
          </form>
        </Fragment>
      );
    } else {
      return <Redirect to={`/users/${this.state.userID}/show`} />;
    }
  }
}

export default UserSearch;
