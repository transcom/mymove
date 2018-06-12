import React from 'react';
import { shallow } from 'enzyme';
import { AppWrapper } from '.';
import Header from 'shared/Header/MyMove';
import Footer from 'shared/Footer';

describe('AppWrapper tests', () => {
  let _wrapper;

  beforeEach(() => {
    _wrapper = shallow(<AppWrapper />);
  });

  it('renders without crashing', () => {
    const appWrapper = _wrapper.find('div');
    expect(appWrapper).toBeDefined();
  });

  it('renders Header component', () => {
    expect(_wrapper.find(Header)).toHaveLength(1);
  });

  it('renders Footer component', () => {
    expect(_wrapper.find(Footer)).toHaveLength(1);
  });
});
