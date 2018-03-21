import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { createPpm } from './ducks';
import WizardPage from 'shared/WizardPage';
import PPMSize from '.';
export class PpmSizeWizardPage extends Component {
  render() {
    const { pages, pageKey, pendingPpmSize, createPpm } = this.props;
    return (
      <WizardPage
        handleSubmit={createPpm}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={pendingPpmSize !== null}
      >
        <PPMSize />
      </WizardPage>
    );
  }
}
PpmSizeWizardPage.propTypes = {
  createPpm: PropTypes.func.isRequired,
  pendingPpmSize: PropTypes.string,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createPpm }, dispatch);
}
function mapStateToProps(state) {
  return state.ppm;
}
export default connect(mapStateToProps, mapDispatchToProps)(PpmSizeWizardPage);
