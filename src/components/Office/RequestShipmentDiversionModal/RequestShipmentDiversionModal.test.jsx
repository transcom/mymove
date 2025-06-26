import React from 'react';
import { mount } from 'enzyme';
import { act } from 'react-dom/test-utils';

import RequestShipmentDiversionModal from 'components/Office/RequestShipmentDiversionModal/RequestShipmentDiversionModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('RequestShipmentDiversionModal', () => {
  const shipmentInfo = {
    shipmentID: '123456',
    ifMatchEtag: 'string',
    shipmentLocator: '123456-01',
    actualPickupDate: new Date('3 / 16 / 2020'),
  };

  it('renders the component', () => {
    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );
    expect(wrapper.find('RequestShipmentDiversionModal').exists()).toBe(true);
    expect(wrapper.find('ModalTitle').exists()).toBe(true);
    expect(wrapper.find('ModalActions').exists()).toBe(true);
    expect(wrapper.find('ModalClose').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalBackButton"]').exists()).toBe(true);
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
    expect(wrapper.find('button[data-testid="modalSubmitButton"]').exists()).toBe(true);
  });

  it('closes the modal when close icon is clicked', () => {
    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[data-testid="modalCloseButton"]').simulate('click');

    expect(onClose.mock.calls.length).toBe(1);
  });

  it('closes the modal when the cancel button is clicked', () => {
    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[data-testid="modalBackButton"]').simulate('click');

    expect(onClose).toHaveBeenCalled();
  });

  it('disables the submit button on initial render', () => {
    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    expect(wrapper.find('button[data-testid="modalSubmitButton"]').prop('disabled')).toBe(true);
  });

  it('does not call the submit function when submit button is clicked and the reason field is empty', async () => {
    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    wrapper.find('button[data-testid="modalSubmitButton"]').simulate('click');

    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('does not call the submit function when submit button is clicked and the current date is before the actual pickup date', async () => {
    const shipmentInfoWithDate = {
      shipmentID: '123456',
      ifMatchEtag: 'string',
      shipmentLocator: '123456-01',
      actualPickupDate: new Date('6 / 11 / 3024'),
    };

    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfoWithDate} />,
    );
    await act(async () => {
      wrapper
        .find('[data-testid="textInput"]')
        .simulate('change', { target: { name: 'diversionReason', value: 'reasonable reason' } });
    });

    expect(wrapper.text()).toContain('Fields marked with * are required.');
    const label = wrapper.find('label[htmlFor="diversionReason"]');
    expect(label.exists()).toBe(true);
    expect(label.text()).toContain('*');

    wrapper.update();

    wrapper.find('button[data-testid="modalSubmitButton"]').simulate('click');

    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('shows validation error on text input blur event', async () => {
    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    await act(async () => {
      wrapper.find('[data-testid="textInput"]').simulate('blur');
    });

    wrapper.update();
    expect(wrapper.find('[data-testid="errorMessage"]').text()).toEqual('Required');
    expect(wrapper.find('button[data-testid="modalSubmitButton"]').prop('disabled')).toBe(true);
  });

  it('calls the submit function when submit button is clicked and the reason field is not empty', async () => {
    const wrapper = mount(
      <RequestShipmentDiversionModal onSubmit={onSubmit} onClose={onClose} shipmentInfo={shipmentInfo} />,
    );

    await act(async () => {
      wrapper
        .find('[data-testid="textInput"]')
        .simulate('change', { target: { name: 'diversionReason', value: 'reasonable reason' } });
    });

    wrapper.update();

    expect(wrapper.find('[data-testid="errorMessage"]').exists()).toBe(false);
    expect(wrapper.find('button[data-testid="modalSubmitButton"]').prop('disabled')).toBe(false);
    await act(async () => {
      wrapper.find('form').simulate('submit');
    });

    expect(onSubmit).toHaveBeenCalledWith(
      shipmentInfo.id,
      shipmentInfo.eTag,
      shipmentInfo.shipmentLocator,
      'reasonable reason',
    );
  });
});
