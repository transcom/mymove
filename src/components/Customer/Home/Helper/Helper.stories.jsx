/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { text, array } from '@storybook/addon-knobs';

import Helper from './index';

import {
  HelperNeedsOrders,
  HelperNeedsShipment,
  HelperNeedsSubmitMove,
  HelperSubmittedMove,
  HelperAmendedOrders,
} from 'pages/MyMove/Home/HomeHelpers';

const title = 'Next step: Add your orders';
const separator = '\n';

export default {
  title: 'Customer Components / Helper',
};

export const Basic = () => (
  <Helper title={text('Title', title)}>
    <p>
      {text(
        'Description',
        'If you have a hard copy, you can take photos of each page. If you have a PDF, you can upload that',
      )}
    </p>
  </Helper>
);

export const UnorderedList = () => (
  <Helper title={text('Title', title)}>
    <ul>
      {array(
        'Help List',
        ['If you have a hard copy, you can take photos of each page', 'If you have a PDF, you can upload that'],
        separator,
      ).map((helpText) => (
        <li key={helpText}>
          <span>{helpText}</span>
        </li>
      ))}
    </ul>
  </Helper>
);

export const NeedsOrders = () => <HelperNeedsOrders />;
export const NeedsShipment = () => <HelperNeedsShipment />;
export const NeedsSubmitMove = () => <HelperNeedsSubmitMove />;
export const SubmittedMove = () => <HelperSubmittedMove />;
export const AmendedOrders = () => <HelperAmendedOrders />;
