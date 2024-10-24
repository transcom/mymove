import React from 'react';
import { screen, within, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { ReviewDocuments } from './ReviewDocuments';

import PPMDocumentsStatus from 'constants/ppms';
import { ppmShipmentStatuses } from 'constants/shipments';
import {
  usePPMShipmentDocsQueries,
  useReviewShipmentWeightsQuery,
  usePPMCloseoutQuery,
  useEditShipmentQueries,
  useGetPPMSITEstimatedCostQuery,
} from 'hooks/queries';
import { renderWithProviders } from 'testUtils';
import {
  createPPMShipmentWithFinalIncentive,
  createPPMShipmentWithExcessWeight,
} from 'utils/test/factories/ppmShipment';
import { createCompleteWeightTicket, createSecondCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import createUpload from 'utils/test/factories/upload';
import { servicesCounselingRoutes, tooRoutes } from 'constants/routes';

Element.prototype.scrollTo = jest.fn();

beforeEach(() => {
  jest.clearAllMocks();
});

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const mockPatchWeightTicket = jest.fn();
const mockPatchProGear = jest.fn();
const mockPatchExpense = jest.fn();
const mockPatchPPMDocumentsSetStatus = jest.fn();

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  patchWeightTicket: (options) => mockPatchWeightTicket(options),
  patchProGearWeightTicket: (options) => mockPatchProGear(options),
  patchExpense: (options) => mockPatchExpense(options),
  patchPPMDocumentsSetStatus: (options) => mockPatchPPMDocumentsSetStatus(options),
}));

