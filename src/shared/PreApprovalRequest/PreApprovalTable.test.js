import React from 'react';
import { shallow } from 'enzyme';
import PreApprovalTable from './PreApprovalTable';

describe('PreApprovalTable tests', () => {
  let wrapper;
  const onEdit = jest.fn();
  const shipmentLineItems = [
    {
      id: 'sldkjf',
      tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: 167000,
      notes: '',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
    {
      id: 'sldsdff',
      tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: 788300,
      notes: 'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
  ];
  describe('When shipmentLineItems exist', () => {
    it('renders without crashing', () => {
      wrapper = shallow(
        <PreApprovalTable
          shipmentLineItems={shipmentLineItems}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
      expect(wrapper.find('PreApprovalRequest').length).toEqual(2);
    });
  });
  describe('When no shipmentLineItems exist', () => {
    it('does not show the table', () => {
      wrapper = shallow(
        <PreApprovalTable
          shipmentLineItems={[]}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
      expect(wrapper.exists('div.pre-approval-panel-table-cont')).toBe(true);
      expect(wrapper.exists('table')).toBe(false);
    });
  });
  describe('When a request is being acted upon', () => {
    it('is the only request that is actionable', () => {
      const onActivation = jest.fn();
      wrapper = shallow(
        <PreApprovalTable
          shipmentLineItems={shipmentLineItems}
          onRequestActivation={onActivation}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
      wrapper.setState({ actionRequestId: shipmentLineItems[0].id });
      const requests = wrapper.find('PreApprovalRequest');
      expect(requests.length).toEqual(2);
      requests.forEach(req => {
        if (req.prop('shipmentLineItem').id === shipmentLineItems[0].id) {
          expect(req.prop('isActionable')).toBe(true);
        } else {
          expect(req.prop('isActionable')).toBe(false);
        }
      });
    });
  });
});
