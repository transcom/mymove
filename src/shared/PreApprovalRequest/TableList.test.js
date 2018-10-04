import React from 'react';
import { shallow } from 'enzyme';
import TableList from './TableList.jsx';

describe('TableList tests', () => {
  let wrapper, icons;
  const onEdit = jest.fn();
  const shipment_accessorials = [
    {
      code: '105D',
      item: 'Unpack Reg Crate',
      location: 'D',
      base_quantity: ' 16.7',
      notes: '',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
    {
      code: '105E',
      item: 'Unpack Reg Crate',
      location: 'D',
      base_quantity: ' 16.7',
      notes:
        'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
  ];
  describe('When on approval is passed in and status is submitted', () => {
    it('renders without crashing', () => {
      wrapper = shallow(
        <TableList
          shipment_accessorials={shipment_accessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
      const childContainer = wrapper.find('.icon');
      expect(wrapper.find('.accessorial-panel').length).toEqual(1);
      expect(childContainer.length).toBe(6);
    });
  });
  describe('When on approval is NOT passed in and status is SUBMITTED', () => {
    beforeEach(() => {
      wrapper = shallow(
        <TableList
          shipment_accessorials={shipment_accessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const childContainer = wrapper.find('.icon');
      expect(childContainer.length).toBe(4);
    });
  });
  describe('When on approval is passed in and status is APPROVED', () => {
    beforeEach(() => {
      shipment_accessorials[0].status = 'APPROVED';
      shipment_accessorials[1].status = 'APPROVED';
      wrapper = shallow(
        <TableList
          shipment_accessorials={shipment_accessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const childContainer = wrapper.find('.icon');
      expect(childContainer.length).toBe(2);
    });
  });
  describe('When on approval is NOT passed in and status is APPROVED', () => {
    beforeEach(() => {
      shipment_accessorials[0].status = 'APPROVED';
      shipment_accessorials[1].status = 'APPROVED';
      wrapper = shallow(
        <TableList
          shipment_accessorials={shipment_accessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const childContainer = wrapper.find('.icon');
      expect(childContainer.length).toBe(2);
    });
  });
});
