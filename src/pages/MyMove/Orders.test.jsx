import React from 'react';
import { mount } from 'enzyme';

import ConnectedOrders from './Orders';

import { MockProviders } from 'testUtils';

describe('Orders page', () => {
  it('renders', () => {
    const initialState = {
      swaggerInternal: {
        spec: {
          definitions: { CreateUpdateOrders: { properties: { orders_type: { enum: [], 'x-display-value': [] } } } },
        },
      },
      serviceMember: {
        currentServiceMember: {
          id: 'testServiceMember123',
        },
      },
    };
    const wrapper = mount(
      <MockProviders initialState={initialState} initialEntries={['/']}>
        <ConnectedOrders pages={[]} pageKey="" updateOrders={jest.fn()} />
      </MockProviders>,
    );
    expect(wrapper.length).toEqual(1);
  });
});
