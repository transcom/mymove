import React from 'react';
import { shallow } from 'enzyme';
import PreApprovalRequest from './PreApprovalRequest';

describe('PreApprovalRequest tests', () => {
  let wrapper;
  const onEdit = jest.fn();
  const shipmentAccessorial = {
    id: 'sldkjf',
    accessorial: { code: '105D', item: 'Reg Shipping' },
    location: 'D',
    base_quantity: 167000,
    notes: '',
    created_at: '2018-09-24T14:05:38.847Z',
    status: 'SUBMITTED',
  };
  describe('When on approval is passed in and status is submitted', () => {
    it('renders without crashing', () => {
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentAccessorial}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
      const icons = wrapper.find('.icon');
      expect(wrapper.find('tr').length).toEqual(1);
      expect(icons.length).toBe(2);
    });
  });
  describe('When on approval is NOT passed in and status is SUBMITTED', () => {
    beforeEach(() => {
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentAccessorial}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(1);
    });
  });
  describe('When on approval is passed in and status is APPROVED', () => {
    beforeEach(() => {
      shipmentAccessorial.status = 'APPROVED';
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentAccessorial}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(1);
    });
  });
  describe('When on approval is NOT passed in and status is APPROVED', () => {
    beforeEach(() => {
      shipmentAccessorial.status = 'APPROVED';
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentAccessorial}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(1);
    });
  });
});
