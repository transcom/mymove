import React from 'react';
import { shallow } from 'enzyme';

import { OfficeWrapper, RenderWithOrWithoutHeader } from './index';

import { Queues } from 'scenes/Office/Queues';
import SomethingWentWrong from 'shared/SomethingWentWrong';

describe('OfficeWrapper tests', () => {
  let wrapper;

  beforeEach(() => {
    wrapper = shallow(<OfficeWrapper getCurrentUserInfo={() => {}} />);
  });

  it('renders without crashing or erroring', () => {
    const officeWrapper = wrapper.find('div');
    expect(officeWrapper).toBeDefined();
    expect(wrapper.find(SomethingWentWrong)).toHaveLength(0);
  });

  it('renders the fail whale', () => {
    wrapper.setState({ hasError: true });
    expect(wrapper.find(SomethingWentWrong)).toHaveLength(1);
  });
});

describe('RenderWithOrWithoutHeader', () => {
  it('renders QueueHeader component', () => {
    const wrapper = shallow(<RenderWithOrWithoutHeader tag="main" component={Queues} withHeader={true} />);
    expect(wrapper.find('QueueHeader').exists()).toBe(true);
  });
  it('renders the component passed to it', () => {
    const wrapper = shallow(<RenderWithOrWithoutHeader tag="main" component={Queues} withHeader={true} />);
    expect(wrapper.find('Queues').exists()).toBe(true);
  });
});

describe('RenderWithOrWithoutHeader', () => {
  it('does not renders QueueHeader component', () => {
    const wrapper = shallow(<RenderWithOrWithoutHeader tag="main" component={Queues} withHeader={false} />);
    expect(wrapper.find('QueueHeader').exists()).toBe(false);
  });
  it('renders the component passed to it', () => {
    const wrapper = shallow(<RenderWithOrWithoutHeader tag="main" component={Queues} withHeader={false} />);
    expect(wrapper.find('Queues').exists()).toBe(true);
  });
});
