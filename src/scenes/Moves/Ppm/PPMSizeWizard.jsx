import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createOrUpdatePpm, loadPpm } from './ducks';
import WizardPage from 'shared/WizardPage';
import PpmSize from './size';
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
      createOrUpdatePpm(moveId, pendingPpmSize);
    }
  };
  render() {
    const { pages, pageKey, pendingPpmSize, currentPpm } = this.props;
    const ppmSize = pendingPpmSize || (currentPpm && currentPpm.size);
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={ppmSize !== null}
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
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdatePpm, loadPpm }, dispatch);
}
function mapStateToProps(state) {
  return { ...state.ppm, move: state.submittedMoves };
}
export default connect(mapStateToProps, mapDispatchToProps)(PpmSizeWizardPage);
