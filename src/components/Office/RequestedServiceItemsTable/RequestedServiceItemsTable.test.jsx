/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow, mount } from 'enzyme';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

const defaultProps = {
  handleShowRejectionDialog: jest.fn(),
  handleUpdateMTOServiceItemStatus: jest.fn(),
};

const serviceItemWithCrating = {
  id: 'abc123',
  createdAt: '2020-11-20',
  serviceItem: 'Domestic crating',
  code: 'DCRT',
  details: {
    description: 'grandfather clock',
    itemDimensions: { length: 7000, width: 2000, height: 3500 },
  },
};

const serviceItemWithContact = {
  id: 'abc1234',
  createdAt: '2020-09-01',
  serviceItem: 'Domestic destination 1st day SIT',
  code: 'DDFSIT',
  details: {
    firstCustomerContact: { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15' },
    secondCustomerContact: { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21' },
    reason: 'Took a detour',
  },
};

const serviceItemWithDetails = {
  id: 'abc12345',
  createdAt: '2020-10-15',
  serviceItem: 'Domestic origin 1st day SIT',
  code: 'DOFSIT',
  details: {
    pickupPostalCode: '20050',
    reason: 'Took a detour',
  },
};

const testDetails = (wrapper) => {
  expect(wrapper.find('.detailCrating p').text()).toBe('grandfather clock');
  expect(wrapper.find('.detailType').at(0).text()).toBe('Item Dimensions:');
  expect(wrapper.find('.detail dd').at(0).text()).toBe('7"x2"x3.5"');

  expect(wrapper.find('.detailType').at(2).text()).toBe('First Customer Contact:');
  expect(wrapper.find('.detail dd').at(2).text().includes('1200Z')).toBe(true);
  expect(wrapper.find('.detailType').at(3).text()).toBe('First Available Delivery Date:');
  expect(wrapper.find('.detail dd').at(3).text().includes('15 Sep 2020')).toBe(true);

  expect(wrapper.find('.detailType').at(4).text()).toBe('Second Customer Contact:');
  expect(wrapper.find('.detail dd').at(4).text().includes('2300Z')).toBe(true);
  expect(wrapper.find('.detailType').at(5).text()).toBe('Second Available Delivery Date:');
  expect(wrapper.find('.detail dd').at(5).text().includes('21 Sep 2020')).toBe(true);

  expect(wrapper.find('.detailType').at(7).text()).toBe('ZIP:');
  expect(wrapper.find('.detail dd').at(7).text().includes('20050')).toBe(true);
  expect(wrapper.find('.detailType').at(6).text()).toBe('Reason:');
  expect(wrapper.find('.detail dd').at(6).text().includes('Took a detour')).toBe(true);
};

describe('RequestedServiceItemsTable', () => {
  it('shows the correct number of service items in the table', () => {
    const serviceItems = [serviceItemWithCrating];

    let wrapper = shallow(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.text().includes('1 item')).toBe(true);

    serviceItems.push(serviceItemWithContact);

    wrapper = shallow(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );
    expect(wrapper.text().includes('2 items')).toBe(true);
  });

  it('displays the service item name and submitted date', () => {
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.find('.codeName').at(0).text()).toBe('Domestic crating');
    expect(wrapper.find('.nameAndDate').at(0).text().includes('20 Nov 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(1).text()).toBe('Domestic destination 1st day SIT');
    expect(wrapper.find('.nameAndDate').at(1).text().includes('1 Sep 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(2).text()).toBe('Domestic origin 1st day SIT');
    expect(wrapper.find('.nameAndDate').at(2).text().includes('15 Oct 2020')).toBe(true);
  });

  it('shows the service item detail text', () => {
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );
    testDetails(wrapper);
  });

  it('displays the approve and reject status buttons', () => {
    serviceItemWithContact.status = 'SUBMITTED';
    serviceItemWithCrating.status = 'SUBMITTED';
    serviceItemWithDetails.status = 'SUBMITTED';

    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    const acceptButtons = wrapper.find({ 'data-testid': 'acceptButton' });
    expect(acceptButtons.at(0).text().includes('Accept')).toBe(true);
    expect(acceptButtons.at(1).text().includes('Accept')).toBe(true);
    expect(acceptButtons.at(2).text().includes('Accept')).toBe(true);

    const rejectButtons = wrapper.find({ 'data-testid': 'rejectButton' });
    expect(rejectButtons.at(0).text().includes('Reject')).toBe(true);
    expect(rejectButtons.at(1).text().includes('Reject')).toBe(true);
    expect(rejectButtons.at(2).text().includes('Reject')).toBe(true);
  });

  it('shows the service item detail text when approved and shows the reject button', () => {
    serviceItemWithDetails.status = 'APPROVED';
    serviceItemWithCrating.status = 'APPROVED';
    serviceItemWithContact.status = 'APPROVED';
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
      />,
    );

    testDetails(wrapper);
    const rejectTextButton = wrapper.find({ 'data-testid': 'rejectTextButton' });
    expect(rejectTextButton.at(0).text().includes('Reject')).toBe(true);
    expect(rejectTextButton.at(1).text().includes('Reject')).toBe(true);
    expect(rejectTextButton.at(2).text().includes('Reject')).toBe(true);
  });

  it('shows the service item detail text when rejected and shows the approve text button', () => {
    serviceItemWithDetails.status = 'REJECTED';
    serviceItemWithCrating.status = 'REJECTED';
    serviceItemWithContact.status = 'REJECTED';
    const serviceItems = [serviceItemWithCrating, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        {...defaultProps}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
      />,
    );

    testDetails(wrapper);
    const approveTextButton = wrapper.find({ 'data-testid': 'approveTextButton' });
    expect(approveTextButton.at(0).text().includes('Approve')).toBe(true);
    expect(approveTextButton.at(1).text().includes('Approve')).toBe(true);
    expect(approveTextButton.at(2).text().includes('Approve')).toBe(true);
  });
});
