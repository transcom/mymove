import { isFinite } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createOrUpdatePpm, getRawWeightInfo, isHHGPPMComboMove } from './ducks';
import WizardHeader from '../WizardHeader';
import WizardPage from 'shared/WizardPage';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import ppmBlack from 'shared/icon/ppm-black.svg';
import PpmSize from './Size';

export class PpmSizeWizardPage extends Component {
  handleSubmit = () => {
    const { pendingPpmSize, createOrUpdatePpm, currentPpm } = this.props;
    //todo: we should make sure this move matches the redux state
    const moveId = this.props.match.params.moveId;
    if (pendingPpmSize) {
      let weight = currentPpm.weight_estimate;
      if (currentPpm.size !== pendingPpmSize || (!isFinite(weight) || weight === 0)) {
        weight = 0;
      }

      return createOrUpdatePpm(moveId, {
        size: pendingPpmSize,
        weight_estimate: weight,
      });
    }
  };
  render() {
    const { pages, pageKey, pendingPpmSize, currentPpm, error, isHHGPPMComboMove } = this.props;
    const ppmSize = pendingPpmSize || (currentPpm && currentPpm.size);
    return (
      <div>
        {isHHGPPMComboMove && (
          <WizardHeader
            icon={ppmBlack}
            title="Move Setup"
            right={
              <ProgressTimeline>
                <ProgressTimelineStep name="Move Setup" current />
                <ProgressTimelineStep name="Review" />
              </ProgressTimeline>
            }
          />
        )}
        <WizardPage
          handleSubmit={this.handleSubmit}
          pageList={pages}
          pageKey={pageKey}
          pageIsValid={Boolean(ppmSize)}
          dirty={Boolean(pendingPpmSize)}
          error={error}
        >
          <PpmSize />
        </WizardPage>
      </div>
    );
  }
}
PpmSizeWizardPage.propTypes = {
  createOrUpdatePpm: PropTypes.func.isRequired,
  pendingPpmSize: PropTypes.string,
  currentPpm: PropTypes.shape({ size: PropTypes.string, id: PropTypes.string }),
  error: PropTypes.object,
  weightInfo: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdatePpm }, dispatch);
}
function mapStateToProps(state) {
  return {
    ...state.ppm,
    isHHGPPMComboMove: isHHGPPMComboMove(state),
    move: state.moves,
    weightInfo: getRawWeightInfo(state),
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(PpmSizeWizardPage);
