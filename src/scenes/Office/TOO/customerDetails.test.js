import { shallow } from 'enzyme/build';
import { OfficeWrapper } from 'scenes/Office/index';
import React from 'react';
import CustomerDetails from './customerDetails';
import { detectFlags } from 'shared/featureFlags';

describe('OfficeWrapper', () => {
  it('renders the /too/customer/6ac40a00-e762-4f5f-b08d-3ea72a8e4b63/details route when the too feature flag is enabled ', () => {
    const flags = detectFlags('development', '');
    const officeWrapper = shallow(<OfficeWrapper getCurrentUserInfo={() => {}} context={{ flags: flags }} />);
    expect(
      officeWrapper
        .find('Connect(PrivateRouteContainer)[path="/too/customer/6ac40a00-e762-4f5f-b08d-3ea72a8e4b63/details"]')
        .first()
        .prop('component'),
    ).toBe(CustomerDetails);
  });

  it('does not render the /too/customer/6ac40a00-e762-4f5f-b08d-3ea72a8e4b63/details route when the too feature flag is disabled', () => {
    const flags = detectFlags('production', 'office.move.mil');
    const officeWrapper = shallow(<OfficeWrapper getCurrentUserInfo={() => {}} context={{ flags: flags }} />);
    expect(
      officeWrapper
        .find('Connect(PrivateRouteContainer)[path="/too/customer/6ac40a00-e762-4f5f-b08d-3ea72a8e4b63/details"]')
        .exists(),
    ).toBe(false);
  });
});
