import React from 'react';
import { shallow, mount } from 'enzyme';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

const handleUpdateServiceItems = jest.fn();

const serviceItemWithImg = {
  id: 'abc123',
  submittedAt: '2020-11-20',
  serviceItem: 'Domestic Crating',
  code: 'DCRT',
  details: {
    description: 'grandfather clock',
    imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
    itemDimensions: { length: 7000, width: 2000, height: 3500 },
  },
};

const serviceItemWithDetails = {
  id: 'abc12345',
  submittedAt: '2020-10-15',
  serviceItem: 'Domestic Origin 1st Day SIT',
  code: 'DOFSIT',
  details: {
    ZIP: '20050',
    Reason: 'Took a detour',
  },
};

const serviceItemWithContact = {
  id: 'abc1234',
  submittedAt: '2020-09-01',
  serviceItem: 'Domestic Destination 1st Day SIT',
  code: 'DDFSIT',
  details: {
    firstCustomerContact: { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15' },
    secondCustomerContact: { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21' },
  },
};

describe('RequestedServiceItemsTable', () => {
  it('show the correct number of service items in the table', () => {
    const serviceItems = [serviceItemWithImg];

    let wrapper = shallow(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
      />,
    );

    expect(wrapper.text().includes('1 item')).toBe(true);

    serviceItems.push(serviceItemWithContact);

    wrapper = shallow(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
      />,
    );
    expect(wrapper.text().includes('2 items')).toBe(true);
  });

  it('displays the service item name and submitted date', () => {
    const serviceItems = [serviceItemWithImg, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
      />,
    );

    expect(wrapper.find('.codeName').at(0).text()).toBe('Domestic Crating');
    expect(wrapper.find('.nameAndDate').at(0).text().includes('20 Nov 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(1).text()).toBe('Domestic Destination 1st Day SIT');
    expect(wrapper.find('.nameAndDate').at(1).text().includes('1 Sep 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(2).text()).toBe('Domestic Origin 1st Day SIT');
    expect(wrapper.find('.nameAndDate').at(2).text().includes('15 Oct 2020')).toBe(true);
  });

  it('displays the approve and reject status buttons', () => {
    const serviceItems = [serviceItemWithImg, serviceItemWithContact, serviceItemWithDetails];
    const wrapper = mount(
      <RequestedServiceItemsTable
        handleUpdateMTOServiceItemStatus={handleUpdateServiceItems}
        serviceItems={serviceItems}
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
});
