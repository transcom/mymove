import React from 'react';
import { MemoryRouter } from 'react-router';

import TXOTabNav from './TXOTabNav';

export default {
  title: 'Components/TXO Tab Navigation',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d23132ee-a6ce-451e-95f9-0a4ef0882ace?mode=design',
    },
  },
  decorators: [
    (Story) => {
      return (
        <MemoryRouter initialEntries={['/']}>
          <div style={{ padding: '20px' }}>
            <Story />
          </div>
        </MemoryRouter>
      );
    },
  ],
};

const basicNavProps = {
  order: {},
  moveCode: 'TESTCO',
};

const moveDetailsAmendedOrders = {
  ...basicNavProps,
  order: {
    uploadedAmendedOrderID: '1234',
  },
};

const moveTaskOrderWithExcessRisk = {
  ...basicNavProps,
  excessWeightRiskCount: 1,
};

export const Default = () => <TXOTabNav {...basicNavProps} />;

export const withMoveDetailsTag = () => <TXOTabNav {...moveDetailsAmendedOrders} />;

export const withMoveTaskOrderTag = () => <TXOTabNav {...moveTaskOrderWithExcessRisk} />;
