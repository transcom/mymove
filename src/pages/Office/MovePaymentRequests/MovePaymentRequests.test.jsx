import React from 'react';
import { mount } from 'enzyme';

import MovePaymentRequests from './MovePaymentRequests';

import { MockProviders } from 'testUtils';

describe('MovePaymentRequests', () => {
  const testMoveId = 'c3952bf1-b689-4fa1-aaef-95f33fa60390';
  const component = mount(
    <MockProviders initialEntries={[`/moves/${testMoveId}/payment-requests`]}>
      <MovePaymentRequests />
    </MockProviders>,
  );

  it('renders without errors', () => {
    expect(component.find('h2').contains('Payment Requests')).toBe(true);
  });
});
