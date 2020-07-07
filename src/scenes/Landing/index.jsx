import React, { Component, Fragment } from 'react';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import { bindActionCreators } from 'redux';
import { withLastLocation } from 'react-router-last-location';

import { MoveSummary } from './MoveSummary';
import PpmAlert from './PpmAlert';
import { selectedMoveType, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { createServiceMember, isProfileComplete } from 'scenes/ServiceMembers/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import {
  selectCurrentUser,
  selectGetCurrentUserIsLoading,
  selectGetCurrentUserIsSuccess,
  selectGetCurrentUserIsError,
} from 'shared/Data/users';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import Alert from 'shared/Alert';
import SignIn from 'shared/User/SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import scrollToTop from 'shared/scrollToTop';
import { updateMove } from 'scenes/Moves/ducks';
import { getPPM } from 'scenes/Moves/Ppm/ducks';
import { loadPPMs } from 'shared/Entities/modules/ppms';
import { selectActiveOrLatestOrders } from 'shared/Entities/modules/orders';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';

export class Landing extends Component {
  componentDidMount() {
    scrollToTop();
  }
  componentDidUpdate(prevProps) {
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
    if (prevProps.move && prevProps.move.id !== this.props.move.id) {
      this.props.loadPPMs(this.props.move.id);
    }
  }
  startMove = (values) => {
    const { serviceMember } = this.props;
    if (isEmpty(serviceMember)) {
      console.error('With no service member, you should have been redirected already.');
    }
    this.props.push(`service-member/${serviceMember.id}/create`);
  };

  editMove = (move) => {
    this.props.push(`moves/${move.id}/edit`);
  };

  resumeMove = () => {
    this.props.push(this.getNextIncompletePage());
  };

  reviewProfile = () => {
    this.props.push('profile-review');
  };

  getNextIncompletePage = () => {
    const { selectedMoveType, lastMoveIsCanceled, serviceMember, orders, move, ppm, backupContacts } = this.props;
    return getNextIncompletePageInternal({
      selectedMoveType,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      move,
      ppm,
      backupContacts,
    });
  };
  render() {
    const {
      isLoggedIn,
      loggedInUserIsLoading,
      loggedInUserSuccess,
      loggedInUserError,
      isProfileComplete,
      createdServiceMemberError,
      moveSubmitSuccess,
      entitlement,
      serviceMember,
      orders,
      move,
      ppm,
      requestPaymentSuccess,
      updateMove,
    } = this.props;
    return (
      <div className="grid-container usa-prose">
        {loggedInUserIsLoading && <LoadingPlaceholder />}
        {!isLoggedIn && !loggedInUserIsLoading && <SignIn location={this.props.location} />}
        {loggedInUserSuccess && (
          <Fragment>
            <div>
              {moveSubmitSuccess && !ppm && (
                <Alert type="success" heading="Success">
                  You've submitted your move
                </Alert>
              )}
              {ppm && moveSubmitSuccess && <PpmAlert heading="Congrats - your move is submitted!" />}
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

            {isLoggedIn && !isEmpty(serviceMember) && isProfileComplete && (
              <MoveSummary
                entitlement={entitlement}
                profile={serviceMember}
                orders={orders}
                move={move}
                ppm={ppm}
                editMove={this.editMove}
                resumeMove={this.resumeMove}
                reviewProfile={this.reviewProfile}
                requestPaymentSuccess={requestPaymentSuccess}
                updateMove={updateMove}
              />
            )}
          </Fragment>
        )}
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const user = selectCurrentUser(state);
  const serviceMember = get(state, 'serviceMember.currentServiceMember');

  const props = {
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    selectedMoveType: selectedMoveType(state),
    isLoggedIn: user.isLoggedIn,
    isProfileComplete: isProfileComplete(state),
    serviceMember: serviceMember || {},
    backupContacts: state.serviceMember.currentBackupContacts || [],
    orders: selectActiveOrLatestOrders(state),
    move: selectActiveOrLatestMove(state),
    ppm: getPPM(state),
    loggedInUser: user,
    loggedInUserIsLoading: selectGetCurrentUserIsLoading(state),
    loggedInUserError: selectGetCurrentUserIsError(state),
    loggedInUserSuccess: selectGetCurrentUserIsSuccess(state),
    createdServiceMemberIsLoading: state.serviceMember.isLoading,
    createdServiceMemberSuccess: state.serviceMember.hasSubmitSuccess,
    createdServiceMemberError: state.serviceMember.error,
    createdServiceMember: state.serviceMember.currentServiceMember,
    moveSubmitSuccess: state.signedCertification.moveSubmitSuccess,
    entitlement: loadEntitlementsFromState(state),
    requestPaymentSuccess: state.ppm.requestPaymentSuccess,
  };
  return props;
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, createServiceMember, updateMove, loadPPMs }, dispatch);
}

export default withLastLocation(connect(mapStateToProps, mapDispatchToProps)(Landing));
