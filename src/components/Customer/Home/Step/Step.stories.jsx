/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { text, boolean } from '@storybook/addon-knobs';

import DocsUploaded from '../DocsUploaded/index';
import ShipmentList from '../ShipmentList/index';

import Step from '.';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export const ProfileComplete = () => (
  <div className="grid-container">
    <h3>Profile Complete</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Profile Complete')}
      editBtnLabel={text('Edit Button Label', 'Edit')}
      headerText={text('Header Text', 'Profile')}
      onEditBtnClick={() => {}}
      step={text('Step', '1')}
    >
      <p>{(text('Description'), 'Make sure to keep your personal information up to date during your move')}</p>
    </Step>
  </div>
);

export const UploadOrders = () => (
  <div className="grid-container">
    <h3>Upload orders</h3>
    <Step
      actionBtnLabel="Add orders"
      onActionBtnClick={() => {}}
      editBtnLabel={text('Edit Button Label', 'Edit')}
      headerText={text('Header Text', 'Upload orders')}
      onEditBtnClick={() => {}}
      step={text('Step', '2')}
    >
      <p>{(text('Description'), 'Upload photos of each page, or upload a PDF')}</p>
    </Step>
  </div>
);

export const OrdersUploaded = () => (
  <div className="grid-container">
    <h3>Upload orders</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Orders uploaded')}
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

export const OrdersAmended = () => (
  <div className="grid-container">
    <h3>Orders</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Orders')}
      editBtnLabel={text('Upload documents Button Label', 'Upload documents')}
      headerText={text('Header Text', 'Orders')}
      onEditBtnClick={() => {}}
      step={text('Step', '2')}
      containerClassName="step-amended-orders"
    >
      <p>If you receive amended orders:</p>
      <ul>
        <li>Upload the new documents here</li>
        <li>Talk directly with your movers about changes</li>
        <li>The transportation office will update your move info to reflect the new orders</li>
      </ul>
    </Step>
  </div>
);

export const ShipmentSelection = () => (
  <div className="grid-container">
    <h3>Shipment Selection</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Orders uploaded')}
      editBtnLabel={text('Edit Button Label', 'Edit')}
      headerText={text('Header Text', 'Upload orders')}
      onEditBtnClick={() => {}}
      step={text('Step', '2')}
      actionBtnLabel="Plan your shipments"
    >
      <p>{(text('Description'), 'Upload photos of each page, or upload a PDF')}</p>
    </Step>
  </div>
);

export const Shipments = () => (
  <div className="grid-container">
    <h3>Shipments</h3>
    <Step
      complete={boolean('Complete', false)}
      completedHeaderText={text('Complete Header Text', 'Shipments')}
      headerText={text('Header Text', 'Upload orders')}
      step={text('Step', '3')}
      actionBtnLabel="Plan your shipments"
      onActionBtnClick={() => {}}
      secondary
    >
      <ShipmentList
        shipments={[
          { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG },
          { id: '0002', shipmentType: SHIPMENT_OPTIONS.NTS },
          { id: '0003', shipmentType: SHIPMENT_OPTIONS.PPM },
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
      completedHeaderText={text('Complete Header Text', 'Move request confirmed')}
      headerText={text('Header Text', 'Confirm move request')}
      step={text('Step', '4')}
      actionBtnLabel="Review and submit"
    >
      <p>
        {
          (text('Description'),
          'Review your move details and sign the legal paperwork, then send the info on to your move counselor')
        }
      </p>
    </Step>
  </div>
);

export const MoveSubmitted = () => (
  <div className="grid-container">
    <h3>Move submitted</h3>
    <Step
      complete={boolean('Complete', true)}
      completedHeaderText={text('Complete Header Text', 'Move request confirmed')}
      headerText={text('Header Text', 'Review your request')}
      step={text('Step', '4')}
      actionBtnLabel="Review your request"
      secondaryBtn
      secondaryBtnStyle={{ boxShadow: 'inset 0 0 0 2px #0050d8' }}
    >
      <p>{(text('Description'), 'Move submitted 03 Nov 2020')}</p>
    </Step>
  </div>
);

export default {
  title: 'Customer Components / Step',
};
