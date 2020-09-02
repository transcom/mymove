/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import { history, store } from '../../shared/store';

import EditShipment, { EditShipmentComponent } from './EditShipment';

const defaultProps = {
  match: { isExact: false, path: '', url: '', params: { moveId: '' } },
  history: { goBack: () => {} },
  updateMTOShipment: () => {},
  push: () => {},
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
};
function mountEditShipment(props = defaultProps) {
  return mount(
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <EditShipment {...props} />
      </ConnectedRouter>
    </Provider>,
  );
}
describe('EditShipment component', () => {
  it('renders expected form components', () => {
    const wrapper = mountEditShipment();
    expect(wrapper.find('EditShipment').length).toBe(1);
    expect(wrapper.find('DatePickerInput').length).toBe(2);
    expect(wrapper.find('AddressFields').length).toBe(1);
    expect(wrapper.find('ContactInfoFields').length).toBe(2);
  });

  it('renders second address field when has delivery address', () => {
    const wrapper = mount(<EditShipmentComponent {...defaultProps} />);
    wrapper.setState({ hasDeliveryAddress: true });
    expect(wrapper.find('AddressFields').length).toBe(2);
  });
});
