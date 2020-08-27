import React from 'react';
import { shallow, mount } from 'enzyme';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

const handleUpdateServiceItems = jest.fn();

const serviceItemWithImg = {
  id: 'abc123',
  submittedAt: '2020-11-20',
  serviceItem: 'Domestic Crating',
  code: 'DCRT',
  details: {
    text: 'grandfather clock 7ft x 2ft x 3.5ft',
    imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
  },
};

const serviceItemWithText = {
  id: 'abc1234',
  submittedAt: '2020-09-01',
  serviceItem: 'Domestic origin SIT',
  code: 'DOMSIT',
  details: {
    text: 'Another service item',
  },
};

const serviceItemWithDetails = {
  id: 'abc1234',
  submittedAt: '2020-10-15',
  serviceItem: 'Fuel Surcharge',
  code: 'FSC',
  details: {
    text: { ZIP: '20050', Reason: 'Took a detour' },
  },
};

const testDetails = (wrapper) => {
  expect(wrapper.find('.detailImage').text()).toBe('grandfather clock 7ft x 2ft x 3.5ft');
  expect(wrapper.find('.detail').at(1).text()).toBe('Another service item');
  expect(wrapper.find('.detailType').at(0).text()).toBe('ZIP:');
  expect(wrapper.find('.detailLine').at(0).text().includes('20050')).toBe(true);
  expect(wrapper.find('.detailType').at(1).text()).toBe('Reason:');
  expect(wrapper.find('.detailLine').at(1).text().includes('Took a detour')).toBe(true);
};

describe('RequestedServiceItemsTable', () => {
  it('show the correct number of service items in the table', () => {
    const serviceItems = [serviceItemWithImg];

    let wrapper = shallow(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.text().includes('1 item')).toBe(true);

    serviceItems.push(serviceItemWithText);

    wrapper = shallow(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );
    expect(wrapper.text().includes('2 items')).toBe(true);
  });

  it('displays the service item name and submitted date', () => {
    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    expect(wrapper.find('.codeName').at(0).text()).toBe('Domestic Crating');
    expect(wrapper.find('.nameAndDate').at(0).text().includes('20 Nov 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(1).text()).toBe('Domestic origin SIT');
    expect(wrapper.find('.nameAndDate').at(1).text().includes('1 Sep 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(2).text()).toBe('Fuel Surcharge');
    expect(wrapper.find('.nameAndDate').at(2).text().includes('15 Oct 2020')).toBe(true);
  });

  it('shows the service item detail text', () => {
    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
        statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      />,
    );

    testDetails(wrapper);
  });

  it('displays the approve and reject status buttons', () => {
    serviceItemWithText.status = 'SUBMITTED';
    serviceItemWithImg.status = 'SUBMITTED';
    serviceItemWithDetails.status = 'SUBMITTED';

    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
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
    serviceItemWithImg.status = 'APPROVED';
    serviceItemWithText.status = 'APPROVED';
    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
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
    serviceItemWithImg.status = 'REJECTED';
    serviceItemWithText.status = 'REJECTED';
    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
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
