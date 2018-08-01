import PropTypes from 'prop-types';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { push } from 'react-router-redux';
import { selectedMoveType, lastMoveIsCanceled } from 'scenes/Moves/ducks';

import Summary from './Summary';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';

class ProfileReview extends Component {
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  resumeMove = () => {
    this.props.push(this.getNextIncompletePage());
  };
  getNextIncompletePage = () => {
    const {
      selectedMoveType,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      move,
      ppm,
      hhg,
      backupContacts,
    } = this.props;
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
    return (
      <WizardPage
        handleSubmit={this.resumeMove}
        pageList={this.props.pages}
        pageKey={this.props.pageKey}
        pageIsValid={true}
      >
        <h1>Review your Profile</h1>
        <p>
          Has anything changed since your last move? Please check your info
          below, especially your Rank.
        </p>
        <Summary />
      </WizardPage>
    );
  }
}

ProfileReview.propTypes = {
  currentServiceMember: PropTypes.object,
};

function mapStateToProps(state) {
  return {
    serviceMember: state.serviceMember.currentServiceMember,
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    selectedMoveType: selectedMoveType(state),
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(ProfileReview);
