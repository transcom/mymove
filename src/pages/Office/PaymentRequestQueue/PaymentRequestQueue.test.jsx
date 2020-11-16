import React from 'react';
import { mount } from 'enzyme';

import PaymentRequestQueue from './PaymentRequestQueue';

import { MockProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  usePaymentRequestQueueQueries: () => {
    return {
      queueResult: {
        data: [
          {
            age: 0.8477863,
            customer: {
              agency: 'ARMY',
              dodID: '3305957632',
              eTag: 'MjAyMC0xMC0xNVQyMzo0ODozNC41ODQxOTZa',
              email: 'leo_spaceman_sm@example.com',
              first_name: 'Leo',
              id: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
              last_name: 'Spacemen',
              phone: '555-555-5555',
              userID: 'c4d59e2b-bff0-4fce-a31f-26a19b1ad34a',
            },
            departmentIndicator: 'AIR_FORCE',
            id: 'a2c34dba-015f-4f96-a38b-0c0b9272e208',
            locator: 'R993T7',
            moveID: '5d4b25bb-eb04-4c03-9a81-ee0398cb779e',
            originGBLOC: 'LKNQ',
            status: 'PENDING',
            submittedAt: '2020-10-15T23:48:35.420Z',
          },
        ],
        totalCount: 1,
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
}));

describe('PaymentRequestQueue', () => {
  const requiredProps = {
    history: { push: jest.fn() },
  };

  const wrapper = mount(
    <MockProviders initialEntries={['invoicing/queue']}>
      {/* eslint-disable-next-line react/jsx-props-no-spreading */}
      <PaymentRequestQueue {...requiredProps} />
    </MockProviders>,
  );

  it('renders the h1', () => {
    expect(wrapper.find('h1').text()).toBe('Payment requests (1)');
  });

  it('renders the correct column headers', () => {
    expect(wrapper.find('thead tr').length).toBe(2);

    expect(wrapper.find('thead tr th').at(0).text()).toBe('Customer name');
    expect(wrapper.find('thead tr th').at(1).text()).toBe('DoD ID');
    expect(wrapper.find('thead tr th').at(2).text()).toContain('Status');
    expect(wrapper.find('thead tr th').at(3).text()).toBe('Age');
    expect(wrapper.find('thead tr th').at(4).text()).toBe('Submitted');
    expect(wrapper.find('thead tr th').at(5).text()).toBe('Move Code');
    expect(wrapper.find('thead tr th').at(6).text()).toContain('Branch');
    expect(wrapper.find('thead tr th').at(7).text()).toBe('Origin GBLOC');
  });

  it('renders the table with data and expected values', () => {
    expect(wrapper.find('Table').exists()).toBe(true);
    expect(wrapper.find('tbody tr').length).toBe(1);

    expect(wrapper.find('tbody tr td').at(0).text()).toBe('Spacemen, Leo');
    expect(wrapper.find('tbody tr td').at(1).text()).toBe('3305957632');
    expect(wrapper.find('tbody tr td').at(2).text()).toBe('Payment requested');
    expect(wrapper.find('tbody tr td').at(3).text()).toBe('Less than 1 day');
    expect(wrapper.find('tbody tr td').at(4).text()).toBe('15 Oct 2020');
    expect(wrapper.find('tbody tr td').at(5).text()).toBe('R993T7');
    expect(wrapper.find('tbody tr td').at(6).text()).toBe('Army');
    expect(wrapper.find('tbody tr td').at(7).text()).toBe('LKNQ');
  });
});
