import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PreApprovalRequestForm, {
  formName as PreApprovalRequestFormName,
} from 'shared/PreApprovalRequestForm';
import { submit, isValid, isSubmitting } from 'redux-form';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

class ScratchPad extends Component {
  onSubmit = values => {
    console.log('onSubmit', values);
  };
  render() {
    return (
      <div className="usa-grid">
        <div className="usa-width-one-whole">
          <BasicPanel title={'TEST TITLE'}>
            <PreApprovalRequestForm
              ref={form => (this.formReference = form)}
              onSubmit={this.onSubmit}
            />
            <button
              disabled={!this.props.formEnabled}
              onClick={this.props.submitForm}
            >
              Submit
            </button>
          </BasicPanel>
        </div>
      </div>
    );
  }
}
function mapStateToProps(state) {
  return {
    formEnabled:
      isValid(PreApprovalRequestFormName)(state) &&
      !isSubmitting(PreApprovalRequestFormName)(state),
  };
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(PreApprovalRequestFormName),
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(ScratchPad);
