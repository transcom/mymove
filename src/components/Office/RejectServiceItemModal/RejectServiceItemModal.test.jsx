import React from 'react';
import { mount } from 'enzyme';
import { act } from 'react-dom/test-utils';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import RejectServiceItemModal from 'components/Office/RejectServiceItemModal/RejectServiceItemModal';

describe('RejectServiceItemModal', () => {
  const submittedServiceItem = {
    id: 'abc123',
    mtoShipmentID: 'xyz789',
    code: 'DCRT',
    serviceItem: 'Domestic Crating',
    status: SERVICE_ITEM_STATUS.SUBMITTED,
    createdAt: '2020-10-31T00:00:00',
    details: {
      description: 'grandfather clock',
      imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
      itemDimensions: { length: 7000, width: 2000, height: 3500 },
    },
  };

  it('renders the component', () => {
    const wrapper = mount(
      <RejectServiceItemModal serviceItem={submittedServiceItem} onSubmit={jest.fn()} onClose={jest.fn()} />,
    );
    expect(wrapper.find('RejectServiceItemModal').exists()).toBe(true);
  });

  it('closes the modal when close icon is clicked', () => {
    const onClose = jest.fn();
    const wrapper = mount(
      <RejectServiceItemModal serviceItem={submittedServiceItem} onSubmit={jest.fn()} onClose={onClose} />,
    );

    wrapper.find('button[data-testid="closeRejectServiceItem"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the back button is clicked', () => {
    const onClose = jest.fn();
    const wrapper = mount(
      <RejectServiceItemModal serviceItem={submittedServiceItem} onSubmit={jest.fn()} onClose={onClose} />,
    );

    wrapper.find('button[data-testid="backButton"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('displays the approved timestamp if previously approved', () => {
    const approvedServiceItem = {
      id: 'abc123',
      code: 'DCRT',
      serviceItem: 'Domestic Crating',
      status: SERVICE_ITEM_STATUS.APPROVED,
      createdAt: '2020-08-28T18:20:28.772634',
      approvedAt: '2020-10-31T00:00:00.12345',
      details: {
        description: 'grandfather clock',
        imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
        itemDimensions: { length: 7000, width: 2000, height: 3500 },
      },
    };
    const wrapper = mount(
      <RejectServiceItemModal serviceItem={approvedServiceItem} onSubmit={jest.fn()} onClose={jest.fn()} />,
    );
    expect(wrapper.find('td').at(0).text().includes('31 Oct 2020')).toBe(true);
  });

  it('disables the submit button on initial render', () => {
    const wrapper = mount(
      <RejectServiceItemModal serviceItem={submittedServiceItem} onSubmit={jest.fn()} onClose={jest.fn()} />,
    );

    expect(wrapper.find('button[data-testid="submitButton"]').prop('disabled')).toBe(true);
  });

  it('shows validation error on text input blur event', async () => {
    const wrapper = mount(
      <RejectServiceItemModal serviceItem={submittedServiceItem} onSubmit={jest.fn()} onClose={jest.fn()} />,
    );

    await act(async () => {
      wrapper.find('[data-testid="textInput"]').simulate('blur');
    });

    wrapper.update();
    expect(wrapper.find('[data-testid="errorMessage"]').text()).toEqual('Required');
    expect(wrapper.find('button[data-testid="submitButton"]').prop('disabled')).toBe(true);
  });

  it('enables the submit button when a rejection reason is entered', async () => {
    const wrapper = mount(
      <RejectServiceItemModal serviceItem={submittedServiceItem} onSubmit={jest.fn()} onClose={jest.fn()} />,
    );

    await act(async () => {
      wrapper
        .find('[data-testid="textInput"]')
        .simulate('change', { target: { name: 'rejectionReason', value: 'good reason' } });
    });

    wrapper.update();
    expect(wrapper.find('[data-testid="errorMessage"]').exists()).toBe(false);
    expect(wrapper.find('button[data-testid="submitButton"]').prop('disabled')).toBe(false);
  });

  // onSubmit is not getting called
  it('calls the submit function when submit button is clicked', async () => {
    const onSubmit = jest.fn();

    const wrapper = mount(
      <RejectServiceItemModal serviceItem={submittedServiceItem} onSubmit={onSubmit} onClose={jest.fn()} />,
    );

    await act(async () => {
      wrapper
        .find('input[name="rejectionReason"]')
        .simulate('change', { target: { name: 'rejectionReason', value: 'good reason' } });
    });
    wrapper.update();

    await act(async () => {
      // the submit button doesn't have an onClick listener explicitly attached but the form does
      wrapper.find('form').simulate('submit');
    });

    expect(onSubmit).toHaveBeenCalledWith('abc123', 'xyz789', SERVICE_ITEM_STATUS.REJECTED, 'good reason');
  });
});
