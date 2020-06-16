import React from 'react';
import { shallow } from 'enzyme';

import DataPoint from '.';

describe('DataPoint', () => {
  it('renders empty header and body', () => {
    const wrapper = shallow(<DataPoint />);
    expect(wrapper.text()).toBe('');
  });

  it('renders with header text', () => {
    const header = 'This is a datapoint header.';
    const wrapper = shallow(<DataPoint header={header} />);
    expect(wrapper.find('thead').text()).toContain(header);
  });

  it('renders with body element', () => {
    const bodyText = 'Body test.';
    const BodyElement = () => <>{bodyText}</>;
    const wrapper = shallow(<DataPoint body={<BodyElement />} />);
    expect(wrapper.find(BodyElement).name()).toBe('BodyElement');
    expect(wrapper.find(BodyElement).dive().text()).toContain(bodyText);
  });
});
