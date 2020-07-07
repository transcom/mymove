import React from 'react';
import { shallow } from 'enzyme';
import { AppWrapper } from '.';
import Header from 'shared/Header/MyMove';
import Footer from 'shared/Footer';
import SomethingWentWrong from 'shared/SomethingWentWrong';

describe('AppWrapper tests', () => {
  let _wrapper;

  const minProps = {
    context: {
      flags: {
        isHhgFlow: false,
      },
    },
  };

  beforeEach(() => {
    _wrapper = shallow(<AppWrapper {...minProps} />);
  });

  it('renders without crashing or erroring', () => {
    const appWrapper = _wrapper.find('div');
    expect(appWrapper).toBeDefined();
    expect(_wrapper.find(SomethingWentWrong)).toHaveLength(0);
  });

  it('renders Header component', () => {
    expect(_wrapper.find(Header)).toHaveLength(1);
  });

  it('renders Footer component', () => {
    expect(_wrapper.find(Footer)).toHaveLength(1);
  });

  it('renders the fail whale', () => {
    _wrapper.setState({ hasError: true });
    expect(_wrapper.find(SomethingWentWrong)).toHaveLength(1);
  });
});
