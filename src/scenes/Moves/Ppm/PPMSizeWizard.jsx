import { isFinite } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createOrUpdatePpm, getRawWeightInfo } from './ducks';
import WizardPage from 'shared/WizardPage';
import PpmSize from './Size';

export class PpmSizeWizardPage extends Component {
  handleSubmit = () => {
    const {
      pendingPpmSize,
      createOrUpdatePpm,
      weightInfo,
      currentPpm,
    } = this.props;
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

      createOrUpdatePpm(moveId, {
        size: pendingPpmSize,
        weight_estimate: weight,
      });
    }
  };
  render() {
    const {
      pages,
      pageKey,
      pendingPpmSize,
      currentPpm,
      hasSubmitSuccess,
      error,
    } = this.props;
    const ppmSize = pendingPpmSize || (currentPpm && currentPpm.size);
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={Boolean(ppmSize)}
        pageIsDirty={Boolean(pendingPpmSize)}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <PpmSize />
      </WizardPage>
    );
  }
}
PpmSizeWizardPage.propTypes = {
  createOrUpdatePpm: PropTypes.func.isRequired,
  pendingPpmSize: PropTypes.string,
  currentPpm: PropTypes.shape({ size: PropTypes.string, id: PropTypes.string }),
  error: PropTypes.object,
  weightInfo: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdatePpm }, dispatch);
}
function mapStateToProps(state) {
  return {
    ...state.ppm,
    move: state.moves,
    weightInfo: getRawWeightInfo(state),
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(PpmSizeWizardPage);
