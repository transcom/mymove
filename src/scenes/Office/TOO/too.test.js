import { shallow } from 'enzyme/build';
import { OfficeWrapper } from 'scenes/Office/index';
import React from 'react';
import TOO from './too';
import { detectFlags } from 'shared/featureFlags';

describe('OfficeWrapper', () => {
  it('renders the /ghc/too route when the too feature flag is enabled ', () => {
    const flags = detectFlags('development', '');
    const officeWrapper = shallow(<OfficeWrapper getCurrentUserInfo={() => {}} context={{ flags: flags }} />);
    expect(
      officeWrapper
        .find('Connect(PrivateRouteContainer)[path="/ghc/too"]')
        .first()
        .prop('component'),
    ).toBe(TOO);
  });

  it('does not render the /ghc/too route when the too feature flag is disabled', () => {
    const flags = detectFlags('production', 'office.move.mil');
    const officeWrapper = shallow(<OfficeWrapper getCurrentUserInfo={() => {}} context={{ flags: flags }} />);
    expect(officeWrapper.find('Connect(PrivateRouteContainer)[path="/ghc/too"]').exists()).toBe(false);
  });
});
