import React, { Component, Fragment } from 'react';
import { isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';
import { withLastLocation } from 'react-router-last-location';

import { MoveSummary } from './MoveSummary';
import { isHHGPPMComboMove } from 'scenes/Moves/Ppm/ducks';
import { selectedMoveType, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { getCurrentShipment } from 'shared/UI/ducks';
import { createServiceMember, isProfileComplete } from 'scenes/ServiceMembers/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { loadLoggedInUser } from 'shared/User/ducks';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import Alert from 'shared/Alert';
import SignIn from 'shared/User/SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import scrollToTop from 'shared/scrollToTop';
import { updateMove } from 'scenes/Moves/ducks';
import { getPPM } from 'scenes/Moves/Ppm/ducks';

export class Landing extends Component {
  componentDidMount() {
    scrollToTop();
  }
  componentDidUpdate() {
    const {
      serviceMember,
      createdServiceMemberIsLoading,
      createdServiceMemberError,
      loggedInUserSuccess,
      createServiceMember,
      isProfileComplete,
    } = this.props;

    if (loggedInUserSuccess) {
      if (!createdServiceMemberIsLoading && isEmpty(serviceMember) && !createdServiceMemberError) {
        // Once the logged in user loads, if the service member doesn't
        // exist we need to dispatch creating one, once.
        createServiceMember({});
      } else if (!isEmpty(serviceMember) && !isProfileComplete) {
        // If the service member exists, but is not complete, redirect to next incomplete page.
        this.resumeMove();
      }
    }
  }
  startMove = values => {
    const { serviceMember } = this.props;
    if (isEmpty(serviceMember)) {
      console.error('With no service member, you should have been redirected already.');
    }
    this.props.push(`service-member/${serviceMember.id}/create`);
  };

  editMove = move => {
    this.props.push(`moves/${move.id}/edit`);
  };

  resumeMove = () => {
    this.props.push(this.getNextIncompletePage());
  };

  reviewProfile = () => {
    this.props.push('profile-review');
  };

  addPPMShipment = moveID => {
    this.props.updateMove(moveID, 'HHG_PPM').then(() => {
      this.props.push(`/moves/${moveID}/hhg-ppm-start`);
    });
  };

  getNextIncompletePage = () => {
    const { selectedMoveType, lastMoveIsCanceled, serviceMember, orders, move, ppm, hhg, backupContacts } = this.props;
    return getNextIncompletePageInternal({
      selectedMoveType,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      move,
      ppm,
      hhg,
      backupContacts,
    });
  };
  render() {
    const {
      isLoggedIn,
      loggedInUserIsLoading,
      loggedInUserSuccess,
      loggedInUserError,
      hasSubmitSuccess,
      isProfileComplete,
      isHHGPPMComboMove,
      createdServiceMemberError,
      moveSubmitSuccess,
      entitlement,
      serviceMember,
      orders,
      move,
      ppm,
      currentShipment,
      requestPaymentSuccess,
      updateMove,
    } = this.props;
    return (
      <div className="usa-grid">
        {loggedInUserIsLoading && <LoadingPlaceholder />}
        {!isLoggedIn && <SignIn />}
        {loggedInUserSuccess && (
          <Fragment>
            <div>
              {moveSubmitSuccess && (
                <Alert type="success" heading="Success">
                  You've submitted your move
                </Alert>
              )}
              {isHHGPPMComboMove &&
                hasSubmitSuccess && (
                  <Alert type="success" heading="You've added a PPM shipment">
                    Next, your shipment is awaiting approval and this can take up to 3 business days
                  </Alert>
                )}
              {loggedInUserError && (
                <Alert type="error" heading="An error occurred">
                  There was an error loading your user information.
                </Alert>
              )}
              {createdServiceMemberError && (
                <Alert type="error" heading="An error occurred">
                  There was an error creating your profile information.
                </Alert>
              )}
            </div>

            {isLoggedIn &&
              !isEmpty(serviceMember) &&
              isProfileComplete && (
                <MoveSummary
                  entitlement={entitlement}
                  profile={serviceMember}
                  orders={orders}
                  move={move}
                  ppm={ppm}
                  shipment={currentShipment}
                  editMove={this.editMove}
                  resumeMove={this.resumeMove}
                  reviewProfile={this.reviewProfile}
                  requestPaymentSuccess={requestPaymentSuccess}
                  updateMove={updateMove}
                  addPPMShipment={this.addPPMShipment}
                />
              )}
          </Fragment>
        )}
      </div>
    );
  }
}

const mapStateToProps = state => {
  const shipment = getCurrentShipment(state);
  const props = {
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    selectedMoveType: selectedMoveType(state),
    isLoggedIn: state.user.isLoggedIn,
    isProfileComplete: isProfileComplete(state),
    isHHGPPMComboMove: isHHGPPMComboMove(state),
    serviceMember: state.serviceMember.currentServiceMember || {},
    backupContacts: state.serviceMember.currentBackupContacts || [],
    orders: state.orders.currentOrders || {},
    move: state.moves.currentMove || state.moves.latestMove || {},
    ppm: getPPM(state),
    currentShipment: shipment || {},
    loggedInUser: state.loggedInUser.loggedInUser,
    loggedInUserIsLoading: state.loggedInUser.isLoading,
    loggedInUserError: state.loggedInUser.error,
    loggedInUserSuccess: state.loggedInUser.hasSucceeded,
    createdServiceMemberIsLoading: state.serviceMember.isLoading,
    createdServiceMemberSuccess: state.serviceMember.hasSubmitSuccess,
    createdServiceMemberError: state.serviceMember.error,
    createdServiceMember: state.serviceMember.currentServiceMember,
    moveSubmitSuccess: state.signedCertification.moveSubmitSuccess,
    hasSubmitSuccess: state.signedCertification.hasSubmitSuccess,
    entitlement: loadEntitlementsFromState(state),
    requestPaymentSuccess: state.ppm.requestPaymentSuccess,
  };
  return props;
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, createServiceMember, loadLoggedInUser, updateMove }, dispatch);
}

export default withLastLocation(connect(mapStateToProps, mapDispatchToProps)(Landing));
