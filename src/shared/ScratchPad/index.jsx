import React, { Component } from 'react';
import PreApprovalPanel from 'shared/PreApprovalRequest/PreApprovalPanel.jsx';

const shipment_accessorials = [
  {
    code: '105D',
    item: 'Unpack Reg Crate',
    location: 'D',
    base_quantity: '  16.7',
    notes: '',
    created_at: '2018-09-24T14:05:38.847Z',
    status: 'SUBMITTED',
  },
  {
    code: '105E',
    item: 'Unpack Reg Crate',
    location: 'D',
    base_quantity: '  16.7',
    notes:
      'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
    created_at: '2018-09-24T14:05:38.847Z',
    status: 'APPROVED',
  },
];
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

class ScratchPad extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide panels-body">
        <div className="usa-width-one-whole">
          <div className="usa-width-two-thirds">
            <PreApprovalPanel
              shipment_accessorials={shipment_accessorials}
              accessorials={accessorials}
            />
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
export default ScratchPad;
