import React from 'react';
import { shallow } from 'enzyme';
import HHGWeightWarning from './HHGWeightWarning';
import Alert from 'shared/Alert';

describe('HHG with too high a weight estimate', function() {
  const shipment = { weight_estimate: 12000, progear_weight_estimate: 300, spouse_progear_weight_estimate: 200 };
  const entitlements = { weight: 10000, pro_gear: 200, pro_gear_spouse: 100, storage_in_transit: 90 };
  const wrapper = shallow(<HHGWeightWarning shipment={shipment} entitlements={entitlements} />);

  it('shows a warning if the estimated weight is too high', function() {
    expect(wrapper.html()).toContain(
      'Your weight estimate of 12,000 is 2,000 lbs over your maximum entitlement of 10,000 lbs.',
    );
  });

  it('shows a warning if the estimated weight is too high', function() {
    expect(wrapper.html()).toContain(
      'Your pro-gear weight estimate of 300 is 100 lbs over your maximum entitlement of 200 lbs.',
    );
  });

  it('shows a warning if the estimated weight is too high', function() {
    expect(wrapper.html()).toContain(
      'Your spouse pro-gear weight estimate of 200 is 100 lbs over your maximum entitlement of 100 lbs.',
    );
  });
});

describe('with valid weights', function() {
  const shipment = { weight_estimate: 1000, progear_weight_estimate: 200, spouse_progear_weight_estimate: 200 };
  const entitlements = { weight: 2000, pro_gear: 300, pro_gear_spouse: 300, storage_in_transit: 90 };
  const wrapper = shallow(<HHGWeightWarning shipment={shipment} entitlements={entitlements} />);

  it('shows no alerts', function() {
    expect(wrapper.containsMatchingElement(Alert)).toEqual(false);
  });
});

describe('with no estimates', function() {
  const shipment = {};
  const entitlements = { weight: 2000, pro_gear: 300, pro_gear_spouse: 300, storage_in_transit: 90 };
  const wrapper = shallow(<HHGWeightWarning shipment={shipment} entitlements={entitlements} />);

  it('shows no alerts', function() {
    expect(wrapper.containsMatchingElement(Alert)).toEqual(false);
  });
});
