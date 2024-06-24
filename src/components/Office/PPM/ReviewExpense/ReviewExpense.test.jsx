import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewExpense from './ReviewExpense';

import ppmDocumentStatus from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';
import { MockProviders } from 'testUtils';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  ppmShipmentInfo: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    expectedDepartureDate: '2022-12-02',
    actualMoveDate: '2022-12-06',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
  },
  tripNumber: 1,
  ppmNumber: 1,
  showAllFields: false,
  categoryIndex: 1,
};

const expenseRequiredProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'boxes, tape, bubble wrap',
    paidWithGtcc: false,
    amount: 123456,
  },
};

const storageProps = {
  expense: {
    ...expenseRequiredProps.expense,
    movingExpenseType: expenseTypes.STORAGE,
    sitStartDate: '2022-12-15',
    sitEndDate: '2022-12-25',
  },
};

const rejectedProps = {
  expense: {
    ...expenseRequiredProps.expense,
    reason: 'Rejection reason',
    status: ppmDocumentStatus.REJECTED,
  },
};

const documentSetsProps = {
  documentSets: [
    {
      documentSetType: 'WEIGHT_TICKET',
      documentSet: {
        adjustedNetWeight: null,
        allowableWeight: 3000,
        createdAt: '2024-05-16T18:50:33.689Z',
        eTag: 'MjAyNC0wNS0yMFQxNzo1MjowMS4xNTA3MDFa',
        emptyDocument: {
          id: '67739c19-37ca-4de2-8412-513b3564b70c',
          service_member_id: '3f13faa6-0cec-4842-887b-8a7a2d54797e',
          uploads: [
            {
              bytes: 5291,
              contentType: 'image/png',
              createdAt: '2024-05-16T18:50:45.453Z',
              filename: 'thumbnail_image001.png',
              id: '14bcfa31-8879-4ceb-83b5-d36f17005774',
              status: 'PROCESSING',
              updatedAt: '2024-05-16T18:50:45.453Z',
              url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/14bcfa31-8879-4ceb-83b5-d36f17005774?contentType=image%2Fpng',
            },
          ],
        },
        emptyDocumentId: '67739c19-37ca-4de2-8412-513b3564b70c',
        emptyWeight: 1000,
        fullDocument: {
          id: '219675be-c530-4ea4-85ae-8eee92215849',
          service_member_id: '3f13faa6-0cec-4842-887b-8a7a2d54797e',
          uploads: [
            {
              bytes: 5291,
              contentType: 'image/png',
              createdAt: '2024-05-16T18:50:51.465Z',
              filename: 'thumbnail_image001.png',
              id: 'c1fb5f37-b633-46b4-a827-fffa8f7f4b92',
              status: 'PROCESSING',
              updatedAt: '2024-05-16T18:50:51.465Z',
              url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/c1fb5f37-b633-46b4-a827-fffa8f7f4b92?contentType=image%2Fpng',
            },
          ],
        },
        fullDocumentId: '219675be-c530-4ea4-85ae-8eee92215849',
        fullWeight: 4000,
        id: 'c6014838-ff2c-4803-ab3d-8cea23689c10',
        missingEmptyWeightTicket: false,
        missingFullWeightTicket: false,
        netWeightRemarks: null,
        ownsTrailer: true,
        ppmShipmentId: '09b1d087-f8e5-4f46-b75a-297cff7a4d33',
        proofOfTrailerOwnershipDocument: {
          id: 'c087764e-6095-43f0-afda-49cd486820f7',
          service_member_id: '3f13faa6-0cec-4842-887b-8a7a2d54797e',
          uploads: [
            {
              bytes: 5291,
              contentType: 'image/png',
              createdAt: '2024-05-16T18:51:01.436Z',
              filename: 'thumbnail_image001.png',
              id: 'b850d264-29fb-4723-84f5-4a6488e7f358',
              status: 'PROCESSING',
              updatedAt: '2024-05-16T18:51:01.436Z',
              url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/b850d264-29fb-4723-84f5-4a6488e7f358?contentType=image%2Fpng',
            },
          ],
        },
        proofOfTrailerOwnershipDocumentId: 'c087764e-6095-43f0-afda-49cd486820f7',
        reason: null,
        status: 'APPROVED',
        trailerMeetsCriteria: true,
        updatedAt: '2024-05-20T17:52:01.150Z',
        vehicleDescription: 'test',
      },
      uploads: [
        {
          bytes: 5291,
          contentType: 'image/png',
          createdAt: '2024-05-16T18:50:45.453Z',
          filename: 'thumbnail_image001.png',
          id: '14bcfa31-8879-4ceb-83b5-d36f17005774',
          status: 'PROCESSING',
          updatedAt: '2024-05-16T18:50:45.453Z',
          url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/14bcfa31-8879-4ceb-83b5-d36f17005774?contentType=image%2Fpng',
        },
        {
          bytes: 5291,
          contentType: 'image/png',
          createdAt: '2024-05-16T18:50:51.465Z',
          filename: 'thumbnail_image001.png',
          id: 'c1fb5f37-b633-46b4-a827-fffa8f7f4b92',
          status: 'PROCESSING',
          updatedAt: '2024-05-16T18:50:51.465Z',
          url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/c1fb5f37-b633-46b4-a827-fffa8f7f4b92?contentType=image%2Fpng',
        },
        {
          bytes: 5291,
          contentType: 'image/png',
          createdAt: '2024-05-16T18:51:01.436Z',
          filename: 'thumbnail_image001.png',
          id: 'b850d264-29fb-4723-84f5-4a6488e7f358',
          status: 'PROCESSING',
          updatedAt: '2024-05-16T18:51:01.436Z',
          url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/b850d264-29fb-4723-84f5-4a6488e7f358?contentType=image%2Fpng',
        },
      ],
      tripNumber: 0,
    },
  ],
};

