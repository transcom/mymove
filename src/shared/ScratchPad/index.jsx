import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PreApprovalRequestForm from 'shared/PreApprovalRequestForm';
import { submit } from 'redux-form';
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
              accessorials={[
                {
                  id: 'sdlfkj',
                  code: 'F9D',
                  item: 'Long Haul',
                },
                {
                  id: 'badfka',
                  code: '19D',
                  item: 'Crate',
                },
              ]}
              ref={form => (this.formReference = form)}
              onSubmit={this.onSubmit}
            />
            <button onClick={this.props.submitForm}>Submit</button>
          </BasicPanel>
        </div>
      </div>
    );
  }
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit('preapproval_request_form'),
    },
    dispatch,
  );
}
export default connect(null, mapDispatchToProps)(ScratchPad);
