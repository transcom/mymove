import React from 'react';
import { shallow } from 'enzyme';
import EntitlementBar from '.';

describe('EntitlementBar', () => {
  const getSummaryHtml = entitlement => {
    const wrapper = shallow(<EntitlementBar entitlement={entitlement} />);
    const text = wrapper
      .find('p')
      .last()
      .html();
    return text;
  };
  it('when the entitlement is missing, there is no text', () => {
    expect(getSummaryHtml()).toEqual('<p></p>');
  });
  it('when the entitlement has spouse pro-gear, it is included in summary', () => {
    expect(
      getSummaryHtml({
        pro_gear: 2000,
        pro_gear_spouse: 500,
        sum: 7000,
        weight: 5000,
        storage_in_transit: 90,
      }),
    ).toEqual(
      '<p>5,000 lbs. + 2,000 lbs. of pro-gear + 500 lbs. of spouse&#x27;s pro-gear = <strong>7,000 lbs.</strong></p>',
    );
  });
  it('when the entitlement does not have spouse pro-gear, it is excluded from summary', () => {
    expect(
      getSummaryHtml({
        pro_gear: 2000,
        pro_gear_spouse: 0,
        sum: 7000,
        weight: 5000,
        storage_in_transit: 90,
      }),
    ).toEqual('<p>5,000 lbs. + 2,000 lbs. of pro-gear = <strong>7,000 lbs.</strong></p>');
  });
});
