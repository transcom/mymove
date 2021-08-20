import React from 'react';
import { mount, shallow } from 'enzyme';
import { ApproveRejectModal } from '.';

// throw errors if component throws warning or errors
let error = console.error;
console.error = function (message) {
  error.apply(console, arguments); // keep default behaviour
  throw message instanceof Error ? message : new Error(message);
};

describe('AcceptRejectModal component test', () => {
  // Positive tests
  it('renders without crashing with required props', () => {
    const wrapper = shallow(<ApproveRejectModal approveBtnOnClick={jest.fn()} rejectBtnOnClick={jest.fn()} />);
    expect(wrapper.find('Fragment')).toHaveLength(1);
  });

  it('renders nothing without crashing with required props and showModal prop false', () => {
    const wrapper = shallow(
      <ApproveRejectModal showModal={false} approveBtnOnClick={jest.fn()} rejectBtnOnClick={jest.fn()} />,
    );
    expect(wrapper.find('Fragment')).toHaveLength(0);
  });

  it('handleApproveClick() is called', () => {
    const mockFn = jest.fn();
    const wrapper = shallow(<ApproveRejectModal approveBtnOnClick={mockFn} rejectBtnOnClick={jest.fn()} />);
    wrapper.find('button[children="Approve"]').simulate('click');
    expect(mockFn).toHaveBeenCalled();
  });

  it('handleRejectionClick() is called', () => {
    const mockFn = jest.fn();
    const wrapper = shallow(<ApproveRejectModal approveBtnOnClick={jest.fn()} rejectBtnOnClick={mockFn} />);
    wrapper.setState({
      showRejectionToggleBtn: false,
      showRejectionInput: true,
      rejectBtnIsDisabled: false,
      rejectionReason: 'rejected',
    });
    wrapper.find('button[children="Reject"]').simulate('click');
    expect(mockFn).toHaveBeenCalled();
  });

  it('handleRejectionChange() is called', () => {
    const wrapper = shallow(<ApproveRejectModal approveBtnOnClick={jest.fn()} rejectBtnOnClick={jest.fn()} />);
    wrapper.setState({
      showRejectionToggleBtn: false,
      showRejectionInput: true,
      rejectBtnIsDisabled: true,
    });
    const spyFn = jest.spyOn(wrapper.instance(), 'handleRejectionChange');
    wrapper.find('input').simulate('change', { target: { value: 'foo' } });
    expect(spyFn).toHaveBeenCalled();
  });

  it('handleRejectionCancelClick() is called', () => {
    const wrapper = shallow(<ApproveRejectModal approveBtnOnClick={jest.fn()} rejectBtnOnClick={jest.fn()} />);
    wrapper.setState({
      showRejectionToggleBtn: false,
      showRejectionInput: true,
      rejectBtnIsDisabled: true,
    });
    const spyFn = jest.spyOn(wrapper.instance(), 'handleRejectionCancelClick');
    wrapper.instance().handleRejectionCancelClick();
    expect(spyFn).toHaveBeenCalled();
  });

  it('handleRejectionToggleClick() is called', () => {
    let wrapper = mount(<ApproveRejectModal approveBtnOnClick={jest.fn()} rejectBtnOnClick={jest.fn()} />);
    const spyFn = jest.spyOn(wrapper.instance(), 'handleRejectionToggleClick');
    wrapper.instance().handleRejectionToggleClick();
    expect(spyFn).toHaveBeenCalled();
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
