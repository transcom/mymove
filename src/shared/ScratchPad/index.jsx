import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PreApprovalRequest from 'shared/PreApprovalRequest';

class ScratchPad extends Component {
  onEdit = () => {};
  onDelete = () => {};
  onApproval = () => {};
  render() {
    const accessorials = [
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
      <div>
        <div className="usa-grid grid-wide panels-body">
          <div className="usa-width-one-whole">
            <div className="usa-width-two-thirds">
              <BasicPanel title={'TEST TITLE'}>
                <PreApprovalRequest
                  accessorials={accessorials}
                  isActionable={true}
                  onEdit={this.onEdit}
                  onDelete={this.onDelete}
                  onApproval={this.onApproval}
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
      </div>
    );
  }
}

export default ScratchPad;
