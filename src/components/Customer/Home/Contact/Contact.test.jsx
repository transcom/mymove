/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Contact from '.';

const defaultProps = {
  header: '',
  dutyStationName: '',
  officeType: '',
  telephone: '',
};
function mountFooter(props = defaultProps) {
  return mount(<Contact {...props} />);
}
describe('Contact component', () => {
  it('renders footer with given required props', () => {
    const header = 'Contact Info';
    const dutyStationName = 'Headquarters';
    const officeType = 'Homebase';
    const telephone = '(777) 777-7777';
    const props = {
      header,
      dutyStationName,
      officeType,
      telephone,
    };
    const wrapper = mountFooter(props);
    expect(wrapper.find('h6').text()).toBe(header);
    expect(wrapper.find('strong').text()).toBe(dutyStationName);
    expect(wrapper.find('span').length).toBe(2);
    expect(wrapper.find('span').at(0).text()).toBe(officeType);
    expect(wrapper.find('span').at(1).text()).toBe(telephone);
  });
});
