import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewExpense from './ReviewExpense';

import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import { ppmShipmentStatuses } from 'constants/shipments';
import createUpload from 'utils/test/factories/upload';
import ppmDocumentStatus from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';
import { MockProviders } from 'testUtils';
import {
  useGetPPMSITEstimatedCostQuery,
  useReviewShipmentWeightsQuery,
  usePPMCloseoutQuery,
  useEditShipmentQueries,
  usePPMShipmentDocsQueries,
} from 'hooks/queries';
import { formatWeight } from 'utils/formatters';
import { getPPMTypeLabel, PPM_TYPES } from 'shared/constants';
import { hasProGearSPR, hasSpouseProGearSPR } from 'utils/ppmCloseout';

beforeEach(() => {
  jest.clearAllMocks();
});

jest.mock('hooks/queries', () => ({
  usePPMShipmentDocsQueries: jest.fn(),
  usePPMCloseoutQuery: jest.fn(),
  useReviewShipmentWeightsQuery: jest.fn(),
  useEditShipmentQueries: jest.fn(),
  useGetPPMSITEstimatedCostQuery: jest.fn(),
}));

const useEditShipmentQueriesReturnValue = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: 'NEEDS SERVICE COUNSELING',
  },
  order: {
    id: '1',
    originDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
      },
    },
    destinationDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: 'PERMANENT_CHANGE_OF_STATION',
    order_type_detail: 'HHG_PERMITTED',
    tac: '9999',
  },
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'shipment123',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      requestedDeliveryDate: '2018-04-15',
      scheduledDeliveryDate: '2014-04-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const mtoShipment = createPPMShipmentWithFinalIncentive({
  ppmShipment: { status: ppmShipmentStatuses.NEEDS_CLOSEOUT },
});

const weightTicketEmptyDocumentCreatedDate = new Date();
// The factory used above doesn't handle overrides for uploads correctly, so we need to do it manually.
const weightTicketEmptyDocumentUpload = createUpload({
  fileName: 'emptyWeightTicket.pdf',
  createdAtDate: weightTicketEmptyDocumentCreatedDate,
});

const weightTicketFullDocumentCreatedDate = new Date(weightTicketEmptyDocumentCreatedDate);
weightTicketFullDocumentCreatedDate.setDate(weightTicketFullDocumentCreatedDate.getDate() + 1);
const weightTicketFullDocumentUpload = createUpload(
  { fileName: 'fullWeightTicket.xls', createdAtDate: weightTicketFullDocumentCreatedDate },
  { contentType: 'application/vnd.ms-excel' },
);

const progearWeightTicketDocumentCreatedDate = new Date(weightTicketFullDocumentCreatedDate);
progearWeightTicketDocumentCreatedDate.setDate(progearWeightTicketDocumentCreatedDate.getDate() + 1);
const progearWeightTicketDocumentUpload = createUpload({
  fileName: 'progearWeightTicket.pdf',
  createdAtDate: progearWeightTicketDocumentCreatedDate,
});

const movingExpenseDocumentCreatedDate = new Date(progearWeightTicketDocumentCreatedDate);
movingExpenseDocumentCreatedDate.setDate(movingExpenseDocumentCreatedDate.getDate() + 1);
const movingExpenseDocumentUpload = createUpload(
  { fileName: 'movingExpense.jpg', createdAtDate: movingExpenseDocumentCreatedDate },
  { contentType: 'image/jpeg' },
);

mtoShipment.ppmShipment.weightTickets[0].emptyDocument.uploads = [weightTicketEmptyDocumentUpload];
mtoShipment.ppmShipment.weightTickets[0].fullDocument.uploads = [weightTicketFullDocumentUpload];
mtoShipment.ppmShipment.proGearWeightTickets[0].document.uploads = [progearWeightTicketDocumentUpload];
mtoShipment.ppmShipment.movingExpenses[0].document.uploads = [movingExpenseDocumentUpload];

const usePPMShipmentDocsQueriesReturnValueAllDocs = {
  mtoShipment,
  documents: {
    MovingExpenses: [...mtoShipment.ppmShipment.movingExpenses],
    ProGearWeightTickets: [...mtoShipment.ppmShipment.proGearWeightTickets],
    WeightTickets: [...mtoShipment.ppmShipment.weightTickets],
  },
  isError: false,
  isLoading: false,
  isSuccess: true,
};

