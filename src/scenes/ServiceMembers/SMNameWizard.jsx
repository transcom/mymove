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
    const pendingSMNameData = this.props.nameForm.values;
    console.log(pendingSMNameData);
    if (pendingSMNameData) {
      const serviceMember = { pendingSMNameData };
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
      nameForm,
    } = this.props;
    const SMNameData =
      nameForm &&
      nameForm.values &&
      (nameForm.values.first_name && nameForm.values.last_name);
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
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
  nameForm: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  return { ...state.serviceMember, nameForm: state.form.service_member_name };
}
export default connect(mapStateToProps, mapDispatchToProps)(SMNameWizardPage);
