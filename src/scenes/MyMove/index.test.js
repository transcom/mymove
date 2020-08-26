import React from 'react';
import { shallow } from 'enzyme';

import { AppWrapper } from './index';

import Header from 'shared/Header/MyMove';
import Footer from 'shared/Footer';
import SomethingWentWrong from 'shared/SomethingWentWrong';

describe('AppWrapper tests', () => {
  let wrapper;

  const minProps = {
    loadInternalSchema: jest.fn(),
    loadUser: jest.fn(),
    context: {
      flags: {
        hhgFlow: false,
      },
    },
  };

  beforeEach(() => {
    wrapper = shallow(<AppWrapper {...minProps} />);
  });

  it('renders without crashing or erroring', () => {
    const appWrapper = wrapper.find('div');
    expect(appWrapper).toBeDefined();
    expect(wrapper.find(SomethingWentWrong)).toHaveLength(0);
  });

  it('renders Header component', () => {
    expect(wrapper.find(Header)).toHaveLength(1);
  });

  it('renders Footer component', () => {
    expect(wrapper.find(Footer)).toHaveLength(1);
  });

  it('fetches initial data', () => {
    expect(minProps.loadUser).toHaveBeenCalled();
    expect(minProps.loadInternalSchema).toHaveBeenCalled();
  });

  it('renders the fail whale', () => {
    wrapper.setState({ hasError: true });
    expect(wrapper.find(SomethingWentWrong)).toHaveLength(1);
  });
});
