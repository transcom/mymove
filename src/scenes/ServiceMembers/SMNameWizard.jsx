import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { updateServiceMember, loadServiceMember } from './ducks';
import WizardPage from 'shared/WizardPage';
import SMName from './SMName';

export class SMNameWizardPage extends Component {
  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    const pendingSMNameData = this.props.pendingSMNameData;
    if (pendingSMNameData) {
      const serviceMember = {
        first_name: pendingSMNameData.first_name,
        middle_initial: pendingSMNameData.middle_initial,
        last_name: pendingSMNameData.last_name,
        suffix: pendingSMNameData.suffix,
      };
      this.props.updateServiceMember(serviceMember);
    }
  };

  render() {
    const {
      pages,
      pageKey,
      pendingSMNameData,
      hasSubmitSuccess,
      error,
    } = this.props;
    const SMNameData = pendingSMNameData;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={Boolean(SMNameData)}
        pageIsDirty={Boolean(pendingSMNameData)}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <SMName />
      </WizardPage>
    );
  }
}
SMNameWizardPage.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  pendingSMNameData: PropTypes.object,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  return { ...state.serviceMember };
}
export default connect(mapStateToProps, mapDispatchToProps)(SMNameWizardPage);
