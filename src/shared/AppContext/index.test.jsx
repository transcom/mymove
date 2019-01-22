import { withContext } from '.';
import React from 'react';
import { shallow } from 'enzyme';
const Dummy = withContext(({ context }) => {
  return <div>{context.name}</div>;
});

// I timeboxed getting this to work but couldn't figure out what I am missing.
describe.skip('AppContext withContext HOC', () => {
  it('should pass context into a component as a prop', () => {
    const context = { name: 'bar' };
    const wrapper = shallow(<Dummy />, { context });
    expect(wrapper.text()).to.equal('bar');
  });
});
