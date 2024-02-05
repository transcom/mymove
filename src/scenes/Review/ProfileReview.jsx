import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import ServiceMemberSummary from './ServiceMemberSummary';

import { withContext } from 'shared/AppContext';
import WizardPage from 'shared/WizardPage';
import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import scrollToTop from 'shared/scrollToTop';
import { selectConusStatus } from 'store/onboarding/selectors';
import { selectServiceMemberFromLoggedInUser, selectHasCanceledMove } from 'store/entities/selectors';
import withRouter from 'utils/routing';

class ProfileReview extends Component {
  componentDidMount() {
    scrollToTop();
  }

  resumeMove = () => {
    const { router } = this.props;
    router.navigate(this.getNextIncompletePage());
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
    const { serviceMember, schemaGrade, schemaAffiliation, schemaOrdersType } = this.props;
    return (
      <WizardPage handleSubmit={this.resumeMove} pageList={this.props.pages} pageKey={this.props.pageKey} pageIsValid>
        <h1>Review your Profile</h1>
        <p>Has anything changed since your last move? Please check your info below, especially your pay grade.</p>
        <ServiceMemberSummary
          serviceMember={serviceMember}
          schemaGrade={schemaGrade}
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
    schemaGrade: getInternalSwaggerDefinition(state, 'OrderPayGrade'),
    schemaOrdersType: getInternalSwaggerDefinition(state, 'OrdersType'),
    schemaAffiliation: getInternalSwaggerDefinition(state, 'Affiliation'),
    backupContacts: serviceMember?.backup_contacts || [],
  };
}

export default withContext(withRouter(connect(mapStateToProps)(ProfileReview)));