// prevents react-fileviewer from throwing errors without mocking relevant DOM elements
jest.mock('components/DocumentViewer/Content/Content', () => {
  const MockContent = () => <div>Content</div>;
  return MockContent;
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
  estimatedCost: 5000,
};

const mockRoutingOptions = {
  path: servicesCounselingRoutes.BASE_REVIEW_SHIPMENT_WEIGHTS_PATH,
  params: { moveCode: 'READY1' },
};

const mockTooRountingOptions = {
  path: tooRoutes.BASE_MOVE_VIEW_PATH,
  params: { moveCode: 'READY1' },
};

describe('ReviewDocuments', () => {
  describe('check loading and error component states', () => {
    const loadingReturnValue = {
      isLoading: true,
      isError: false,
      isSuccess: false,
    };

    const errorReturnValue = {
      isLoading: false,
      isError: true,
      isSuccess: false,
    };

    it('renders the Loading Placeholder when the PPMCloseout query is still loading', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      usePPMCloseoutQuery.mockReturnValue(loadingReturnValue);
      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      const acceptOption = screen.getByTestId('approveRadio');
      expect(acceptOption).toBeInTheDocument();

      const rejectOption = screen.getByTestId('rejectRadio');
      expect(rejectOption).toBeInTheDocument();
      await userEvent.click(acceptOption);

      const continueButton = screen.getByTestId('reviewDocumentsContinueButton');
      expect(continueButton).toBeInTheDocument();
      await userEvent.click(continueButton);

      expect(screen.queryByText('Reviewing this weight ticket is required')).not.toBeInTheDocument();

      await waitFor(() => {
        expect(mockPatchWeightTicket).toHaveBeenCalled();
      });
      expect(await screen.findByRole('heading', { name: 'Send to customer?', level: 3 })).toBeInTheDocument();

      await userEvent.click(screen.getByTestId('shipmentInfo-showRequestDetailsButton'));
      await waitFor(() => {
        expect(screen.getByText('Show Details', { exact: false })).toBeInTheDocument();
      });

      const h2 = await screen.findByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Loading Placeholder when the PPMShipmentDocs query is still loading', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(loadingReturnValue);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);
      const h2 = await screen.findByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });
    it('renders the Something Went Wrong component when the PPMShipmentDocs query errors', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(errorReturnValue);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      const errorMessage = await screen.findByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });
  describe('with a single weight ticket loaded', () => {
    it('renders the DocumentViewer', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      const docMenuButton = await screen.findByRole('button', { name: /open menu/i });
      expect(docMenuButton).toBeInTheDocument();

      // We don't really have a better way to grab the DocumentViewerMenu to check its visibility because css isn't
      // loaded in the test environment. Instead, we'll grab it by its test id and check that it has the correct class.
      const docViewer = screen.getByTestId('DocViewerMenu');
      expect(docViewer).toHaveClass('collapsed');

      expect(within(docViewer).getByRole('heading', { level: 3, name: 'Documents' })).toBeInTheDocument();

      await userEvent.click(docMenuButton);

      expect(docViewer).not.toHaveClass('collapsed');

      const uploadList = within(docViewer).getByRole('list');
      expect(uploadList).toBeInTheDocument();

      expect(within(uploadList).getAllByRole('listitem').length).toBe(2);
      expect(within(uploadList).getByRole('button', { name: /emptyWeightTicket\.pdf.*/i })).toBeInTheDocument();
      expect(within(uploadList).getByRole('button', { name: /fullWeightTicket\.xls.*/i })).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 2, name: '1 of 1 Document Sets' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();

      expect(screen.getByRole('button', { name: /close sidebar/i })).toBeInTheDocument();
    });

    it('renders and handles the Continue button with the appropriate payload', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      const weightTicket = mtoShipmentWithOneWeightTicket.ppmShipment.weightTickets[0];

      const newEmptyWeight = 14500;
      const emptyWeightInput = screen.getByTestId('emptyWeight');
      await userEvent.clear(emptyWeightInput);
      await userEvent.type(emptyWeightInput, newEmptyWeight.toString());

      const newFullWeight = 18500;
      const fullWeightInput = screen.getByTestId('fullWeight');
      await userEvent.clear(fullWeightInput);
      await userEvent.type(fullWeightInput, newFullWeight.toString());

      const netWeightDisplay = screen.getByTestId('net-weight-display');
      expect(netWeightDisplay).toHaveTextContent('4,000 lbs');

      const acceptOption = screen.getByTestId('approveRadio');
      expect(acceptOption).toBeInTheDocument();

      const rejectOption = screen.getByTestId('rejectRadio');
      expect(rejectOption).toBeInTheDocument();
      await userEvent.click(rejectOption);

      const rejectionReason = 'Not legible';
      const reasonTextBox = screen.getByLabelText('Reason');
      expect(reasonTextBox).toBeInTheDocument();
      await userEvent.type(reasonTextBox, rejectionReason);

      const continueButton = screen.getByTestId('reviewDocumentsContinueButton');
      expect(continueButton).toBeInTheDocument();
      await userEvent.click(continueButton);

      expect(screen.queryByText('Reviewing this weight ticket is required')).not.toBeInTheDocument();

      const expectedPayload = {
        ppmShipmentId: mtoShipmentWithOneWeightTicket.ppmShipment.id,
        weightTicketId: weightTicket.id,
        eTag: weightTicket.eTag,
        payload: {
          ppmShipmentId: mtoShipmentWithOneWeightTicket.ppmShipment.id,
          vehicleDescription: weightTicket.vehicleDescription,
          emptyWeight: newEmptyWeight,
          missingEmptyWeightTicket: weightTicket.missingEmptyWeightTicket,
          fullWeight: newFullWeight,
          missingFullWeightTicket: weightTicket.missingFullWeightTicket,
          ownsTrailer: weightTicket.ownsTrailer,
          trailerMeetsCriteria: weightTicket.trailerMeetsCriteria,
          status: PPMDocumentsStatus.REJECTED,
          reason: rejectionReason,
        },
      };
      await waitFor(() => {
        expect(mockPatchWeightTicket).toHaveBeenCalledWith(expectedPayload);
      });

      expect(await screen.findByRole('heading', { name: 'Send to customer?', level: 3 })).toBeInTheDocument();

      await userEvent.click(screen.getByRole('button', { name: 'Confirm' }));
      const confirmPayload = {
        ppmShipmentId: mtoShipmentWithOneWeightTicket.ppmShipment.id,
        eTag: mtoShipmentWithOneWeightTicket.ppmShipment.eTag,
      };
      expect(mockPatchPPMDocumentsSetStatus).toHaveBeenCalledWith(confirmPayload);
      expect(mockNavigate).toHaveBeenCalled();
    });

    it('renders and handles the Close button', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      const closeSidebarButton = await screen.findByRole('button', { name: /close sidebar/i });

      expect(closeSidebarButton).toBeInTheDocument();

      await userEvent.click(closeSidebarButton);
      expect(mockNavigate).toHaveBeenCalled();
    });

    it('shows an error if submissions fails', async () => {
      jest.spyOn(console, 'error').mockImplementation(() => {});

      mockPatchWeightTicket.mockRejectedValueOnce('fatal error');
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      expect(await screen.findByRole('button', { name: 'Continue' })).toBeInTheDocument();

      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(screen.getByText('There was an error submitting the form. Please try again later.')).toBeInTheDocument();
    });

    it('handles navigation properly using the continue/back buttons', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      expect(await screen.findByRole('heading', { level: 2, name: '1 of 1 Document Sets' }));

      expect(await screen.findByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();

      // Need to accept the document before we can move forward without errors.
      await userEvent.click(screen.getByLabelText('Accept'));

      const continueButton = screen.getByRole('button', { name: 'Continue' });
      expect(continueButton).toBeEnabled();

      const backButton = screen.getByRole('button', { name: 'Back' });
      expect(backButton).not.toBeEnabled();

      await userEvent.click(continueButton);

      expect(await screen.findByRole('heading', { name: 'Send to customer?', level: 3 })).toBeInTheDocument();

      expect(backButton).toBeEnabled();
      await userEvent.click(backButton);

      expect(await screen.findByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();
    });

    it('handles navigation properly using the continue/back buttons', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      renderWithProviders(<ReviewDocuments />, mockTooRountingOptions);

      expect(await screen.findByRole('heading', { level: 2, name: '1 of 1 Document Sets' }));

      expect(await screen.findByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();

      // Need to accept the document before we can move forward without errors.
      await userEvent.click(screen.getByLabelText('Accept'));

      const continueButton = screen.getByRole('button', { name: 'Continue' });
      expect(continueButton).toBeEnabled();

      const backButton = screen.getByRole('button', { name: 'Back' });
      expect(backButton).not.toBeEnabled();

      await userEvent.click(continueButton);

      expect(await screen.findByRole('heading', { name: 'Send to customer?', level: 3 })).toBeInTheDocument();

      expect(backButton).toBeEnabled();
      await userEvent.click(backButton);

      expect(await screen.findByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();
    });
  });
  describe('with multiple document sets loaded', () => {
    let usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets = {
      ...usePPMShipmentDocsQueriesReturnValueAllDocs,
      documents: {
        ...usePPMShipmentDocsQueriesReturnValueAllDocs.documents,
        WeightTickets: [
          ...mtoShipment.ppmShipment.weightTickets,
          createCompleteWeightTicket({ serviceMemberId: mtoShipment.ppmShipment.serviceMemberId }),
        ],
      },
    };

    it('renders and handles the Accept button', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      expect(await screen.findByRole('heading', { level: 2, name: '1 of 4 Document Sets' }));
      expect(screen.getByLabelText('Accept')).toBeInTheDocument();
      expect(screen.getByLabelText('Reject')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByRole('heading', { level: 2, name: '2 of 4 Document Sets' }));
    });

    it('renders and handles the Back button', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      expect(screen.findByRole('heading', { level: 2, name: '1 of 4 Document Sets' }));
      expect(screen.getByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));

      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByRole('heading', { level: 2, name: '2 of 4 Document Sets' }));
      expect(screen.getByRole('heading', { level: 3, name: /trip 2/ })).toBeInTheDocument();

      await userEvent.click(screen.getByRole('button', { name: 'Back' }));
      expect(screen.getByRole('heading', { level: 2, name: '1 of 4 Document Sets' }));
      expect(screen.getByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();
    });

    it('renders weights correctly for multiple trips/weight tickets', async () => {
      usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets = {
        ...usePPMShipmentDocsQueriesReturnValueAllDocs,
        documents: {
          ...usePPMShipmentDocsQueriesReturnValueAllDocs.documents,
          WeightTickets: [
            createCompleteWeightTicket({ serviceMemberId: mtoShipment.ppmShipment.serviceMemberId }),
            createSecondCompleteWeightTicket({ serviceMemberId: mtoShipment.ppmShipment.serviceMemberId }),
          ],
        },
      };

      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      expect(screen.findByRole('heading', { level: 2, name: '1 of 4 Document Sets' }));
      expect(screen.getByRole('heading', { level: 3, name: /trip 1/ })).toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));

      // render the empty, full and allowable weights correctly for the first trip/weight ticket
      const weightTicketOne = usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets.documents.WeightTickets[0];
      expect(weightTicketOne.emptyWeight).toBe(14500);
      expect(screen.getByTestId('emptyWeight')).toHaveValue('14,500');
      expect(weightTicketOne.fullWeight).toBe(18500);
      expect(screen.getByTestId('fullWeight')).toHaveValue('18,500');

      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByRole('heading', { level: 2, name: '2 of 4 Document Sets' }));
      expect(screen.getByRole('heading', { level: 3, name: /trip 2/ })).toBeInTheDocument();

      // render the empty, full and allowable weights correctly for the second trip/weight ticket
      const weightTicketTwo = usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets.documents.WeightTickets[1];
      expect(weightTicketTwo.emptyWeight).toBe(10000);
      expect(screen.getByTestId('emptyWeight')).toHaveValue('10,000');
      expect(weightTicketTwo.fullWeight).toBe(12000);
      expect(screen.getByTestId('fullWeight')).toHaveValue('12,000');
    });

    it('only shows uploads for the document set being reviewed', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueAllDocs);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      const docMenuButton = await screen.findByRole('button', { name: /open menu/i });
      expect(docMenuButton).toBeInTheDocument();

      // We don't really have a great way to grab the list of uploads so we'll grab the parent element and go from there
      const docViewer = screen.getByTestId('DocViewerMenu');

      await userEvent.click(docMenuButton);

      expect(docViewer).not.toHaveClass('collapsed');

      const uploadList = within(docViewer).getByRole('list');
      expect(uploadList).toBeInTheDocument();

      expect(within(uploadList).getAllByRole('listitem').length).toBe(2);
      expect(within(uploadList).getByRole('button', { name: /emptyWeightTicket\.pdf.*/i })).toBeInTheDocument();
      expect(within(uploadList).getByRole('button', { name: /fullWeightTicket\.xls.*/i })).toBeInTheDocument();
      expect(within(uploadList).queryByRole('button', { name: /progearWeightTicket\.pdf.*/i })).not.toBeInTheDocument();
      expect(within(uploadList).queryByRole('button', { name: /movingExpense\.jpg.*/i })).not.toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 2, name: '1 of 3 Document Sets' })).toBeInTheDocument();
    });

    it('shows uploads for all documents on the summary page', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueAllDocs);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      expect(await screen.findByRole('heading', { name: 'Trip 1', level: 3 })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(await screen.findByRole('heading', { name: 'Pro-gear 1', level: 3 })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(await screen.findByRole('heading', { name: 'Receipt 1', level: 3 })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(await screen.findByRole('heading', { name: 'Send to customer?', level: 3 })).toBeInTheDocument();

      const docMenuButton = await screen.findByRole('button', { name: /open menu/i });
      expect(docMenuButton).toBeInTheDocument();

      // We don't really have a great way to grab the list of uploads so we'll grab the parent element and go from there
      const docViewer = screen.getByTestId('DocViewerMenu');

      await userEvent.click(docMenuButton);

      expect(docViewer).not.toHaveClass('collapsed');

      const uploadList = within(docViewer).getByRole('list');
      expect(uploadList).toBeInTheDocument();

      const uploadsButtons = within(uploadList).getAllByRole('listitem');
      expect(uploadsButtons.length).toBe(4);

      const allUploads = [
        weightTicketEmptyDocumentUpload,
        weightTicketFullDocumentUpload,
        progearWeightTicketDocumentUpload,
        movingExpenseDocumentUpload,
      ];

      // we expect uploads to be sorted in descending order by updatedAt
      allUploads.sort((a, b) => {
        if (a.updatedAt < b.updatedAt) {
          return 1;
        }

        if (a.updatedAt > b.updatedAt) {
          return -1;
        }

        return 0;
      });

      for (let i = 0; i < allUploads.length; i += 1) {
        // checking for text content because otherwise we'd have to form a regex to use the {name:} option of getByRole
        // and our linters don't like a regex that is formed using a variable because of the
        // security/detect-non-literal-regexp rule. Not super important to use it here, so we'll do this instead of
        // doing the IS3 process.
        expect(within(uploadsButtons[i]).getByRole('button')).toHaveTextContent(allUploads[i].filename);
      }
    });

    it('handles moving from weight tickets the summary page when there are multiple types of documents', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueAllDocs);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);

      expect(await screen.findByRole('heading', { name: 'Trip 1', level: 3 })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(await screen.findByRole('heading', { name: 'Pro-gear 1', level: 3 })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(await screen.findByRole('heading', { name: 'Receipt 1', level: 3 })).toBeInTheDocument();
      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(await screen.findByRole('heading', { name: 'Send to customer?', level: 3 })).toBeInTheDocument();
      expect(await screen.getByRole('button', { name: 'Back' })).toBeEnabled();

      expect(screen.getByRole('heading', { level: 2, name: /All Document Sets/ })).toBeInTheDocument();
    });

    const usePPMShipmentDocsQueriesReturnValueProGearOnly = {
      ...usePPMShipmentDocsQueriesReturnValueAllDocs,
      mtoShipment,
      documents: {
        MovingExpenses: [],
        ProGearWeightTickets: [...mtoShipment.ppmShipment.proGearWeightTickets],
        WeightTickets: [],
      },
    };

    it('shows an error when submitting without a status selected', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueProGearOnly);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByText('Reviewing this pro-gear is required')).toBeInTheDocument();
    });

    it('shows an error when pro-gear is rejected and submitted without a written reason', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueProGearOnly);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);
      const rejectionButton = screen.getByTestId('rejectRadio');
      expect(rejectionButton).toBeInTheDocument();
      await userEvent.click(rejectionButton);
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByText('Add a reason why this pro-gear is rejected'));
    });

    it('shows an error when a rejected expense is submitted with no reason', async () => {
      const usePPMShipmentDocsQueriesReturnValueExpensesOnly = {
        ...usePPMShipmentDocsQueriesReturnValueAllDocs,
        mtoShipment,
        documents: {
          MovingExpenses: [...mtoShipment.ppmShipment.movingExpenses],
          ProGearWeightTickets: [],
          WeightTickets: [],
        },
      };
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueExpensesOnly);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);
      await userEvent.click(screen.getByLabelText('Reject'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(screen.getByText('Add a reason why this receipt is rejected')).toBeInTheDocument();
    });

    it('shows an error when an excluded expense is submitted with no reason', async () => {
      const usePPMShipmentDocsQueriesReturnValueExpensesOnly = {
        ...usePPMShipmentDocsQueriesReturnValueAllDocs,
        mtoShipment,
        documents: {
          MovingExpenses: [...mtoShipment.ppmShipment.movingExpenses],
          ProGearWeightTickets: [],
          WeightTickets: [],
        },
      };
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueExpensesOnly);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);
      useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);
      await userEvent.click(screen.getByLabelText('Exclude'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByText('Add a reason why this receipt is excluded')).toBeInTheDocument();
    });
  });
  describe('check over weight alerts', () => {
    it('does not display an alert when move is not over weight', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueAllDocs);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueAll);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);
      const alert = screen.queryByText('This move has excess weight. Edit the PPM net weight to resolve.');
      expect(alert).toBeNull();
    });

    it('displays an alert when move is over weight', async () => {
      const excessWeightPPMShipment = createPPMShipmentWithExcessWeight({
        ppmShipment: { status: ppmShipmentStatuses.NEEDS_CLOSEOUT },
      });
      const useReviewShipmentWeightsQueryReturnValueExcessWeight = {
        ...useReviewShipmentWeightsQueryReturnValueAll,
        mtoShipments: [excessWeightPPMShipment],
      };
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueAllDocs);
      usePPMCloseoutQuery.mockReturnValue(usePPMCloseoutQueryReturnValue);
      useReviewShipmentWeightsQuery.mockReturnValue(useReviewShipmentWeightsQueryReturnValueExcessWeight);

      renderWithProviders(<ReviewDocuments />, mockRoutingOptions);
      const alert = screen.getByText('This move has excess weight. Edit the PPM net weight to resolve.');
      expect(alert).toBeInTheDocument();
    });
  });
});
