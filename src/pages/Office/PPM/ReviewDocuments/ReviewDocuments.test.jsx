import React from 'react';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { ReviewDocuments } from './ReviewDocuments';

import PPMDocumentsStatus from 'constants/ppms';
import { ppmShipmentStatuses } from 'constants/shipments';
import { usePPMShipmentDocsQueries } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import createUpload from 'utils/test/factories/upload';

Element.prototype.scrollTo = jest.fn();

beforeEach(() => {
  jest.clearAllMocks();
});

const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
  useLocation: jest.fn(),
}));

const mockPatchWeightTicket = jest.fn();
const mockPatchProGear = jest.fn();
const mockPatchExpense = jest.fn();

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  patchWeightTicket: (options) => mockPatchWeightTicket(options),
  patchProGearWeightTicket: (options) => mockPatchProGear(options),
  patchExpense: (options) => mockPatchExpense(options),
}));

// prevents react-fileviewer from throwing errors without mocking relevant DOM elements
jest.mock('components/DocumentViewer/Content/Content', () => {
  const MockContent = () => <div>Content</div>;
  return MockContent;
});

jest.mock('hooks/queries', () => ({
  usePPMShipmentDocsQueries: jest.fn(),
}));

const mtoShipment = createPPMShipmentWithFinalIncentive({
  ppmShipment: { status: ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL },
});

// The factory used above doesn't handle overrides for uploads correctly, so we need to do it manually.
const weightTicketEmptyDocumentUpload = createUpload({ fileName: 'emptyWeightTicket.pdf' });
const weightTicketFullDocumentUpload = createUpload(
  { fileName: 'fullWeightTicket.xls' },
  { contentType: 'application/vnd.ms-excel' },
);
const progearWeightTicketDocumentUpload = createUpload({ fileName: 'progearWeightTicket.pdf' });
const movingExpenseDocumentUpload = createUpload({ fileName: 'movingExpense.jpg' }, { contentType: 'image/jpeg' });

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

const requiredProps = {
  match: { params: { shipmentId: mtoShipment.id, moveCode: 'READY1' } },
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

    it('renders the Loading Placeholder when the query is still loading', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(loadingReturnValue);
      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

      const h2 = await screen.findByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });
    it('renders the Something Went Wrong component when the query errors', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(errorReturnValue);
      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

      const errorMessage = await screen.findByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });
  describe('with a single weight ticket loaded', () => {
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

    it('renders the DocumentViewer', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

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
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

      const newEmptyWeight = 14500;
      const emptyWeightInput = screen.getByRole('textbox', { name: 'Empty weight' });
      await userEvent.clear(emptyWeightInput);
      await userEvent.type(emptyWeightInput, newEmptyWeight.toString());

      const newFullWeight = 18500;
      const fullWeightInput = screen.getByRole('textbox', { name: 'Full weight' });
      await userEvent.clear(fullWeightInput);
      await userEvent.type(fullWeightInput, newFullWeight.toString());

      expect(screen.getByLabelText(/net weight/i)).toHaveTextContent('4,000 lbs');

      expect(await screen.findByLabelText('Accept')).toBeInTheDocument();

      const rejectOption = screen.getByLabelText('Reject');
      expect(rejectOption).toBeInTheDocument();
      await userEvent.click(rejectOption);

      expect(screen.getByLabelText('Reason')).toBeInTheDocument();

      const rejectionReason = 'Not legible';
      await userEvent.type(screen.getByLabelText('Reason'), rejectionReason);

      const continueButton = screen.getByRole('button', { name: 'Continue' });
      expect(continueButton).toBeInTheDocument();
      await userEvent.click(continueButton);

      expect(screen.queryByText('Reviewing this weight ticket is required')).not.toBeInTheDocument();

      const weightTicket = mtoShipmentWithOneWeightTicket.ppmShipment.weightTickets[0];
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

      expect(mockPatchWeightTicket).toHaveBeenCalledWith(expectedPayload);

      expect(await screen.findByRole('heading', { name: 'Send to customer?', level: 3 })).toBeInTheDocument();

      await userEvent.click(screen.getByRole('button', { name: 'Confirm' }));
      expect(mockPush).toHaveBeenCalled();
    });

    it('renders and handles the Close button', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

      const closeSidebarButton = await screen.findByRole('button', { name: /close sidebar/i });

      expect(closeSidebarButton).toBeInTheDocument();

      await userEvent.click(closeSidebarButton);
      expect(mockPush).toHaveBeenCalled();
    });

    it('shows an error if submissions fails', async () => {
      jest.spyOn(console, 'error').mockImplementation(() => {});

      mockPatchWeightTicket.mockRejectedValueOnce('fatal error');
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

      expect(await screen.findByRole('button', { name: 'Continue' })).toBeInTheDocument();

      await userEvent.click(screen.getByLabelText('Accept'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(screen.getByText('There was an error submitting the form. Please try again later.')).toBeInTheDocument();
    });

    it('handles navigation properly using the continue/back buttons', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueWithOneWeightTicket);
      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

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
    const usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets = {
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
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

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
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueMultipleWeightTickets);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

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

    it('only shows uploads for the document set being reviewed', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueAllDocs);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

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

    it('handles moving from weight tickets the summary page when there are multiple types of documents', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueAllDocs);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });

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
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueProGearOnly);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByText('Reviewing this pro-gear is required')).toBeInTheDocument();
    });

    it('shows an error when pro-gear is rejected and submitted without a written reason', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueProGearOnly);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });
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
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueExpensesOnly);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });
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
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueExpensesOnly);

      render(<ReviewDocuments {...requiredProps} />, { wrapper: MockProviders });
      await userEvent.click(screen.getByLabelText('Exclude'));
      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      expect(screen.getByText('Add a reason why this receipt is excluded')).toBeInTheDocument();
    });
  });
});
