import React from 'react';
import { mount } from 'enzyme';

import { CompletePPMCloseoutConfirmationModal } from 'components/Office/PPM/CompletePPMCloseoutConfirmationModal/CompletePPMCloseoutConfirmationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('CompletePPMCloseoutConfirmationModal', () => {
  it('renders the component', () => {
    const wrapper = mount(<CompletePPMCloseoutConfirmationModal onSubmit={onSubmit} onClose={onClose} />);
    expect(wrapper.find('[data-testid="CompletePPMCloseoutConfirmationModal"]').exists()).toBe(true);
    expect(wrapper.find('ModalTitle').exists()).toBe(true);
    expect(wrapper.find('ModalActions').exists()).toBe(true);
    expect(wrapper.find('ModalClose').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalBackButton"]').exists()).toBe(true);
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
  });

  it('closes the modal when close icon is clicked', () => {
    const wrapper = mount(<CompletePPMCloseoutConfirmationModal onSubmit={onSubmit} onClose={onClose} />);

    wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the cancel button is clicked', () => {
    const wrapper = mount(<CompletePPMCloseoutConfirmationModal onSubmit={onSubmit} onClose={onClose} />);

    wrapper.find('button[data-testid="modalBackButton"]').simulate('click');

    expect(onClose).toHaveBeenCalled();
  });

  it('calls the submit function when submit button is clicked', async () => {
    const wrapper = mount(<CompletePPMCloseoutConfirmationModal onSubmit={onSubmit} onClose={onClose} />);

    wrapper.find('button[type="submit"]').simulate('click');

    expect(onSubmit).toHaveBeenCalled();
  });
});
