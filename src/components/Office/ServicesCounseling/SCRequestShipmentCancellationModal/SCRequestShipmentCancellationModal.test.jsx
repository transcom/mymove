import React from 'react';
import { mount } from 'enzyme';

import { SCRequestShipmentCancellationModal } from 'components/Office/ServicesCounseling/SCRequestShipmentCancellationModal/SCRequestShipmentCancellationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('SCRequestShipmentCancellationModal', () => {
  const shipmentID = '123456';

  it('renders the component', () => {
    const wrapper = mount(
      <SCRequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />,
    );
    expect(wrapper.find('SCRequestShipmentCancellationModal').exists()).toBe(true);
    expect(wrapper.find('ModalTitle').exists()).toBe(true);
    expect(wrapper.find('ModalActions').exists()).toBe(true);
    expect(wrapper.find('ModalClose').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalBackButton"]').exists()).toBe(true);
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
  });

  it('closes the modal when close icon is clicked', () => {
    const wrapper = mount(
      <SCRequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />,
    );

    wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the cancel button is clicked', () => {
    const wrapper = mount(
      <SCRequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />,
    );

    wrapper.find('button[data-testid="modalBackButton"]').simulate('click');

    expect(onClose).toHaveBeenCalled();
  });

  it('calls the submit function when submit button is clicked', async () => {
    const wrapper = mount(
      <SCRequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />,
    );

    wrapper.find('button[type="submit"]').simulate('click');

    expect(onSubmit).toHaveBeenCalled();
  });
});
