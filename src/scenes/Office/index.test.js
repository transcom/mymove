import React from 'react';
import { shallow } from 'enzyme';
import { OfficeWrapper, RenderWithHeader, RenderWithoutHeader, Queues } from '.';
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

  it('renders the fail whale', () => {
    _wrapper.setState({ hasError: true });
    expect(_wrapper.find(SomethingWentWrong)).toHaveLength(1);
  });
});

describe('RenderWithHeader', () => {
  it('renders QueueHeader component', () => {
    const wrapper = shallow(<RenderWithHeader component={Queues} />);
    expect(wrapper.find('QueueHeader').exists()).toBe(true);
  });
  it('renders the component passed to it', () => {
    const wrapper = shallow(<RenderWithHeader component={Queues} />);
    expect(wrapper.find('Queues').exists()).toBe(true);
  });
});

describe('RenderWithoutHeader', () => {
  it('does not renders QueueHeader component', () => {
    const wrapper = shallow(<RenderWithoutHeader component={Queues} />);
    expect(wrapper.find('QueueHeader').exists()).toBe(false);
  });
  it('renders the component passed to it', () => {
    const wrapper = shallow(<RenderWithoutHeader component={Queues} />);
    expect(wrapper.find('Queues').exists()).toBe(true);
  });
});