const mtoShipmentWithOneWeightTicket = {
  ...mtoShipment,
  ppmShipment: {
    ...mtoShipment.ppmShipment,
    proGearWeightTickets: [],
    movingExpenses: [],
  },
};

const usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket = {
  ...usePPMShipmentDocsQueriesReturnValueAllDocs,
  mtoShipment: mtoShipmentWithOneWeightTicket,
  documents: {
    MovingExpenses: [],
    ProGearWeightTickets: [],
    WeightTickets: [...mtoShipment.ppmShipment.weightTickets],
  },
};

/**
 * @constant {Object} useReviewShipmentWeightsQueryReturnValueAll
 * @description The mocked return values from the useReviewShipmentWeightsQuery
 * that is being used by the EditPPMNetWeight component inside of the
 * ReviewWeightTicket component
 * */
const useReviewShipmentWeightsQueryReturnValueAll = {
  orders: {
    orderID: {
      entitlement: {
        totalWeight: 1000,
      },
    },
  },
  mtoShipments: [],
};

const usePPMCloseoutQueryReturnValue = {
  ppmCloseout: {
    SITReimbursement: 0,
    actualMoveDate: '2020-03-16',
    actualWeight: 4002,
    aoa: 340000,
    ddp: 33297,
    dop: 15048,
    estimatedWeight: 4000,
    gcc: 17102245,
    grossIncentive: 4855170,
    haulFSC: 403,
    haulPrice: 4529083,
    id: '1a719536-02ba-44cd-b97d-5a0548237dc5',
    miles: 415,
    packPrice: 253447,
    plannedMoveDate: '2020-03-15',
    proGearWeightCustomer: 500,
    proGearWeightSpouse: 0,
    remainingIncentive: 4515170,
    unpackPrice: 23892,
  },
  isError: false,
  isLoading: false,
  isSuccess: true,
};

const useGetPPMSITEstimatedCostQueryReturnValue = {
  estimatedCost: {
    sitCost: 5000,
  },
  isError: false,
  isLoading: false,
  isSuccess: true,
};

const useGetPPMSITEstimatedCostQueryLoading = {
  ...useGetPPMSITEstimatedCostQueryReturnValue,
  isError: false,
  isLoading: true,
  isSuccess: false,
};

const useGetPPMSITEstimatedCostQueryError = {
  ...useGetPPMSITEstimatedCostQueryReturnValue,
  isError: true,
  isLoading: false,
  isSuccess: false,
};

const defaultProps = {
  ppmShipmentInfo: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    expectedDepartureDate: '2022-12-02',
    actualMoveDate: '2022-12-06',
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
    estimatedCost: 3000,
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
    weightStored: 2000,
    sitLocation: 'ORIGIN',
    amount: 456,
  },
};

