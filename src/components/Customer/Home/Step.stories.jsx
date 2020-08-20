/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { withKnobs, text, boolean } from '@storybook/addon-knobs';

import Step from './Step';
import DocsUploaded from './DocsUploaded';
import ShipmentList from './ShipmentList';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export const ProfileComplete = () => (
  <div className="grid-container">
    <h3>Profile Complete</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Profile Complete')}
      description={(text('Description'), 'Make sure to keep your personal information up to date during your move')}
      editBtnLabel={text('Edit Button Label', 'Edit')}
      headerText={text('Header Text', 'Profile')}
      onEditBtnClick={() => {}}
      step={text('Step', '1')}
    />
  </div>
);

export const UploadOrders = () => (
  <div className="grid-container">
    <h3>Upload orders</h3>
    <Step
      actionBtnLabel="Add orders"
      onActionBtnClick={() => {}}
      description={(text('Description'), 'Upload photos of each page, or upload a PDF')}
      editBtnLabel={text('Edit Button Label', 'Edit')}
      headerText={text('Header Text', 'Upload orders')}
      onEditBtnClick={() => {}}
      step={text('Step', '2')}
    />
  </div>
);

export const OrdersUploaded = () => (
  <div className="grid-container">
    <h3>Upload orders</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Orders uploaded')}
      description={(text('Description'), 'Upload photos of each page, or upload a PDF')}
      editBtnLabel={text('Edit Button Label', 'Edit')}
      headerText={text('Header Text', 'Upload orders')}
      onEditBtnClick={() => {}}
      step={text('Step', '2')}
    >
      <DocsUploaded
        files={[
          { filename: 'Screen Shot 2020-09-11 at 12.56.58 PM.png' },
          { filename: 'Screen Shot 2020-09-11 at 12.58.12 PM.png' },
          { filename: 'orderspage3_20200723.png' },
        ]}
      />
    </Step>
  </div>
);

export const ShipmentSelection = () => (
  <div className="grid-container">
    <h3>Shipment Selection</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Orders uploaded')}
      description={(text('Description'), 'Upload photos of each page, or upload a PDF')}
      editBtnLabel={text('Edit Button Label', 'Edit')}
      headerText={text('Header Text', 'Upload orders')}
      onEditBtnClick={() => {}}
      step={text('Step', '2')}
      actionBtnLabel="Plan your shipments"
    />
  </div>
);

export const Shipments = () => (
  <div className="grid-container">
    <h3>Shipments</h3>
    <Step
      complete={boolean('Complete', false)}
      completedHeaderText={text('Complete Header Text', 'Shipments')}
      description={
        (text('Description'),
        `Tell us where you're going and when you want to get there. We'll help you set up shipments to make it work`)
      }
      headerText={text('Header Text', 'Upload orders')}
      step={text('Step', '3')}
      actionBtnLabel="Plan your shipments"
      onActionBtnClick={() => {}}
      secondary
    >
      <ShipmentList
        shipments={[
          { id: '#0001', type: SHIPMENT_OPTIONS.HHG },
          { id: '#0002', type: SHIPMENT_OPTIONS.NTS },
          { id: '#0003', type: SHIPMENT_OPTIONS.PPM },
        ]}
        onShipmentClick={() => {}}
      />
    </Step>
  </div>
);

export const ConfirmMove = () => (
  <div className="grid-container">
    <h3>Confirm move</h3>
    <Step
      complete={boolean('Complete', false)}
      completedHeaderText={text('Complete Header Text', 'Shipments')}
      description={
        (text('Description'),
        'Review your move details and sign the legal paperwor , then send the info on to your move counselor')
      }
      headerText={text('Header Text', 'Confirm and move request')}
      step={text('Step', '4')}
      actionBtnLabel="Review and submit"
      actionBtnDisabled
    />
  </div>
);

export default {
  title: 'Customer Components | Step',
  decorators: [withKnobs],
};
