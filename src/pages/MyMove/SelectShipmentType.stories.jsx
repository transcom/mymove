import React from 'react';

import { SelectShipmentType } from './SelectShipmentType';

import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Customer Components / Pages / Select Shipment Type',
  parameters: { layout: 'fullscreen' },
};

const noop = () => {};

const defaultProps = {
  push: noop,
  updateMove: noop,
  loadMTOShipments: noop,
  move: {
    status: 'DRAFT',
  },
  mtoShipments: [],
};

export const Submitted = () => {
  const props = {
    ...defaultProps,
    move: {
      ...defaultProps.move,
      status: 'SUBMITTED',
    },
  };

  return (
    <MockProviders>
      <SelectShipmentType {...props} />
    </MockProviders>
  );
};

export const NoSelectedShipments = () => {
  return (
    <MockProviders>
      <SelectShipmentType {...defaultProps} />
    </MockProviders>
  );
};

export const WithPPMComplete = () => {
  const props = {
    ...defaultProps,
    move: {
      ...defaultProps.move,
      personally_procured_moves: [{}],
    },
  };

  return (
    <MockProviders>
      <SelectShipmentType {...props} />
    </MockProviders>
  );
};

export const WithNTSComplete = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [
      {
        shipmentType: SHIPMENT_OPTIONS.NTS,
      },
    ],
  };

  return (
    <MockProviders>
      <SelectShipmentType {...props} />
    </MockProviders>
  );
};

export const WithNTSRComplete = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [
      {
        shipmentType: SHIPMENT_OPTIONS.NTSR,
      },
    ],
  };

  return (
    <MockProviders>
      <SelectShipmentType {...props} />
    </MockProviders>
  );
};

export const WithNTSAndNTSRComplete = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [
      {
        shipmentType: SHIPMENT_OPTIONS.NTS,
      },
      {
        shipmentType: SHIPMENT_OPTIONS.NTSR,
      },
    ],
  };

  return (
    <MockProviders>
      <SelectShipmentType {...props} />
    </MockProviders>
  );
};
