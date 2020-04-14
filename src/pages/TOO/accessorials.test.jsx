import React from 'react';
import { shallow } from 'enzyme';
import { GovBanner } from '@trussworks/react-uswds';
import Accessorials from './accessorials';

describe('Accessorials', () => {
  const wrapper = shallow(<Accessorials />);

  it('should render the GovBanner', () => {
    expect(wrapper.find(<GovBanner />).length).toEqual(1);
  });

  it('should render the h1', () => {
    expect(wrapper.text()).toContain('This is where we will put our accessorial components!');
  });
});
