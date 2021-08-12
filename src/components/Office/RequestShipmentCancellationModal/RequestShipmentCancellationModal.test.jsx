import React from 'react';
import { mount } from 'enzyme';

import RequestShipmentCancellationModal from 'components/Office/RequestShipmentCancellationModal/RequestShipmentCancellationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('RequestShipmentCancellationModal', () => {
  const shipmentInfo = {
    shipmentID: '123456',
    moveTaskOrderID: '98765',
    ifMatchEtag: 'string',
  };

  it('renders the component', () => {
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );
    expect(wrapper.find('RequestShipmentCancellationModal').exists()).toBe(true);
    expect(wrapper.find('ModalTitle').exists()).toBe(true);
    expect(wrapper.find('ModalActions').exists()).toBe(true);
    expect(wrapper.find('ModalClose').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalBackButton"]').exists()).toBe(true);
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
  });

  it('closes the modal when close icon is clicked', () => {
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the cancel button is clicked', () => {
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[data-testid="modalBackButton"]').simulate('click');

    expect(onClose).toHaveBeenCalled();
  });

  it('calls the submit function when submit button is clicked', async () => {
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[type="submit"]').simulate('click');

    expect(onSubmit).toHaveBeenCalled();
  });
});
