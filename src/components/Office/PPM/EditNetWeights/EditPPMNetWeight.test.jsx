import React from 'react';
import { render, screen } from '@testing-library/react';

import EditPPMNetWeight from './EditPPMNetWeight';

import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { ReactQueryWrapper } from 'testUtils';

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  patchWeightTicket: jest.fn(),
}));

const shipments = [
  // moveWeightTotal = 7000
  { primeActualWeight: 1000, reweigh: null, status: 'APPROVED' },
  { primeActualWeight: 2000, reweigh: { weight: 3000 }, status: 'APPROVED' },
  {
    ppmShipment: {
      weightTickets: [
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
      ],
    },
    status: 'APPROVED',
  },
  {
    ppmShipment: {
      weightTickets: [
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
      ],
    },
    status: 'APPROVED',
  },
];

const defaultProps = {
  // moveWeightTotal = 7000
  netWeightRemarks: '',
  weightAllowance: 9000,
  weightTicket: { fullWeight: 1200, emptyWeight: 200 },
  shipments,
};

const excessWeight = {
  // moveWeightTotal = 9500
  ...defaultProps,
  weightTicket: { fullWeight: 2700, emptyWeight: 200 },
  shipments: [
    ...shipments,
    {
      ppmShipment: {
        weightTickets: [createCompleteWeightTicket({}, { fullWeight: 2700, emptyWeight: 200 })],
      },
      status: 'APPROVED',
    },
  ],
};

const reduceWeight = {
  // moveWeightTotal = 10000
  ...defaultProps,
  shipments: [
    ...shipments,
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
  ],
};

describe('EditNetPPMWeight', () => {
  it('renders weights and edit button initially', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...defaultProps} />
      </ReactQueryWrapper>,
    );
    expect(screen.getByText('Net weight')).toBeInTheDocument();
    expect(screen.getByText('| original weight')).toBeInTheDocument();
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
  });
  it('renders correct labels and no calculations if there is no excess weight', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...defaultProps} />
      </ReactQueryWrapper>,
    );
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.queryByText('| to fit within weight allowance')).not.toBeInTheDocument();
    expect(screen.queryByText('| to reduce excess weight')).not.toBeInTheDocument();
  });
  it('renders correct labels and calculations if there is excess weight', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...excessWeight} />
      </ReactQueryWrapper>,
    );
    expect(screen.getByText('| to fit within weight allowance')).toBeInTheDocument();
    expect(screen.getAllByText('2,500 lbs')).toHaveLength(2);
    expect(screen.getByText('-500 lbs')).toBeInTheDocument();
  });
  it('renders correct labels and calculations if there is excess weight and reducing the full weight ticket is necessary', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...reduceWeight} />
      </ReactQueryWrapper>,
    );
    expect(screen.getByText('| to reduce excess weight')).toBeInTheDocument();
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.getByText('-1,000 lbs')).toBeInTheDocument();
  });
  it('renders warning when if move total is higher than weight allowance', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...excessWeight} />
      </ReactQueryWrapper>,
    );
    expect(await screen.findByTestId('warning')).toBeInTheDocument();
  });
});
