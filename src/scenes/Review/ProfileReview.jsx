import PropTypes from 'prop-types';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { push } from 'connected-react-router';
import { withContext } from 'shared/AppContext';

import { selectedMoveType, selectedConusStatus, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';

import ServiceMemberSummary from './ServiceMemberSummary';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import scrollToTop from 'shared/scrollToTop';

class ProfileReview extends Component {
  componentDidMount() {
    scrollToTop();
  }
  resumeMove = () => {
    this.props.push(this.getNextIncompletePage());
  };
  getNextIncompletePage = () => {
    const {
      selectedMoveType,
      conusStatus,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      uploads,
      move,
      ppm,
      mtoShipment,
      backupContacts,
      context,
    } = this.props;
    return getNextIncompletePageInternal({
      selectedMoveType,
      conusStatus,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      uploads,
      move,
      ppm,
      mtoShipment,
      backupContacts,
      context,
    });
  };
  render() {
    const { backupContacts, serviceMember, schemaRank, schemaAffiliation, schemaOrdersType } = this.props;
    return (
      <WizardPage
        handleSubmit={this.resumeMove}
        pageList={this.props.pages}
        pageKey={this.props.pageKey}
        pageIsValid={true}
      >
        <h1>Review your Profile</h1>
        <p>Has anything changed since your last move? Please check your info below, especially your Rank.</p>
        <ServiceMemberSummary
          backupContacts={backupContacts}
          serviceMember={serviceMember}
          schemaRank={schemaRank}
          schemaAffiliation={schemaAffiliation}
          schemaOrdersType={schemaOrdersType}
        />
      </WizardPage>
    );
  }
}

ProfileReview.propTypes = {
  context: PropTypes.shape({
    flags: PropTypes.shape({
      hhgFlow: PropTypes.bool,
      ghcFlow: PropTypes.bool,
    }),
  }).isRequired,
};

ProfileReview.propTypes = {
  serviceMember: PropTypes.object,
  context: {
    flags: {
      hhgFlow: false,
      ghcFlow: false,
    },
  },
};

function mapStateToProps(state) {
  return {
    serviceMember: state.serviceMember.currentServiceMember,
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    selectedMoveType: selectedMoveType(state),
    conusStatus: selectedConusStatus(state),
    schemaRank: getInternalSwaggerDefinition(state, 'ServiceMemberRank'),
    schemaOrdersType: getInternalSwaggerDefinition(state, 'OrdersType'),
    schemaAffiliation: getInternalSwaggerDefinition(state, 'Affiliation'),
    backupContacts: state.serviceMember.currentBackupContacts,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}
export default withContext(connect(mapStateToProps, mapDispatchToProps)(ProfileReview));
