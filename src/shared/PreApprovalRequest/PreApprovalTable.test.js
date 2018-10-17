import React from 'react';
import { shallow } from 'enzyme';
import PreApprovalTable from './PreApprovalTable';

describe('PreApprovalTable tests', () => {
  let wrapper;
  const onEdit = jest.fn();
  const shipmentAccessorials = [
    {
      id: 'sldkjf',
      accessorial: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: 167000,
      notes: '',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
    {
      id: 'sldsdff',
      accessorial: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: 788300,
      notes: 'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
  ];
  describe('When on approval is passed in and status is submitted', () => {
    it('renders without crashing', () => {
      wrapper = shallow(
        <PreApprovalTable
          shipmentAccessorials={shipmentAccessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
      expect(wrapper.find('PreApprovalRequest').length).toEqual(2);
    });
  });
});
