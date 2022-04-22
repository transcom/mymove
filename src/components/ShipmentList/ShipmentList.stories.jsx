import React from 'react';

import ShipmentList from './ShipmentList';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Components / ShipmentList',
  component: ShipmentList,
  args: {
    moveSubmitted: false,
  },
  argTypes: {
    onShipmentClick: { action: 'edit shipment clicked' },
    onDeleteClick: { action: 'delete shipment clicked' },
  },
  decorators: [
    (Story) => (
      <div className="grid-container">
        <Story />
      </div>
    ),
  ],
};

const Template = (args) => <ShipmentList {...args} />;

const generateDecorator = (text) => [
  (Story) => (
    <>
      <h3>{text}</h3>
      <Story />
    </>
  ),
];

export const BasicSingle = Template.bind({});

BasicSingle.args = {
  shipments: [{ id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG }],
};
BasicSingle.decorators = generateDecorator('Single Shipment');

export const BasicMultiple = Template.bind({});
BasicMultiple.args = {
  shipments: [
    { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG },
    { id: '0002', shipmentType: SHIPMENT_OPTIONS.NTS },
    { id: '0003', shipmentType: SHIPMENT_OPTIONS.PPM },
    { id: '0004', shipmentType: SHIPMENT_OPTIONS.NTSR },
  ],
};
BasicMultiple.decorators = generateDecorator('Multiple Shipments');

export const WithWeightsSingle = Template.bind({});
WithWeightsSingle.args = {
  shipments: [
    {
      id: '0001',
      shipmentType: SHIPMENT_OPTIONS.HHG,
      calculatedBillableWeight: 4600,
      estimatedWeight: 5000,
      primeEstimatedWeight: 300,
      reweigh: { id: '1236', weight: 200 },
    },
  ],
  showShipmentWeight: true,
};
WithWeightsSingle.decorators = generateDecorator('Single Shipment');

export const WithWeightsMultiple = Template.bind({});
WithWeightsMultiple.args = {
  shipments: [
    { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG, calculatedBillableWeight: 6161, estimatedWeight: 5600 },
    { id: '0002', shipmentType: SHIPMENT_OPTIONS.HHG, calculatedBillableWeight: 3200, reweigh: { id: '1234' } },
    {
      id: '0003',
      shipmentType: SHIPMENT_OPTIONS.HHG,
      calculatedBillableWeight: 3400,
      estimatedWeight: 5000,
      primeEstimatedWeight: 300,
      reweigh: { id: '1236', weight: 200 },
    },
  ],
  showShipmentWeight: true,
};
WithWeightsMultiple.decorators = generateDecorator('Multiple Shipments');

export const WithPpms = Template.bind({});
WithPpms.args = {
  shipments: [
    {
      id: '0001',
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: 'completePPM',
        advanceRequested: false,
      },
    },
    {
      id: '0002',
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: 'completePPM',
        advanceRequested: null,
      },
    },
  ],
};
WithPpms.decorators = generateDecorator('With PPM Shipments');
