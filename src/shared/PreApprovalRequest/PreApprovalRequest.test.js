import React from 'react';
import { shallow } from 'enzyme';
import { BasicPreApprovalRequest as PreApprovalRequest } from './PreApprovalRequest';

function shipmentLineItem(status) {
  return {
    id: 'sldkjf',
    tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
    location: 'D',
    quantity_1: 167000,
    notes: '',
    created_at: '2018-09-24T14:05:38.847Z',
    status: status || 'SUBMITTED',
  };
}

describe('PreApprovalRequest tests', () => {
  let wrapper;
  let onDelete = jest.fn();
  let onApproval = jest.fn();
  const dummyFn = function() {};
  describe('When on approval is passed in and status is submitted', () => {
    beforeEach(() => {
      onApproval.mockClear();
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentLineItem()}
          isActionable={true}
          isActive={dummyFn}
          onDelete={dummyFn}
          onApproval={onApproval}
        />,
      );
    });
    it('renders without crashing', () => {
      const icons = wrapper.find('.icon');
      expect(wrapper.find('tr').length).toEqual(1);
      expect(icons.length).toBe(3);
    });
    it('it calls onApproval with the correct ID.', () => {
      wrapper.find('[data-test="approve-request"]').simulate('click');
      expect(onApproval.mock.calls[0][0]).toBe(shipmentLineItem().id);
    });
  });
  describe('When on approval is NOT passed in and status is SUBMITTED', () => {
    beforeEach(() => {
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentLineItem()}
          isActionable={true}
          isActive={dummyFn}
          onDelete={dummyFn}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(2);
    });
  });
  describe('When on approval is passed in and status is APPROVED', () => {
    beforeEach(() => {
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentLineItem('APPROVED')}
          isActionable={true}
          isActive={dummyFn}
          onDelete={dummyFn}
          onApproval={dummyFn}
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
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentLineItem('APPROVED')}
          isActionable={true}
          isActive={dummyFn}
          onDelete={dummyFn}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(1);
    });
  });
  describe('When on delete is passed in', () => {
    beforeEach(() => {
      onDelete.mockClear();
      wrapper = shallow(
        <PreApprovalRequest
          shipmentLineItem={shipmentLineItem('APPROVED')}
          isActionable={true}
          isActive={dummyFn}
          onDelete={onDelete}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(1);
    });
    it('it shows a confirmation prompt when delete icon is clicked.', () => {
      wrapper.find('[data-test="delete-request"]').simulate('click');
      const buttons = wrapper.find('button');
      expect(wrapper.find('.delete-confirm').length).toBe(1);
      expect(buttons.length).toBe(2);
    });
    it('it dismisses the delete confirmation when no is clicked.', () => {
      wrapper.find('[data-test="delete-request"]').simulate('click');
      const confirm = wrapper.find('td.delete-confirm').first();
      confirm.find('[data-test="cancel-delete"]').simulate('click');
      expect(wrapper.find('.delete-confirm').length).toBe(0);
    });
    it('it calls the delete callback when yes is clicked.', () => {
      wrapper.find('[data-test="delete-request"]').simulate('click');
      const confirm = wrapper.find('td.delete-confirm').first();
      confirm.find('[data-test="approve-delete"]').simulate('click');
      expect(onDelete.mock.calls.length).toBe(1);
      expect(wrapper.find('.delete-confirm').length).toBe(0);
    });
  });
});
