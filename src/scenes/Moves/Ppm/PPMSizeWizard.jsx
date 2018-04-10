import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createOrUpdatePpm, loadPpm } from './ducks';
import WizardPage from 'shared/WizardPage';
import PpmSize from './Size';

export class PpmSizeWizardPage extends Component {
  componentDidMount() {
    this.props.loadPpm(this.props.match.params.moveId);
  }
  handleSubmit = () => {
    const { pendingPpmSize, createOrUpdatePpm } = this.props;
    //todo: we should make sure this move matches the redux state
    const moveId = this.props.match.params.moveId;
    if (pendingPpmSize) {
      //don't update a ppm unless the size has changed
      createOrUpdatePpm(moveId, { size: pendingPpmSize });
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
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdatePpm, loadPpm }, dispatch);
}
function mapStateToProps(state) {
  return { ...state.ppm, move: state.submittedMoves };
}
export default connect(mapStateToProps, mapDispatchToProps)(PpmSizeWizardPage);
