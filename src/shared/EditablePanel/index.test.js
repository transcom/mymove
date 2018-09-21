import React from 'react';
import { shallow } from 'enzyme';
import { PanelField } from '.';

describe('EditablePanel tests', () => {
  let wrapper;
  it('PanelField renders without crashing', () => {
    const div = document.createElement('div');
    wrapper = shallow(<PanelField title="test" />, div);
    expect(wrapper.find('.panel-field').length).toEqual(1);
  });

  it('PanelField renders empty span value for optional fields', () => {
    const div = document.createElement('div');
    wrapper = shallow(<PanelField title="test" />, div);
    expect(wrapper.find('.field-value').length).toEqual(1);
    expect(wrapper.find('.field-value').text()).toEqual('');
  });

  it('PanelField renders span value "missing" for required field', () => {
    const div = document.createElement('div');
    wrapper = shallow(<PanelField title="test" required />, div);
    expect(wrapper.find('.field-value').length).toEqual(1);
    expect(wrapper.find('.field-value').text()).toEqual('missing');
  });
});
