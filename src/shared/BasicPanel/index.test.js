import React from 'react';
import { shallow } from 'enzyme';
import BasicPanel from '.';

describe('BasicPanel tests', () => {
  let wrapper;
  it('renders without crashing', () => {
    const div = document.createElement('div');
    wrapper = shallow(
      <BasicPanel title="Test title">
        <div>CHILDREN</div>
      </BasicPanel>,
    );
    expect(wrapper.find('.basic-panel').length).toEqual(1);
  });
});
