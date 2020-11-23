import React, { Component, Fragment } from 'react';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withLastLocation } from 'react-router-last-location';
import { withContext } from 'shared/AppContext';

import { PpmSummary } from './PpmSummary';
import PpmAlert from './PpmAlert';
import { selectedMoveType, lastMoveIsCanceled, updateMove } from 'scenes/Moves/ducks';
import { isProfileComplete } from 'scenes/ServiceMembers/ducks';
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
import { getPPM } from 'scenes/Moves/Ppm/ducks';
import { loadPPMs } from 'shared/Entities/modules/ppms';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectActiveOrLatestOrders, selectUploadsForActiveOrders } from 'shared/Entities/modules/orders';
import { loadMTOShipments, selectMTOShipmentForMTO } from 'shared/Entities/modules/mtoShipments';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';

export class PpmLanding extends Component {
  componentDidMount() {
    // Load user into entities
    const { isLoggedIn, showLoggedInUser } = this.props;
    if (isLoggedIn) {
      showLoggedInUser();
    }

    scrollToTop();
  }

  componentDidUpdate(prevProps) {
    const { serviceMember, loggedInUserSuccess, isProfileComplete } = this.props;
    if (loggedInUserSuccess) {
      if (!isEmpty(serviceMember) && !isProfileComplete) {
        // If the service member exists, but is not complete, redirect to next incomplete page.
        this.resumeMove();
      }
    }
    if (prevProps.move && prevProps.move.id !== this.props.move.id) {
      this.props.loadMTOShipments(this.props.move.id);
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
    const excludeHomePage = true;
    this.props.push(this.getNextIncompletePage(excludeHomePage));
  };

  reviewProfile = () => {
    this.props.push('profile-review');
  };

  getNextIncompletePage = (excludeHomePage) => {
    const {
      selectedMoveType,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      uploads,
      move,
      ppm,
      backupContacts,
      context,
    } = this.props;
    return getNextIncompletePageInternal({
      selectedMoveType,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      uploads,
      move,
      ppm,
      backupContacts,
      context,
      excludeHomePage,
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
      <div className="grid-container">
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
              <PpmSummary
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

PpmLanding.propTypes = {
  context: PropTypes.shape({
    flags: PropTypes.shape({
      hhgFlow: PropTypes.bool,
      ghcFlow: PropTypes.bool,
    }),
  }).isRequired,
};

PpmLanding.defaultProps = {
  context: {
    flags: {
      hhgFlow: false,
      ghcFlow: false,
    },
  },
};

const mapStateToProps = (state) => {
  const user = selectCurrentUser(state);
  const serviceMember = get(state, 'serviceMember.currentServiceMember');
  const move = selectActiveOrLatestMove(state);

  const props = {
    mtoShipment: selectMTOShipmentForMTO(state, get(move, 'id', '')),
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    selectedMoveType: selectedMoveType(state),
    isLoggedIn: user.isLoggedIn,
    isProfileComplete: isProfileComplete(state),
    serviceMember: serviceMember || {},
    backupContacts: state.serviceMember.currentBackupContacts || [],
    orders: selectActiveOrLatestOrders(state),
    uploads: selectUploadsForActiveOrders(state),
    move: move,
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
  return bindActionCreators(
    { push, updateMove, loadPPMs, loadMTOShipments, showLoggedInUser: showLoggedInUserAction },
    dispatch,
  );
}

export default withContext(withLastLocation(connect(mapStateToProps, mapDispatchToProps)(PpmLanding)));