const storagePropsLowerGovEstimate = {
  expense: {
    ...expenseRequiredProps.expense,
    movingExpenseType: expenseTypes.STORAGE,
    sitStartDate: '2022-12-15',
    sitEndDate: '2022-12-25',
    weightStored: 2000,
    sitLocation: 'ORIGIN',
    amount: 6000,
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
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
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

      expect(screen.getAllByText('Expense Type')).toHaveLength(2);
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByLabelText('Amount Requested')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: `Review Packing Materials #1` })).toBeInTheDocument();

      expect(screen.getByText(/Add a review for this/i)).toBeInTheDocument();

      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Exclude')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing expense values', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
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
        expect(screen.getAllByText('Expense Type')).toHaveLength(2);
      });
      expect(screen.getByText('Packing materials')).toBeInTheDocument();
      expect(screen.getByDisplayValue('boxes, tape, bubble wrap'));
      expect(screen.getByLabelText('Amount Requested')).toHaveDisplayValue('1,234.56');
    });

    it('shows SIT fields when expense type is Storage', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      await useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);

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

      expect(screen.getByLabelText('Origin')).toBeChecked();
      expect(screen.getByLabelText('Destination')).not.toBeChecked();
      expect(screen.getByText('Actual SIT Reimbursement')).toBeInTheDocument();
      // Actual sit reimbursement is the lower of the requested values between gov estimate and user request
      await waitFor(() => {
        expect(screen.getByTestId('actual-sit-reimbursement')).toHaveTextContent('$4.56');
      });
      expect(screen.getByLabelText('Weight Stored')).toHaveDisplayValue('2,000');

      await waitFor(() => {
        expect(screen.getByTestId('costAmountSuccess')).toBeInTheDocument();
      });

      expect(screen.getByTestId('costAmountSuccess')).toHaveTextContent('$50.00');
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toBeInstanceOf(HTMLInputElement);
      });
      expect(screen.getByLabelText('End date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Total days in SIT')).toBeInTheDocument();
      await waitFor(() => {
        expect(screen.getByTestId('days-in-sit')).toHaveTextContent('11');
      });
    });

    it('actual sit reimbursement uses gov estimate when lower than user request', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      await useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);

      render(
        <ReviewExpense
          {...defaultProps}
          {...storagePropsLowerGovEstimate}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        {
          wrapper: MockProviders,
        },
      );
      expect(screen.getByText('Actual SIT Reimbursement')).toBeInTheDocument();
      await waitFor(() => {
        expect(screen.getByTestId('actual-sit-reimbursement')).toHaveTextContent('$50.00');
      });
    });

    it('renders the $0 cost when the query is still loading', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryLoading);
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

      const costAmount = screen.getByTestId('costAmount');
      expect(costAmount).toHaveTextContent('$0.00');
    });

    it('renders $0 cost when the query errors', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryError);
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

      const errorMessage = screen.getByTestId('costAmount');
      expect(errorMessage).toHaveTextContent('$0.00');
    });

    it('correctly updates days in SIT', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);
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
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
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

      expect(screen.getByLabelText('Amount Requested')).toBeDisabled();

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
        expect(screen.getAllByText('Expense Type')).toHaveLength(2);
      });
      expect(screen.getByText('Packing materials')).toBeDisabled();
      expect(screen.getByDisplayValue('boxes, tape, bubble wrap'));
      expect(screen.getByLabelText('Amount Requested')).toHaveDisplayValue('1,234.56');
      expect(screen.getByLabelText('Amount Requested')).toBeDisabled();
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

  describe('ReviewExpense - small package expense changes', () => {
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
    usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
    usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
    useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
    it('renders small package details when ppmType is SMALL_PACKAGE', async () => {
      const smallPackagePPMInfo = {
        ...defaultProps.ppmShipmentInfo,
        ppmType: PPM_TYPES.SMALL_PACKAGE,
        allowableWeight: 4000,
        movingExpenses: [
          { weightShipped: 1000, isProGear: false },
          { weightShipped: 500, isProGear: true, proGearBelongsToSelf: false },
        ],
      };

      render(
        <ReviewExpense
          {...defaultProps}
          ppmShipmentInfo={smallPackagePPMInfo}
          {...expenseRequiredProps}
          {...documentSetsProps}
          documentSetIndex={documentSetIndex}
        />,
        { wrapper: MockProviders },
      );

      await waitFor(() => {
        expect(screen.getByTestId('smallPackageTag')).toBeInTheDocument();
      });

      expect(screen.getByTestId('smallPackageTag')).toHaveTextContent(getPPMTypeLabel(PPM_TYPES.SMALL_PACKAGE));
      expect(screen.getByText('Allowable Weight')).toBeInTheDocument();
      expect(screen.getByText(formatWeight(4000))).toBeInTheDocument();

      // add up the total weight shipped from moving expenses
      const totalWeightShipped =
        smallPackagePPMInfo.movingExpenses[0].weightShipped + smallPackagePPMInfo.movingExpenses[1].weightShipped; // 1500 lbs
      expect(screen.getByText('Total Weight Shipped')).toBeInTheDocument();
      expect(screen.getByText(formatWeight(totalWeightShipped))).toBeInTheDocument();

      expect(screen.getByText('Pro-gear')).toBeInTheDocument();
      expect(screen.getByText('Spouse Pro-gear')).toBeInTheDocument();

      // there will be TWO of these (returns 'Yes')
      // hasProGearSPR should return "Yes" because one expense has isProGear true
      // hasSpouseProGearSPR should return "Yes" since the expense with pro-gear has proGearBelongsToSelf false
      expect(screen.getAllByText(hasProGearSPR(smallPackagePPMInfo.movingExpenses))).toHaveLength(2);
      expect(screen.getAllByText(hasSpouseProGearSPR(smallPackagePPMInfo.movingExpenses))).toHaveLength(2);
    });
  });
});
