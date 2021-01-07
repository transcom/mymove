import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import { withLastLocation } from 'react-router-last-location';

import { withContext } from 'shared/AppContext';
import { PpmSummary } from './PpmSummary';
import {
  selectServiceMemberFromLoggedInUser,
  selectIsProfileComplete,
  selectCurrentOrders,
  selectCurrentMove,
  selectHasCanceledMove,
  selectMoveType,
} from 'store/entities/selectors';
import { updatePPMs } from 'store/entities/actions';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { selectCurrentUser, selectGetCurrentUserIsLoading, selectGetCurrentUserIsSuccess } from 'shared/Data/users';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import SignIn from 'shared/User/SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import scrollToTop from 'shared/scrollToTop';
import { getPPM } from 'scenes/Moves/Ppm/ducks';
import { getPPMsForMove } from 'services/internalApi';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { loadMTOShipments } from 'shared/Entities/modules/mtoShipments';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';

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
      if (serviceMember && !isProfileComplete) {
        // If the service member exists, but is not complete, redirect to next incomplete page.
        this.resumeMove();
      }
    }

    if (prevProps.move && prevProps.move.id !== this.props.move.id) {
      this.props.loadMTOShipments(this.props.move.id);
      getPPMsForMove(this.props.move.id).then((response) => this.props.updatePPMs(response));
    }
  }

  startMove = (values) => {
    const { serviceMember } = this.props;
    if (!serviceMember) {
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
      isProfileComplete,
      entitlement,
      serviceMember,
      orders,
      move,
      ppm,
      location,
    } = this.props;

    // early return if loading user
    // TODO - handle this at the top level MyMove/index instead
    if (loggedInUserIsLoading) {
      return (
        <div className="grid-container">
          <LoadingPlaceholder />
        </div>
      );
    }

    // early return if not logged in
    // TODO - handle this at the top level MyMove/index instead, and use a redirect instead
    if (!isLoggedIn && !loggedInUserIsLoading) {
      return (
        <div className="grid-container">
          <SignIn location={location} />
        </div>
      );
    }

    return (
      <div className="grid-container">
        <ConnectedFlashMessage />

        {isProfileComplete && (
          <PpmSummary
            entitlement={entitlement}
            profile={serviceMember}
            orders={orders}
            move={move}
            ppm={ppm}
            editMove={this.editMove}
            resumeMove={this.resumeMove}
            reviewProfile={this.reviewProfile}
          />
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
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const move = selectCurrentMove(state) || {};

  const props = {
    lastMoveIsCanceled: selectHasCanceledMove(state),
    selectedMoveType: selectMoveType(state),
    isLoggedIn: user.isLoggedIn,
    isProfileComplete: selectIsProfileComplete(state),
    serviceMember,
    backupContacts: serviceMember?.backup_contacts || [],
    orders: selectCurrentOrders(state) || {},
    move: move,
    ppm: getPPM(state),
    loggedInUser: user,
    loggedInUserIsLoading: selectGetCurrentUserIsLoading(state),
    loggedInUserSuccess: selectGetCurrentUserIsSuccess(state),
    entitlement: loadEntitlementsFromState(state),
  };
  return props;
};

const mapDispatchToProps = {
  push,
  loadMTOShipments,
  updatePPMs,
  showLoggedInUser: showLoggedInUserAction,
};

export default withContext(withLastLocation(connect(mapStateToProps, mapDispatchToProps)(PpmLanding)));
