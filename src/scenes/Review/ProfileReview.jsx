import PropTypes from 'prop-types';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { push } from 'react-router-redux';

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
    const { serviceMember, orders, move, ppm, backupContacts } = this.props;
    return getNextIncompletePageInternal(
      serviceMember,
      orders,
      move,
      ppm,
      backupContacts,
    );
  };
  render() {
    return (
      <WizardPage
        handleSubmit={this.resumeMove}
        pageList={[]}
        pageKey={''}
        pageIsValid={true}
      >
        <h1>Profile Review</h1>
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
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(ProfileReview);
