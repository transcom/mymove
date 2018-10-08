import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PreApprovalRequestForm, {
  formName as PreApprovalRequestFormName,
} from 'shared/PreApprovalRequestForm';
import { submit, isValid, isSubmitting } from 'redux-form';
import PreApprovalRequest from 'shared/PreApprovalRequest';
import { connect } from 'react-redux';
import Creator from 'shared/PreApprovalRequest/Creator';
import { bindActionCreators } from 'redux';
class ScratchPad extends Component {
  onSubmit = values => {
    return new Promise(function(resolve, reject) {
      // do a thing, possibly async, then…
      setTimeout(function() {
        console.log('onSubmit async', values);
        resolve('success');
      }, 50);
    });
  };
  onEdit = () => {};
  onDelete = () => {};
  onApproval = () => {};
  render() {
    const accessorials = [
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
    ];
    const shipment_accessorials = [
      {
        code: '105D',
        item: 'Unpack Reg Crate',
        location: 'D',
        base_quantity: '	16.7',
        notes: '',
        created_at: '2018-09-24T14:05:38.847Z',
        status: 'SUBMITTED',
      },
      {
        code: '105E',
        item: 'Unpack Reg Crate',
        location: 'D',
        base_quantity: '	16.7',
        notes:
          'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
        created_at: '2018-09-24T14:05:38.847Z',
        status: 'APPROVED',
      },
    ];
    return (
      <div className="usa-grid grid-wide panels-body">
        <div className="usa-width-one-whole">
          <div className="usa-width-two-thirds">
            <BasicPanel title={'TEST TITLE'}>
              <PreApprovalRequest
                shipment_accessorials={shipment_accessorials}
                isActionable={true}
                onEdit={this.onEdit}
                onDelete={this.onDelete}
                onApproval={this.onApproval}
              />
              <PreApprovalRequestForm
                accessorials={accessorials}
                onSubmit={this.onSubmit}
              />
              <button
                disabled={!this.props.formEnabled}
                onClick={this.props.submitForm}
              >
                Submit
              </button>
            </BasicPanel>
            <BasicPanel title="Creator Test">
              <Creator
                accessorials={accessorials}
                savePreApprovalRequest={this.onSubmit}
              />
            </BasicPanel>
          </div>
          <div className="usa-width-one-third">
            <button className="usa-button-primary">
              Click Me (I do nothing)
            </button>
          </div>
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
