import React from 'react';
import { mount } from 'enzyme';

import { RequestShipmentCancellationModal } from 'components/Office/RequestShipmentCancellationModal/RequestShipmentCancellationModal';

describe('RequestShipmentCancellationModal', () => {
  const shipmentInfo = {
    shipmentID: '123456',
    moveTaskOrderID: '98765',
    ifMatchEtag: 'string',
  };

  it('renders the component', () => {
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={jest.fn()} onClose={jest.fn()} shipmentInfo={shipmentInfo} />,
    );
    expect(wrapper.find('RequestShipmentCancellationModal').exists()).toBe(true);
    expect(wrapper.find('ModalTitle').exists()).toBe(true);
    expect(wrapper.find('ModalActions').exists()).toBe(true);
    expect(wrapper.find('ModalClose').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalBackButton"]').exists()).toBe(true);
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
  });

  it('closes the modal when close icon is clicked', () => {
    const onClose = jest.fn();
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={jest.fn()} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the cancel button is clicked', () => {
    const onClose = jest.fn();
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={jest.fn()} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[data-testid="modalBackButton"]').simulate('click');

    expect(onClose).toHaveBeenCalled();
  });

  it('calls the submit function when submit putton is clicked', async () => {
    const onSubmit = jest.fn();
    const wrapper = mount(
      <RequestShipmentCancellationModal onSubmit={onSubmit} onClose={jest.fn()} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[type="submit"]').simulate('click');

    expect(onSubmit).toHaveBeenCalled();
  });
});
