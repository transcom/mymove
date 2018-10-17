import React, { Component } from 'react';
import PreApprovalPanel from 'shared/PreApprovalRequest/PreApprovalPanel.jsx';

const shipment_accessorials = [
  {
    id: 'sldkjf',
    accessorial: { code: '105D', item: 'Reg Shipping' },
    location: 'D',
    quantity_1: 8660000,
    notes: '',
    created_at: '2018-09-24T14:05:38.847Z',
    status: 'SUBMITTED',
  },
  {
    id: 'sldsdff',
    accessorial: { code: '105D', item: 'Reg Shipping' },
    location: 'D',
    quantity_1: 167000,
    notes: 'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
    created_at: '2018-09-24T14:05:38.847Z',
    status: 'APPROVED',
  },
];
const accessorials = [
  {
    id: '23j4u9',
    code: 'F9D',
    item: 'Long Haul',
  },
  {
    id: '2348djfl',
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
            <PreApprovalPanel shipment_accessorials={shipment_accessorials} tariff400ngItems={accessorials} />
          </div>
          <div className="usa-width-one-third">
            <button className="usa-button-primary">Click Me (I do nothing)</button>
          </div>
        </div>
      </div>
    );
  }
}
export default ScratchPad;
