/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';

import DocsUploaded from '../DocsUploaded/index';
import ShipmentList from '../../../ShipmentList/index';

import Step from '.';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Customer Components / Step',
  component: Step,
  decorators: [
    (Story) => (
      <div className="grid-container">
        <Story />
      </div>
    ),
  ],
};

const Template = (args) => <Step {...args}>{args.children}</Step>;

export const ProfileComplete = Template.bind({});
ProfileComplete.args = {
  complete: true,
  completedHeaderText: 'Profile Complete',
  children: <p>Make sure to keep your personal information up to date during your move</p>,
  editBtnLabel: 'Edit',
  headerText: 'Profile',
  onEditBtnClick: action('edit button clicked'),
  step: '1',
};

export const UploadOrders = Template.bind({});
UploadOrders.args = {
  actionBtnLabel: 'Add orders',
  onActionBtnClick: action('action button clicked'),
  children: <p>Upload photos of each page, or upload a PDF</p>,
  editBtnLabel: 'Edit',
  headerText: 'Upload orders',
  onEditBtnClick: action('edit button clicked'),
  step: '2',
};

export const OrdersUploaded = Template.bind({});
OrdersUploaded.args = {
  complete: true,
  completedHeaderText: 'Orders uploaded',
  children: (
    <DocsUploaded
      files={[
        { filename: 'Screen Shot 2020-09-11 at 12.56.58 PM.png', id: '1' },
        { filename: 'Screen Shot 2020-09-11 at 12.58.12 PM.png', id: '2' },
        { filename: 'orderspage3_20200723.png', id: '3' },
      ]}
    />
  ),
  editBtnLabel: 'Edit',
  headerText: 'Upload orders',
  onEditBtnClick: action('edit button clicked'),
  step: '2',
};

export const OrdersAmended = Template.bind({});
OrdersAmended.args = {
  complete: true,
  completedHeaderText: 'Orders',
  children: (
    <>
      <p>If you receive amended orders:</p>
      <ul>
        <li>Upload the new documents here</li>
        <li>Talk directly with your movers about changes</li>
        <li>The transportation office will update your move info to reflect the new orders</li>
      </ul>
    </>
  ),
  editBtnLabel: 'Upload documents',
  headerText: 'Orders',
  onEditBtnClick: action('edit button clicked'),
  step: '2',
  containerClassName: 'step-amended-orders',
};

export const ShipmentSelection = Template.bind({});
ShipmentSelection.args = {
  complete: false,
  children: (
    <p>
      We&apos;ll collect addresses, dates, and how you want to move your things.
      <br />
      Note: You can change these details later by talking to a move counselor or your movers.
    </p>
  ),
  headerText: 'Set up shipments',
  step: '3',
  actionBtnLabel: 'Set up your shipments',
};

export const Shipments = Template.bind({});
Shipments.args = {
  complete: true,
  completedHeaderText: 'Shipments',
  children: (
    <ShipmentList
      shipments={[
        { id: '0001', shipmentType: SHIPMENT_OPTIONS.HHG },
        { id: '0002', shipmentType: SHIPMENT_OPTIONS.NTS },
        { id: '0003', shipmentType: SHIPMENT_OPTIONS.PPM },
      ]}
      onShipmentClick={action('shipment edit icon clicked')}
      moveSubmitted={false}
    />
  ),
  headerText: 'Shipments',
  step: '3',
  actionBtnLabel: 'Add another shipment',
  onActionBtnClick: action('action button clicked'),
  secondaryBtn: true,
};

export const ConfirmMove = Template.bind({});
ConfirmMove.args = {
  complete: false,
  completedHeaderText: 'Move request confirmed',
  children: <p>Review your move details and sign the legal paperwork, then send the info on to your move counselor</p>,
  headerText: 'Confirm move request',
  step: '4',
  actionBtnLabel: 'Review and submit',
  onActionBtnClick: action('action button clicked'),
};

export const MoveSubmitted = Template.bind({});
MoveSubmitted.args = {
  complete: true,
  completedHeaderText: 'Move request confirmed',
  children: (
    <p>
      Move submitted 03 Nov 2020.
      <br />
      <Button unstyled onClick={action('print button clicked')} style={{ paddingLeft: 0, textDecoration: 'underline' }}>
        Print the legal agreement
      </Button>
    </p>
  ),
  headerText: 'Review your request',
  step: '4',
  actionBtnLabel: 'Review your request',
  onActionBtnClick: action('action button clicked'),
  secondaryBtn: true,
};
