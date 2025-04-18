import React from 'react';
import { mount } from 'enzyme';

import { SubmitMoveConfirmationModal } from 'components/Office/SubmitMoveConfirmationModal/SubmitMoveConfirmationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('SubmitMoveConfirmationModal', () => {
  it('renders the component', () => {
    const wrapper = mount(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);
    expect(wrapper.find('SubmitMoveConfirmationModal').exists()).toBe(true);
    expect(wrapper.find('ModalTitle').exists()).toBe(true);
    expect(wrapper.find('ModalActions').exists()).toBe(true);
    expect(wrapper.find('ModalClose').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalCancelButton"]').exists()).toBe(true);
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
    expect(wrapper.containsMatchingElement(<p>You can’t make changes after you submit the move.</p>)).toEqual(true);
  });

  it('renders the component with shipment text', () => {
    const wrapper = mount(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} isShipment />);
    expect(wrapper.find('SubmitMoveConfirmationModal').exists()).toBe(true);
    expect(wrapper.find('ModalTitle').exists()).toBe(true);
    expect(wrapper.find('ModalActions').exists()).toBe(true);
    expect(wrapper.find('ModalClose').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalCancelButton"]').exists()).toBe(true);
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
    expect(wrapper.containsMatchingElement(<p>You can’t make changes after you submit the shipment.</p>)).toEqual(true);
  });

  it('closes the modal when close icon is clicked', () => {
    const wrapper = mount(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);

    wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the cancel button is clicked', () => {
    const wrapper = mount(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);

    wrapper.find('button[data-testid="modalCancelButton"]').simulate('click');

    expect(onClose).toHaveBeenCalled();
  });

  it('calls the submit function when submit button is clicked', async () => {
    const wrapper = mount(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);

    wrapper.find('button[type="submit"]').simulate('click');

    expect(onSubmit).toHaveBeenCalled();
  });
});
