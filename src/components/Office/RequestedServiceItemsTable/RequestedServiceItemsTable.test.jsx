import React from 'react';
import { shallow, mount } from 'enzyme';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

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

describe('RequestedServiceItemsTable', () => {
  it('show the correct number of service items in the table', () => {
    const serviceItems = [serviceItemWithImg];

    let wrapper = shallow(<RequestedServiceItemsTable serviceItems={serviceItems} />);

    expect(wrapper.text().includes('1 item')).toBe(true);

    serviceItems.push(serviceItemWithText);

    wrapper = shallow(<RequestedServiceItemsTable serviceItems={serviceItems} />);
    expect(wrapper.text().includes('2 items')).toBe(true);
  });

  it('displays the service item name and submitted date', () => {
    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(<RequestedServiceItemsTable serviceItems={serviceItems} />);

    expect(wrapper.find('.codeName').at(0).text()).toBe('Domestic Crating');
    expect(wrapper.find('.nameAndDate').at(0).text().includes('20 Nov 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(1).text()).toBe('Domestic origin SIT');
    expect(wrapper.find('.nameAndDate').at(1).text().includes('1 Sep 2020')).toBe(true);

    expect(wrapper.find('.codeName').at(2).text()).toBe('Fuel Surcharge');
    expect(wrapper.find('.nameAndDate').at(2).text().includes('15 Oct 2020')).toBe(true);
  });

  it('shows the service item detail text', () => {
    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(<RequestedServiceItemsTable serviceItems={serviceItems} />);

    expect(wrapper.find('.detailImage').text()).toBe('grandfather clock 7ft x 2ft x 3.5ft');
    expect(wrapper.find('.detail').at(1).text()).toBe('Another service item');
    expect(wrapper.find('.detailType').at(0).text()).toBe('ZIP:');
    expect(wrapper.find('.detailLine').at(0).text().includes('20050')).toBe(true);
    expect(wrapper.find('.detailType').at(1).text()).toBe('Reason:');
    expect(wrapper.find('.detailLine').at(1).text().includes('Took a detour')).toBe(true);
  });

  it('displays the approve and reject status buttons', () => {
    const serviceItems = [serviceItemWithImg, serviceItemWithText, serviceItemWithDetails];
    const wrapper = mount(<RequestedServiceItemsTable serviceItems={serviceItems} />);

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
