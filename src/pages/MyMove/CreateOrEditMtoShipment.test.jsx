/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import CreateOrEditMtoShipment from './CreateOrEditMtoShipment';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { history, store } from 'shared/store';

function mockWizardPage(path = '') {
  return {
    match: {
      path,
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
  };
}

const defaultProps = {
  wizardPage: mockWizardPage(),
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
      const wrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.HHG,
        wizardPage: mockWizardPage('/moves/:moveId/hhg-start'),
      });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('EditShipment').length).toBe(0);
    });

    it('or renders only the the EditShipment component', () => {
      const wrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.HHG,
      });
      expect(wrapper.find('EditShipment').length).toBe(1);
      expect(wrapper.find('MtoShipmentForm').length).toBe(0);
    });
  });

  describe('when shipmentType is NTS', () => {
    it('renders only the MtoShipmentForm component', () => {
      const wrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.NTS,
        wizardPage: mockWizardPage('/moves/:moveId/nts-start'),
      });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('EditShipment').length).toBe(0);
    });
  });

  describe('when shipmentType is NTSr', () => {
    it('renders only the NTSDetailsForm component', () => {
      const wrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.NTSR,
        wizardPage: mockWizardPage('/moves/:moveId/ntsr-start'),
      });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('EditShipment').length).toBe(0);
    });
  });
});
