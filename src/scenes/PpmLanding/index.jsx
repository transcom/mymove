import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { withContext } from 'shared/AppContext';
import { PpmSummary } from './PpmSummary';
import {
  selectServiceMemberFromLoggedInUser,
  selectIsProfileComplete,
  selectCurrentOrders,
  selectCurrentMove,
  selectHasCanceledMove,
  selectCurrentPPM,
} from 'store/entities/selectors';
import { updatePPMs } from 'store/entities/actions';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import { getPPMsForMove } from 'services/internalApi';
import { loadMTOShipments } from 'shared/Entities/modules/mtoShipments';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import withRouter from 'utils/routing';

export class PpmLanding extends Component {
  componentDidUpdate(prevProps) {
    if (prevProps.move && prevProps.move.id !== this.props.move.id) {
      this.props.loadMTOShipments(this.props.move.id);
      getPPMsForMove(this.props.move.id).then((response) => this.props.updatePPMs(response));
    }
  }

  editMove = (move) => {
    const {
      router: { navigate },
    } = this.props;
    navigate(`moves/${move.id}/edit`);
  };

  resumeMove = () => {
    const {
      router: { navigate },
    } = this.props;

    const excludeHomePage = true;
    navigate(this.getNextIncompletePage(excludeHomePage));
  };

  reviewProfile = () => {
    const {
      router: { navigate },
    } = this.props;
    navigate('profile-review');
  };

  getNextIncompletePage = (excludeHomePage) => {
    const { selectedMoveType, lastMoveIsCanceled, serviceMember, orders, move, ppm, backupContacts, context } =
      this.props;
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
    const { isProfileComplete, entitlement, serviceMember, orders, move, ppm } = this.props;

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
  router: PropTypes.shape({}),
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
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const move = selectCurrentMove(state) || {};

  const props = {
    lastMoveIsCanceled: selectHasCanceledMove(state),
    isProfileComplete: selectIsProfileComplete(state),
    serviceMember,
    backupContacts: serviceMember?.backup_contacts || [],
    orders: selectCurrentOrders(state) || {},
    move: move,
    ppm: selectCurrentPPM(state) || {},
    entitlement: loadEntitlementsFromState(state),
  };
  return props;
};

const mapDispatchToProps = {
  loadMTOShipments,
  updatePPMs,
};

export default withContext(
  withRouter(
    connect(
      mapStateToProps,
      mapDispatchToProps,
    )(requireCustomerState(PpmLanding, profileStates.BACKUP_CONTACTS_COMPLETE)),
  ),
);
