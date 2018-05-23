import React, { Component } from 'react';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';

import { MoveSummary } from './MoveSummary';

import { createServiceMember } from 'scenes/ServiceMembers/ducks';
import { loadEntitlements } from 'scenes/Orders/ducks';
import { loadLoggedInUser } from 'shared/User/ducks';
import { getNextIncompletePage } from 'scenes/MyMove/getWorkflowRoutes';
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
    const { service_member } = this.props;
    if (this.props.loggedInUserSuccess) {
      if (
        !this.props.createdServiceMemberIsLoading &&
        isEmpty(service_member)
      ) {
        // Once the logged in user loads, if the service member doesn't
        // exist we need to dispatch creating one, once.
        this.props.createServiceMember({});
      } else if (
        !isEmpty(service_member) &&
        !service_member.is_profile_complete
      ) {
        // If the service member exists, but is not complete, redirect to next incomplete page.
        this.resumeMove();
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

  resumeMove = () => {
    this.props.push(getNextIncompletePage(this.props.service_member));
  };
  render() {
    const {
      isLoggedIn,
      loggedInUserIsLoading,
      loggedInUserError,
      createdServiceMemberError,
      loggedInUser,
      moveSubmitSuccess,
      entitlement,
    } = this.props;

    const profile = get(loggedInUser, 'service_member', {});
    const orders = get(profile, 'orders.0');
    const move = get(orders, 'moves.0');
    const ppm = get(move, 'personally_procured_moves.0', {});

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
        {!isLoggedIn && <LoginButton />}
        {isLoggedIn && (
          <MoveSummary
            entitlement={entitlement}
            profile={profile}
            orders={orders}
            move={move}
            ppm={ppm}
            editMove={this.editMove}
            resumeMove={this.resumeMove}
          />
        )}
      </div>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.user.isLoggedIn,
  service_member: get(state, 'loggedInUser.loggedInUser.service_member', {}),
  loggedInUser: state.loggedInUser.loggedInUser,
  loggedInUserIsLoading: state.loggedInUser.isLoading,
  loggedInUserError: state.loggedInUser.error,
  loggedInUserSuccess: state.loggedInUser.hasSucceeded,
  createdServiceMemberIsLoading: state.serviceMember.isLoading,
  createdServiceMemberSuccess: state.serviceMember.hasSubmitSuccess,
  createdServiceMemberError: state.serviceMember.error,
  createdServiceMember: state.serviceMember.currentServiceMember,
  moveSubmitSuccess: state.signedCertification.moveSubmitSuccess,
  entitlement: loadEntitlements(state),
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { push, createServiceMember, loadLoggedInUser },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(Landing);
