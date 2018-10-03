import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PreApprovalRequest from 'shared/PreApprovalRequest';
import PreApprovalRequestForm from 'shared/PreApprovalRequestForm';
import { submit } from 'redux-form';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

class ScratchPad extends Component {
  onSubmit = values => {
    console.log('onSubmit', values);
  };
  onEdit = () => {};
  onDelete = () => {};
  onApproval = () => {};

  render() {
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
      <div className="usa-grid">
        <div className="usa-width-one-whole">
          <BasicPanel title={'TEST TITLE'}>
            <PreApprovalRequest
              shipment_accessorials={shipment_accessorials}
              isActionable={true}
              onEdit={this.onEdit}
              onDelete={this.onDelete}
              onApproval={this.onApproval}
            />
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
        <div className="usa-width-one-third">
          <button className="usa-button-primary">
            Click Me (I do nothing)
          </button>
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
