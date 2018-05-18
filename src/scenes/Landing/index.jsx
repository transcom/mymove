import React, { Component } from 'react';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';

import { MoveSummary } from './MoveSummary';

import { createServiceMember } from 'scenes/ServiceMembers/ducks';
import { loadLoggedInUser } from 'shared/User/ducks';
import Alert from 'shared/Alert';
import LoginButton from 'shared/User/LoginButton';

export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
    if (!this.props.loggedInUserIsLoading) {
      this.props.loadLoggedInUser();
    }
    window.scrollTo(0, 0);
  }
  componentDidUpdate() {
    if (this.props.loggedInUserSuccess) {
      // Once the logged in user loads, if the service member doesn't
      // exist we need to dispatch creating one, once.
      if (
        !this.props.createdServiceMemberIsLoading &&
        !this.props.loggedInUser.service_member
      ) {
        this.props.createServiceMember({}).then(something => {
          this.props.push(
            `service-member/${
              this.props.loggedInUser.service_member.id
            }/create`,
          );
        });
        // If the service member exists, but is not complete, redirect as well.
      } else if (
        this.props.loggedInUser &&
        this.props.loggedInUser.service_member &&
        !this.props.loggedInUser.service_member.is_profile_complete
      ) {
        this.props.push(
          `service-member/${this.props.loggedInUser.service_member.id}/create`,
        );
      }
    }
  }
  startMove = values => {
    if (!this.props.loggedInUser.service_member) {
      console.error(
        'With no service member, you should have been redirected already.',
      );
    }
    this.props.push(
      `service-member/${this.props.loggedInUser.service_member.id}/create`,
    );
  };

  editMove = move => {
    this.props.push(`moves/${move.id}/review`);
  };

  render() {
    const {
      isLoggedIn,
      loggedInUserIsLoading,
      loggedInUserSuccess,
      loggedInUserError,
      createdServiceMemberError,
      loggedInUser,
      moveSubmitSuccess,
    } = this.props;

    let profile = get(loggedInUser, 'service_member');
    let orders = get(profile, 'orders.0');
    let move = get(orders, 'moves.0');
    let ppm = get(move, 'personally_procured_moves.0');

    const displayMove = !!ppm;

    return (
      <div className="usa-grid">
        <div>
          {moveSubmitSuccess && (
            <Alert type="success" heading="Success">
              You've submitted your move
            </Alert>
          )}
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
        </div>
        {displayMove && (
          <MoveSummary
            profile={profile}
            orders={orders}
            move={move}
            ppm={ppm}
            editMove={this.editMove}
          />
        )}

        {!isLoggedIn && <LoginButton />}
        {loggedInUserSuccess && (
          <button onClick={this.startMove}>Start a move</button>
        )}
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
  createdServiceMemberIsLoading: state.serviceMember.isLoading,
  createdServiceMemberSuccess: state.serviceMember.hasSubmitSuccess,
  createdServiceMemberError: state.serviceMember.error,
  createdServiceMember: state.serviceMember.currentServiceMember,
  moveSubmitSuccess: state.signedCertification.moveSubmitSuccess,
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { push, createServiceMember, loadLoggedInUser },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(Landing);
