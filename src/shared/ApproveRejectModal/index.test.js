import React from 'react';
import { shallow } from 'enzyme';
import { ApproveRejectModal } from '.';

// throw errors if component throws warning or errors
let error = console.error;
console.error = function(message) {
  error.apply(console, arguments); // keep default behaviour
  throw message instanceof Error ? message : new Error(message);
};

describe('AcceptRejectModal component test', () => {
  // Positive tests
  it('renders without crashing with required props', () => {
    const wrapper = shallow(<ApproveRejectModal approveBtnOnClick={jest.fn()} rejectBtnOnClick={jest.fn()} />);
    expect(wrapper.find('Fragment')).toHaveLength(1);
  });

  it('renders nothing without crashing with required props and hideModal prop', () => {
    const wrapper = shallow(
      <ApproveRejectModal hideModal={true} approveBtnOnClick={jest.fn()} rejectBtnOnClick={jest.fn()} />,
    );
    expect(wrapper.find('Fragment')).toHaveLength(0);
  });

  // Negative tests
  it('tries renders but crashes with no required props', () => {
    let error = null;
    try {
      shallow(<ApproveRejectModal />);
    } catch (err) {
      error = err;
    }
    expect(error).toBeTruthy();
  });
});
