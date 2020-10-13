/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import CreateOrEditMtoShipment, { CreateOrEditMtoShipmentComponent } from './CreateOrEditMtoShipment';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { history, store } from 'shared/store';

const defaultProps = {
  wizardPage: {
    pageList: ['page1', 'anotherPage/:foo/:bar'],
    pageKey: 'page1',
    match: { isExact: false, path: '', url: '', params: { moveId: '' } },
    history: {
      goBack: jest.fn(),
      push: jest.fn(),
    },
  },
  showLoggedInUser: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
  loadMTOShipments: jest.fn(),
};

describe('CreateOrEditMtoShipment component', () => {
  describe('when shipmentType is HHG', () => {
    const props = { selectedMoveType: SHIPMENT_OPTIONS.HHG };

    it('can render the HHGDetailsForm component', () => {
      return new Error('not implemented');
    });

    it('can render the HHGDetailsForm component', () => {
      return new Error('not implemented');
    });

    it('can render the EditShipment component', () => {
      return new Error('not implemented');
    });
  });

  describe('when shipmentType is NTS', () => {
    it('can render the NTSDetailsForm component', () => {
      return new Error('not implemented');
    });
  });

  describe('when shipmentType is NTSr', () => {
    it('can render the NTSrDetailsForm component', () => {
      return new Error('not implemented');
    });
  });
});
