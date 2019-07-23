import React from 'react';
import { shallow } from 'enzyme';
import { OfficeWrapper } from '.';
import QueueHeader from 'shared/Header/Office';
import SomethingWentWrong from 'shared/SomethingWentWrong';

describe('OfficeWrapper tests', () => {
  let _wrapper;

  beforeEach(() => {
    _wrapper = shallow(<OfficeWrapper getCurrentUserInfo={() => {}} />);
  });

  it('renders without crashing or erroring', () => {
    const officeWrapper = _wrapper.find('div');
    expect(officeWrapper).toBeDefined();
    expect(_wrapper.find(SomethingWentWrong)).toHaveLength(0);
  });

  it('renders Queue Header component', () => {
    expect(_wrapper.find(QueueHeader)).toHaveLength(1);
  });

  it('renders the fail whale', () => {
    _wrapper.setState({ hasError: true });
    expect(_wrapper.find(SomethingWentWrong)).toHaveLength(1);
  });
});
