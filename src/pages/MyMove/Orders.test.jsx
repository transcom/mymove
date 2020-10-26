import React from 'react';
import { mount } from 'enzyme';

//  import/no-named-as-default
import Orders from './Orders';

import { MockProviders } from 'testUtils';

describe('Orders page', () => {
  it('renders', () => {});
  const initialState = {
    swaggerInternal: {
      spec: {
        definitions: { CreateUpdateOrders: { properties: { orders_type: { enum: [], 'x-display-value': [] } } } },
      },
    },
  };
  const wrapper = mount(
    <MockProviders initialState={initialState} initialEntries={['/']}>
      <Orders pages={[]} pageKey="" updateOrders={jest.fn()} />
    </MockProviders>,
  );
  expect(wrapper.length).toEqual(1);
});
