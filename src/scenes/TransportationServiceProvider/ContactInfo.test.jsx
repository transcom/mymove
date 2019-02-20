import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { TransportationServiceProviderContactInfo } from './ContactInfo';
import { getTspForShipmentLabel, getTspForShipment } from 'shared/Entities/modules/transportationServiceProviders';
import { getPublicShipment } from 'shared/Entities/modules/shipments';

let store;
const mockStore = configureStore();

const transportationServiceProvider = {
  id: 'f21cfbfa-3735-4166-97fb-bbc069e52637',
  name: 'Best moving company',
  poc_general_phone: '222-333-4444',
};

const shipment = {
  id: '7b7606b8-a6f7-4a4f-b450-7be340a1fa55',
  transportation_service_provider_id: transportationServiceProvider.id,
};

const props = {
  getTspForShipment: getTspForShipment(getTspForShipmentLabel, shipment.id),
  getPublicShipment: getPublicShipment('Shipments.getPublicShipment', shipment.id),
};

describe('ContactInfo tests', () => {
  let wrapper;
  store = mockStore({});
  wrapper = mount(
    <Provider store={store}>
      <TransportationServiceProviderContactInfo
        transportationServiceProvider={transportationServiceProvider}
        {...props}
      />
    </Provider>,
  );

  it('renders the correct information', () => {
    expect(wrapper.contains(transportationServiceProvider.name)).toEqual(true);
    expect(wrapper.contains(transportationServiceProvider.poc_general_phone)).toEqual(true);
  });
});