const documentSetIndex = 0;

describe('ReviewExpenseForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...expenseRequiredProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Receipt 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Expense Type')).toBeInTheDocument();
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByLabelText('Amount')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: `Review Packing Materials #1` })).toBeInTheDocument();

      expect(screen.getByText(/Add a review for this/i)).toBeInTheDocument();

      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Exclude')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing expense values', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...expenseRequiredProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );

      await waitFor(() => {
        expect(screen.getByText('Expense Type')).toBeInTheDocument();
      });
      expect(screen.getByText('Packing materials')).toBeInTheDocument();
      expect(screen.getByDisplayValue('boxes, tape, bubble wrap'));
      expect(screen.getByLabelText('Amount')).toHaveDisplayValue('1,234.56');
    });

    it('shows SIT fields when expense type is Storage', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...storageProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toBeInstanceOf(HTMLInputElement);
      });
      expect(screen.getByLabelText('End date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Total days in SIT')).toBeInTheDocument();
    });

    it('populates edit form with existing storage values', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...storageProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toHaveDisplayValue('15 Dec 2022');
      });
      expect(screen.getByLabelText('End date')).toHaveDisplayValue('25 Dec 2022');
    });

    it('correctly displays days in SIT', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...storageProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );
      await waitFor(() => {
        expect(screen.getByTestId('days-in-sit')).toHaveTextContent('11');
      });
    });

    it('correctly updates days in SIT', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...storageProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );
      const startDateInput = screen.getByLabelText('Start date');
      // clearing a date input throws an error about `children` being set to NaN
      // this replaces the '5' in '15 Dec 2022' with a '7' -> '17 Dec 2022'
      await userEvent.type(startDateInput, '7', {
        initialSelectionStart: 1,
        initialSelectionEnd: 2,
      });
      await waitFor(() => {
        expect(screen.getByTestId('days-in-sit')).toHaveTextContent('9');
      });
    });

    it('populates edit form with existing status and reason', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...rejectedProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeChecked();
      });
      expect(screen.getByText('Rejection reason')).toBeInTheDocument();
      expect(screen.getByText('484 characters')).toBeInTheDocument();
    });
  });

  describe('displays read only form', () => {
    it('renders disabled blank form on load with defaults', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...expenseRequiredProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
          readOnly
        />,
        {
          wrapper: MockProviders,
        },
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Receipt 1' })).toBeInTheDocument();
      });

      expect(screen.getByLabelText('Amount')).toBeDisabled();

      expect(screen.getByRole('heading', { level: 3, name: `Review Packing Materials #1` })).toBeInTheDocument();

      expect(screen.getByText(/Add a review for this/i)).toBeInTheDocument();

      expect(screen.getByLabelText('Accept')).toBeDisabled();
      expect(screen.getByLabelText('Exclude')).toBeDisabled();
      expect(screen.getByLabelText('Reject')).toBeDisabled();
    });

    it('populates disabled edit form with existing expense values', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...expenseRequiredProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
          readOnly
        />,
        {
          wrapper: MockProviders,
        },
      );

      await waitFor(() => {
        expect(screen.getByText('Expense Type')).toBeInTheDocument();
      });
      expect(screen.getByText('Packing materials')).toBeDisabled();
      expect(screen.getByDisplayValue('boxes, tape, bubble wrap'));
      expect(screen.getByLabelText('Amount')).toHaveDisplayValue('1,234.56');
      expect(screen.getByLabelText('Amount')).toBeDisabled();
    });

    it('populates disabled edit form with existing storage values', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...storageProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
          readOnly
        />,
        {
          wrapper: MockProviders,
        },
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toHaveDisplayValue('15 Dec 2022');
        expect(screen.getByLabelText('Start date')).toBeDisabled();
      });
      expect(screen.getByLabelText('End date')).toHaveDisplayValue('25 Dec 2022');
      expect(screen.getByLabelText('End date')).toBeDisabled();
    });

    it('populates disabled edit form with existing status and reason', async () => {
      render(
        <ReviewExpense
          {...defaultProps}
          {...rejectedProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
          readOnly
        />,
        {
          wrapper: MockProviders,
        },
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeChecked();
        expect(screen.getByLabelText('Reject')).toBeDisabled();
      });
      expect(screen.getByText('Rejection reason')).toBeInTheDocument();
      expect(screen.getByText('484 characters')).toBeInTheDocument();
    });
  });
});
