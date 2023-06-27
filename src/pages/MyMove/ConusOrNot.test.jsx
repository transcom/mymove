import React from 'react';
import { mount } from 'enzyme';
import { act } from 'react-dom/test-utils';
import { Radio } from '@trussworks/react-uswds';

import { ConusOrNot } from 'pages/MyMove/ConusOrNot';
import { CONUS_STATUS } from 'shared/constants';
import { getFeatureFlagForUser } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getFeatureFlagForUser: jest.fn(),
}));

describe('ConusOrNot', () => {
  const minProps = {
    conusStatus: CONUS_STATUS.CONUS,
    setLocation: () => {},
  };
  it('should render radio buttons', async () => {
    getFeatureFlagForUser.mockResolvedValue({ enabled: true, value: 'enabled' });
    // eslint-disable-next-line react/jsx-props-no-spreading
    const wrapper = mount(<ConusOrNot {...minProps} />);

    // wait for the feature flag hook to be resolved
    await act(async () => {
      await Promise.resolve(wrapper);
      await new Promise((resolve) => {
        setTimeout(resolve, 0);
      });
      wrapper.update();
    });

    expect(wrapper.find(Radio).length).toBe(2);

    // PPM button should be checked on page load
    expect(wrapper.find(Radio).at(0).text()).toContain('CONUS');
    expect(wrapper.find(Radio).at(0).find('.usa-radio__input').html()).toContain('checked');

    // HHG button should be disabled
    expect(wrapper.find(Radio).at(1).text()).toContain('OCONUS');
  });
});
