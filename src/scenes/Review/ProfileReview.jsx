import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { push } from 'connected-react-router';

import ServiceMemberSummary from './ServiceMemberSummary';

import { withContext } from 'shared/AppContext';
import WizardPage from 'shared/WizardPage';
import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import scrollToTop from 'shared/scrollToTop';
import { selectConusStatus } from 'store/onboarding/selectors';
import { selectServiceMemberFromLoggedInUser, selectHasCanceledMove } from 'store/entities/selectors';

class ProfileReview extends Component {
  componentDidMount() {
    scrollToTop();
  }
  resumeMove = () => {
    this.props.push(this.getNextIncompletePage());
  };
  getNextIncompletePage = () => {
    const {
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
    const { serviceMember, schemaRank, schemaAffiliation, schemaOrdersType } = this.props;
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
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    serviceMember,
    lastMoveIsCanceled: selectHasCanceledMove(state),
    conusStatus: selectConusStatus(state),
    schemaRank: getInternalSwaggerDefinition(state, 'ServiceMemberRank'),
    schemaOrdersType: getInternalSwaggerDefinition(state, 'OrdersType'),
    schemaAffiliation: getInternalSwaggerDefinition(state, 'Affiliation'),
    backupContacts: serviceMember?.backup_contacts || [],
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}

export default withContext(connect(mapStateToProps, mapDispatchToProps)(ProfileReview));
