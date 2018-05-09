import React, { Component } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';

import { createServiceMember } from 'scenes/ServiceMembers/ducks';
import Alert from 'shared/Alert';
import LoginButton from 'shared/User/LoginButton';

export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
  }
  componentDidUpdate() {
    if (this.props.createdServiceMemberSuccess) {
      this.props.push(
        `service-member/${this.props.createdServiceMember.id}/create`,
      );
    }
  }
  startMove = values => {
    if (this.props.loggedInUser.service_member) {
      this.props.push(
        `service-member/${this.props.loggedInUser.service_member.id}/create`,
      );
    } else {
      this.props.createServiceMember({});
    }
  };
  render() {
    const {
      isLoggedIn,
      loggedInUserIsLoading,
      loggedInUserSuccess,
      loggedInUserError,
      createdServiceMemberError,
    } = this.props;
    return (
      <div className="usa-grid">
        <h1>Welcome! </h1>
        <div>
          {loggedInUserError && (
            <Alert type="error" heading="An error occurred">
              There was an error loading your user information.
            </Alert>
          )}
          {createdServiceMemberError && (
            <Alert type="error" heading="An error occurred">
              There was an error creating your move.
            </Alert>
          )}
          {loggedInUserIsLoading && <span> Loading... </span>}
          {!isLoggedIn && <LoginButton />}
          {loggedInUserSuccess && (
            <button onClick={this.startMove}>Start a move</button>
          )}
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.user.isLoggedIn,
  loggedInUser: state.loggedInUser.loggedInUser,
  loggedInUserIsLoading: state.loggedInUser.isLoading,
  loggedInUserError: state.loggedInUser.error,
  loggedInUserSuccess: state.loggedInUser.hasSucceeded,
  createdServiceMemberSuccess: state.serviceMember.hasSubmitSuccess,
  createdServiceMemberError: state.serviceMember.error,
  createdServiceMember: state.serviceMember.currentServiceMember,
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, createServiceMember }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Landing);
