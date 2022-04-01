import React from 'react';
import { render, screen } from '@testing-library/react';

import MoveHistory from './MoveHistory';

import { MockProviders } from 'testUtils';

const testMoveLocator = 'DVRS0N';

jest.mock('hooks/queries', () => ({
  useGHCGetMoveHistory: () => {
    return {
      isLoading: false,
      isError: false,
      queueResult: {
        totalCount: 2,
        data: [
          {
            action: 'UPDATE',
            actionTstampClk: '2022-03-09T15:33:38.623Z',
            actionTstampStm: '2022-03-09T15:33:38.622Z',
            actionTstampTx: '2022-03-09T15:33:38.579Z',
            changedValues: { postal_code: '90213', updated_at: '2022-03-08T19:08:44.664709' },
            clientQuery:
              'UPDATE "moves" AS moves SET "available_to_prime_at" = $1, "billable_weights_reviewed_at" = $2, "cancel_reason" = $3, "contractor_id" = $4, "excess_weight_acknowledged_at" = $5, "excess_weight_qualified_at" = $6, "excess_weight_upload_id" = $7, "financial_review_flag" = $8, "financial_review_flag_set_at" = $9, "financial_review_remarks" = $10, "locator" = $11, "orders_id" = $12, "ppm_estimated_weight" = $13, "ppm_type" = $14, "reference_id" = $15, "selected_move_type" = $16, "service_counseling_completed_at" = $17, "show" = $18, "status" = $19, "submitted_at" = $20, "tio_remarks" = $21, "updated_at" = $22 WHERE moves.id = $23',
            eventName: 'updateOrder',
            id: '7ce7c1ac-a1d7-4caf-858c-09674a00f273',
            objectId: 'abe92574-53a8-4026-a75c-45ff9eea9bc6',
            oldValues: {
              city: 'Beverly Hills',
              country: 'US',
              created_at: '2022-02-24T23:45:28.8116',
              id: '8dd3d021-101e-442f-83d7-1b5b91554e5e',
              postal_code: '90215',
              state: 'CA',
              street_address_1: '123 Any Street',
              street_address_2: 'P.O. Box 12345',
              street_address_3: 'c/o Some Person',
              updated_at: '2022-03-08T19:01:46.151732',
            },
            relId: 16592,
            schemaName: 'public',
            tableName: 'moves',
            transactionId: 26993,
          },
          {
            action: 'UPDATE',
            actionTstampClk: '2022-03-08T18:28:58.271Z',
            actionTstampStm: '2022-03-08T18:28:58.220Z',
            actionTstampTx: '2022-03-08T18:28:58.152Z',
            changedValues: { postal_code: '90213', updated_at: '2022-03-08T19:08:44.664709' },
            clientQuery:
              'UPDATE orders\nSET origin_duty_location_id = origin_duty_station_id\nWHERE origin_duty_location_id IS NULL\nAND origin_duty_station_id IS NOT NULL;',
            id: '34752aeb-f658-4afa-b1c0-dcdcb5fb3a73',
            objectId: '18b3725c-529c-4add-811b-7345ece8847f',
            oldValues: {
              city: 'Beverly Hills',
              country: 'US',
              created_at: '2022-02-24T23:45:28.8116',
              id: '8dd3d021-101e-442f-83d7-1b5b91554e5e',
              postal_code: '90215',
              state: 'CA',
              street_address_1: '123 Any Street',
              street_address_2: 'P.O. Box 12345',
              street_address_3: 'c/o Some Person',
              updated_at: '2022-03-08T19:01:46.151732',
            },
            relId: 16879,
            schemaName: 'public',
            tableName: 'orders',
            transactionId: 26950,
          },
        ],
        id: 'abe92574-53a8-4026-a75c-45ff9eea9bc6',
        locator: testMoveLocator,
        referenceId: '5037-3728',
      },
    };
  },
}));

describe('MoveHistory', () => {
  render(
    <MockProviders initialEntries={[`/moves/${testMoveLocator}/history`]}>
      <MoveHistory moveCode={testMoveLocator} />,
    </MockProviders>,
  );

  it('renders the different elements of the Move history tab', () => {
    expect(screen.getByText('Move history (2)')).toBeInTheDocument();
    expect(screen.getByRole('table')).toBeInTheDocument();

    expect(screen.getByTestId('move-history-date-time-0')).toHaveTextContent('09 Mar 22 15:33');
    expect(screen.getByTestId('move-history-event-0')).toHaveTextContent('Updated orders');
    expect(screen.getByTestId('move-history-details-0')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-modified-by-0')).toBeInTheDocument();

    expect(screen.getByTestId('move-history-date-time-1')).toHaveTextContent('08 Mar 22 18:28');
    expect(screen.getByTestId('move-history-event-1')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-details-1')).toBeInTheDocument();
    expect(screen.getByTestId('move-history-modified-by-1')).toBeInTheDocument();
  });
});
