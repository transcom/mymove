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
    const isLoggedOut =
      this.props.loggedInUserError &&
      this.props.loggedInUserError.statusCode === 401;
    const unknownUserError = this.props.loggedInUserError && !isLoggedOut;
    return (
      <div className="usa-grid">
        <h1>Welcome! </h1>
        <div>
          {unknownUserError && (
            <Alert type="error" heading="An error occurred">
              There was an error starting your move.
            </Alert>
          )}
          {this.props.loggedInUserIsLoading && <span> Loading... </span>}
          {!this.props.loggedInUserIsLoading && isLoggedOut && <LoginButton />}
          <div />
        </div>
        {!this.props.loggedInUserIsLoading &&
          this.props.loggedInUser && (
            <button onClick={this.startMove}>Start a move</button>
          )}
      </div>
    );
  }
}

const mapStateToProps = state => ({
  loggedInUser: state.loggedInUser.loggedInUser,
  loggedInUserIsLoading: state.loggedInUser.isLoading,
  loggedInUserError: state.loggedInUser.error,
  loggedInUserSuccess: state.loggedInUser.hasSucceeded,
  createdServiceMemberSuccess: state.serviceMember.hasSubmitSuccess,
  createdServiceMember: state.serviceMember.currentServiceMember,
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, createServiceMember }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Landing);
