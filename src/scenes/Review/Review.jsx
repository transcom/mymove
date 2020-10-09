import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import Summary from './Summary';
import { connect } from 'react-redux';
import scrollToTop from 'shared/scrollToTop';
import { hasShortHaulError } from 'shared/incentive';
import styles from './Review.module.scss';
import './Review.css';

class Review extends Component {
  componentDidMount() {
    scrollToTop();
  }

  render() {
    const { pages, pageKey } = this.props;

    return (
      <div className="review-move-container">
        <WizardPage
          handleSubmit={no_op}
          pageList={pages}
          pageKey={pageKey}
          pageIsValid={true}
          canMoveNext={this.props.canMoveNext}
          hideBackBtn
          showFinishLaterBtn
        >
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
