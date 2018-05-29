import React, { Component, Fragment } from 'react';
import { isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';
import { withLastLocation } from 'react-router-last-location';
import { MoveSummary } from './MoveSummary';

import { createServiceMember } from 'scenes/ServiceMembers/ducks';
import { loadEntitlements } from 'scenes/Orders/ducks';
import { loadLoggedInUser } from 'shared/User/ducks';
import { getNextIncompletePage } from 'scenes/MyMove/getWorkflowRoutes';
import Alert from 'shared/Alert';
import LoginButton from 'shared/User/LoginButton';

export class Landing extends Component {
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  componentDidUpdate() {
    const {
      serviceMember,
      createdServiceMemberIsLoading,
      loggedInUserSuccess,
      createServiceMember,
    } = this.props;
    if (loggedInUserSuccess) {
      if (!createdServiceMemberIsLoading && isEmpty(serviceMember)) {
        // Once the logged in user loads, if the service member doesn't
        // exist we need to dispatch creating one, once.
        createServiceMember({});
      } else if (
        !isEmpty(serviceMember) &&
        !serviceMember.is_profile_complete
      ) {
        // If the service member exists, but is not complete, redirect to next incomplete page.
        this.resumeMove();
      }
    }
  }
  startMove = values => {
    const { serviceMember } = this.props;
    if (isEmpty(serviceMember)) {
      console.error(
        'With no service member, you should have been redirected already.',
      );
    }
    this.props.push(`service-member/${serviceMember.id}/create`);
  };

  editMove = move => {
    this.props.push(`moves/${move.id}/edit`);
  };

  resumeMove = () => {
    const { serviceMember, orders, move, ppm } = this.props;
    this.props.push(getNextIncompletePage(serviceMember, orders, move, ppm));
  };
  render() {
    const {
      isLoggedIn,
      loggedInUserIsLoading,
      loggedInUserSuccess,
      loggedInUserError,
      createdServiceMemberError,
      moveSubmitSuccess,
      entitlement,
      serviceMember,
      orders,
      move,
      ppm,
    } = this.props;

    return (
      <div className="usa-grid">
        {loggedInUserIsLoading && <span> Loading... </span>}
        {loggedInUserSuccess && (
          <Fragment>
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
            </div>
            {!isLoggedIn && <LoginButton />}
            {isLoggedIn &&
              !isEmpty(serviceMember) &&
              serviceMember.is_profile_complete && (
                <MoveSummary
                  entitlement={entitlement}
                  profile={serviceMember}
                  orders={orders}
                  move={move}
                  ppm={ppm}
                  editMove={this.editMove}
                  resumeMove={this.resumeMove}
                />
              )}
          </Fragment>
        )}
      </div>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.user.isLoggedIn,
  serviceMember: state.serviceMember.currentServiceMember || {},
  orders: state.orders.currentOrders || {},
  move: state.moves.currentMove || {},
  ppm: state.ppm.currentPpm || {},
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

export default withLastLocation(
  connect(mapStateToProps, mapDispatchToProps)(Landing),
);
