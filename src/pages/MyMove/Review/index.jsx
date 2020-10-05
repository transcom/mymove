import React, { Component } from 'react';
import { arrayOf, bool, string } from 'prop-types';
import { connect } from 'react-redux';

import styles from './Review.module.scss';

import { hasShortHaulError } from 'shared/incentive';
import { no_op as noOp } from 'shared/utils';
import scrollToTop from 'shared/scrollToTop';
/* eslint-disable import/no-named-as-default */
import WizardPage from 'shared/WizardPage';
import Summary from 'components/Customer/Review/Summary';
/* eslint-enable import/no-named-as-default */
import 'scenes/Review/Review.css';

class Review extends Component {
  componentDidMount() {
    scrollToTop();
  }

  render() {
    const { pages, pageKey, canMoveNext } = this.props;

    return (
      <div className={styles.reviewMoveContainer}>
        <WizardPage handleSubmit={noOp} pageList={pages} pageKey={pageKey} pageIsValid canMoveNext={canMoveNext}>
          <div className={`${styles.reviewMoveHeaderContainer} grid-row`}>
            <h2 className="tablet:grid-col-10" data-testid="review-move-header">
              Review your details
            </h2>
            <p className="tablet:grid-col-9">
              You’re almost done setting up your move. Double&#8209;check that your information is accurate, then move
              on to the final step.
            </p>
          </div>
          <Summary />
        </WizardPage>
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
