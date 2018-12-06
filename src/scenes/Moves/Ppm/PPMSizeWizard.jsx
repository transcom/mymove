import { isFinite } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createOrUpdatePpm, getRawWeightInfo, isHHGPPMComboMove } from './ducks';
import WizardHeader from '../WizardHeader';
import WizardPage from 'shared/WizardPage';
import ppmBlack from 'shared/icon/ppm-black.svg';
import PpmSize from './Size';

export class PpmSizeWizardPage extends Component {
  handleSubmit = () => {
    const { pendingPpmSize, createOrUpdatePpm, weightInfo, currentPpm } = this.props;
    //todo: we should make sure this move matches the redux state
    const moveId = this.props.match.params.moveId;
    if (pendingPpmSize) {
      let weight = currentPpm.weight_estimate;
      if (!isFinite(weight)) {
        // Initialize weight to be mid-range
        // eslint-disable-next-line security/detect-object-injection
        let weightRange = weightInfo[pendingPpmSize];
        weight = weightRange.min + (weightRange.max - weightRange.min) / 2;
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
          <WizardHeader icon={ppmBlack} title="Move Setup" right={<p>status tracker goes here</p>} />
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
