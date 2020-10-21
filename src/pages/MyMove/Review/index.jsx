import React, { Component } from 'react';
import { arrayOf, bool, string } from 'prop-types';
import { connect } from 'react-redux';

import styles from './Review.module.scss';

import { hasShortHaulError } from 'shared/incentive';
import { no_op as noOp } from 'shared/utils';
import scrollToTop from 'shared/scrollToTop';
import ConnectedWizardPage from 'shared/WizardPage/index';
import ConnectedSummary from 'components/Customer/Review/Summary';
import 'scenes/Review/Review.css';

class Review extends Component {
  componentDidMount() {
    scrollToTop();
  }

  render() {
    const { pages, pageKey, canMoveNext } = this.props;

    return (
      <div className={styles.reviewMoveContainer}>
        <ConnectedWizardPage
          handleSubmit={noOp}
          pageList={pages}
          pageKey={pageKey}
          pageIsValid
          canMoveNext={canMoveNext}
          hideBackBtn
          showFinishLaterBtn
        >
          <div className={`${styles.reviewMoveHeaderContainer} grid-row margin-bottom-3`}>
            <h1 className="tablet:grid-col-10" data-testid="review-move-header">
              Review your details
            </h1>
            <p className="tablet:grid-col-9">
              You’re almost done setting up your move. Double&#8209;check that your information is accurate, then move
              on to the final step.
            </p>
          </div>
          <ConnectedSummary />
        </ConnectedWizardPage>
      </div>
    );
  }
}

Review.propTypes = {
  canMoveNext: bool.isRequired,
  pageKey: string.isRequired,
  pages: arrayOf(string).isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const ppmEstimate = {
    hasEstimateError: state.ppm.hasEstimateError,
    hasEstimateSuccess: state.ppm.hasEstimateSuccess,
    hasEstimateInProgress: state.ppm.hasEstimateInProgress,
    rateEngineError: state.ppm.rateEngineError || null,
    originDutyStationZip: state.serviceMember.currentServiceMember.current_station.address.postal_code,
  };
  return {
    ...ownProps,
    ppmEstimate,
    canMoveNext: !hasShortHaulError(ppmEstimate.rateEngineError),
  };
};

export default connect(mapStateToProps)(Review);
