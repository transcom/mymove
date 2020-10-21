/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import CreateOrEditMtoShipment from './CreateOrEditMtoShipment';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { history, store } from 'shared/store';

function getMockMatchProp(path = '') {
  return {
    path,
    isExact: false,
    url: '',
    params: { moveId: 'move123' },
  };
}

const defaultProps = {
  match: {
    path: '',
    isExact: false,
    url: '',
    params: { moveId: 'move123' },
  },
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  history: {
    goBack: jest.fn(),
    push: jest.fn(),
  },
  showLoggedInUser: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
  loadMTOShipments: jest.fn(),
  mtoShipment: {},
};

function mountCreateOrEditMtoShipment(props) {
  return mount(
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <CreateOrEditMtoShipment {...defaultProps} {...props} />
      </ConnectedRouter>
    </Provider>,
  );
}

describe('CreateOrEditMtoShipment component', () => {
  describe('when shipmentType is HHG', () => {
    it('renders only the MtoShipmentForm component', () => {
      const createWrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.HHG,
        match: getMockMatchProp('/moves/:moveId/hhg-start'),
      });
      expect(createWrapper.find('MtoShipmentForm').length).toBe(1);

      const editWrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.HHG,
        match: getMockMatchProp('/moves/:moveId/mto-shipments/:mtoShipmentId/edit'),
      });
      expect(editWrapper.find('MtoShipmentForm').length).toBe(1);
    });
  });

  describe('when shipmentType is NTS', () => {
    it('renders only the MtoShipmentForm component', () => {
      const wrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.NTS,
        match: getMockMatchProp('/moves/:moveId/nts-start'),
      });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('EditShipment').length).toBe(0);
    });
  });

  describe('when shipmentType is NTSr', () => {
    it('renders only the NTSDetailsForm component', () => {
      const wrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.NTSR,
        match: getMockMatchProp('/moves/:moveId/ntsr-start'),
      });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('EditShipment').length).toBe(0);
    });
  });
});
