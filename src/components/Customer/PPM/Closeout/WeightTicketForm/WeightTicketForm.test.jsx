import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import WeightTicketForm from 'components/Customer/PPM/Closeout/WeightTicketForm/WeightTicketForm';
import { DocumentAndImageUploadInstructions } from 'content/uploads';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  weightTicket: {
    id: '58350bae-8e87-4e83-bd75-74027fb4853f',
    shipmentId: '8be77cb9-e8af-4ff0-b0a2-ade17cf6653c',
    emptyWeightDocumentId: '27d70a0d-7f20-42af-ab79-f74350412823',
    fullWeightDocumentId: '1ec00b40-447d-4c22-ac73-708b98b8bc20',
    trailerOwnershipDocumentId: '5bf3ed20-08dd-4d8e-92ad-7603bb6377a5',
  },
  tripNumber: 2,
  onCreateUpload: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
  onBack: jest.fn(),
  onSubmit: jest.fn(),
};

const weightTicketRequiredProps = {
  weightTicket: {
    id: '58350bae-8e87-4e83-bd75-74027fb4853f',
    shipmentId: '8be77cb9-e8af-4ff0-b0a2-ade17cf6653c',
    vehicleDescription: 'DMC Delorean',
    emptyWeight: 3999,
    emptyWeightDocumentId: '27d70a0d-7f20-42af-ab79-f74350412823',
    emptyWeightTickets: [
      {
        id: '299e2fb4-432d-4261-bbed-d8280c6090af',
        created_at: '2022-06-22T23:25:50.490Z',
        bytes: 819200,
        url: 'a/fake/path',
        filename: 'empty_weight.jpg',
        content_type: 'image/jpg',
      },
    ],
    fullWeight: 7111,
    fullWeightDocumentId: '1ec00b40-447d-4c22-ac73-708b98b8bc20',
    fullWeightTickets: [
      {
        id: 'f70af8a1-38e9-4ae2-a837-3c0c61069a0d',
        created_at: '2022-06-23T23:25:50.490Z',
        bytes: 409600,
        url: 'a/fake/path',
        filename: 'full_weight.pdf',
        content_type: 'application/pdf',
      },
    ],
    hasOwnTrailer: false,
    trailerOwnershipDocumentId: '5bf3ed20-08dd-4d8e-92ad-7603bb6377a5',
  },
};

const weightTicketUploadsOnlyProps = {
  weightTicket: {
    id: '58350bae-8e87-4e83-bd75-74027fb4853f',
    shipmentId: '8be77cb9-e8af-4ff0-b0a2-ade17cf6653c',
    emptyWeightDocumentId: '27d70a0d-7f20-42af-ab79-f74350412823',
    emptyWeightTickets: [
      {
        id: '299e2fb4-432d-4261-bbed-d8280c6090af',
        created_at: '2022-06-22T23:25:50.490Z',
        bytes: 819200,
        url: 'a/fake/path',
        filename: 'empty_weight.jpg',
        content_type: 'image/jpg',
      },
    ],
    fullWeightDocumentId: '1ec00b40-447d-4c22-ac73-708b98b8bc20',
    fullWeightTickets: [
      {
        id: 'f70af8a1-38e9-4ae2-a837-3c0c61069a0d',
        created_at: '2022-06-23T23:25:50.490Z',
        bytes: 409600,
        url: 'a/fake/path',
        filename: 'full_weight.pdf',
        content_type: 'application/pdf',
      },
    ],
  },
};

const constructedWeightTrailerProps = {
  weightTicket: {
    id: '58350bae-8e87-4e83-bd75-74027fb4853f',
    shipmentId: '8be77cb9-e8af-4ff0-b0a2-ade17cf6653c',
    vehicleDescription: 'DMC Delorean',
    emptyWeight: 3999,
    missingEmptyWeightTicket: true,
    emptyWeightDocumentId: '27d70a0d-7f20-42af-ab79-f74350412823',
    emptyWeightTickets: [
      {
        id: '299e2fb4-432d-4261-bbed-d8280c6090af',
        created_at: '2022-06-22T23:25:50.490Z',
        bytes: 819200,
        url: 'a/fake/path',
        filename: 'weight estimator.xls',
        content_type: 'application/vnd.ms-excel',
      },
    ],
    fullWeight: 7111,
    missingFullWeightTicket: true,
    fullWeightDocumentId: '1ec00b40-447d-4c22-ac73-708b98b8bc20',
    fullWeightTickets: [
      {
        id: 'f70af8a1-38e9-4ae2-a837-3c0c61069a0d',
        created_at: '2022-06-23T23:25:50.490Z',
        bytes: 409600,
        url: 'a/fake/path',
        filename: 'weight estimator.xlsx',
        content_type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      },
    ],
    hasOwnTrailer: true,
    hasClaimedTrailer: true,
    trailerOwnershipDocumentId: '5bf3ed20-08dd-4d8e-92ad-7603bb6377a5',
    trailerOwnershipDocs: [
      {
        id: 'fd4e80f8-d025-44b2-8c33-15240fac51ab',
        created_at: '2022-06-24T23:25:50.490Z',
        bytes: 204800,
        url: 'a/fake/path',
        filename: 'trailer_title.pdf',
        content_type: 'application/pdf',
      },
    ],
  },
};

describe('WeightTicketForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<WeightTicketForm {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Trip 2' })).toBeInTheDocument();
      });

      expect(screen.getByRole('heading', { level: 3, name: 'Vehicle' })).toBeInTheDocument();
      expect(screen.getByLabelText('Vehicle description')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Car make and model, type of truck or van, etc.')).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Empty weight' })).toBeInTheDocument();
      expect(screen.getByLabelText('Empty weight')).toBeInstanceOf(HTMLInputElement);
      const missingWeightInput = screen.getAllByLabelText("I don't have this weight ticket");
      expect(missingWeightInput[0]).toBeInstanceOf(HTMLInputElement);
      expect(missingWeightInput[0]).not.toBeChecked();
      // getByLabelText will fail because the file upload input adds an aria-labeledby that points to the container text
      expect(screen.getByText('Upload empty weight ticket')).toBeInstanceOf(HTMLLabelElement);
      const uploadFileTypeHints = screen.getAllByText(DocumentAndImageUploadInstructions);
      expect(uploadFileTypeHints[0]).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Full weight' })).toBeInTheDocument();
      expect(screen.getByLabelText('Full weight')).toBeInstanceOf(HTMLInputElement);
      expect(missingWeightInput[1]).toBeInstanceOf(HTMLInputElement);
      expect(missingWeightInput[1]).not.toBeChecked();
      // getByLabelText will fail because the file upload input adds an aria-labeledby that points to the container text
      expect(screen.getByText('Upload full weight ticket')).toBeInstanceOf(HTMLLabelElement);
      expect(uploadFileTypeHints[1]).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Trip weight:' })).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Trailer' })).toBeInTheDocument();
      expect(screen.getByText('On this trip, were you using a trailer that you own?')).toBeInstanceOf(
        HTMLLegendElement,
      );
      expect(screen.getByLabelText('No')).toBeChecked();

      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('populates edit form with existing weight ticket values', async () => {
      render(<WeightTicketForm {...defaultProps} {...weightTicketRequiredProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('Empty weight')).toHaveDisplayValue('3,999');
      });

      expect(screen.getByText('empty_weight.jpg')).toBeInTheDocument();
      const deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons[0]).toBeInTheDocument();
      expect(screen.getByText('800KB')).toBeInTheDocument();
      expect(screen.getByText('Uploaded 22 Jun 2022 11:25 PM')).toBeInTheDocument();

      expect(screen.getByLabelText('Full weight')).toHaveDisplayValue('7,111');
      expect(screen.getByText('full_weight.pdf')).toBeInTheDocument();
      expect(deleteButtons[1]).toBeInTheDocument();
      expect(screen.getByText('400KB')).toBeInTheDocument();
      expect(screen.getByText('Uploaded 23 Jun 2022 11:25 PM')).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Trip weight: 3,112 lbs' })).toBeInTheDocument();

      expect(screen.getByLabelText('No')).toBeChecked();

      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('populates edit form with constructed weight and trailer values', async () => {
      render(<WeightTicketForm {...defaultProps} {...constructedWeightTrailerProps} />);

      let missingWeightInput;
      await waitFor(() => {
        missingWeightInput = screen.getAllByLabelText("I don't have this weight ticket");
        expect(missingWeightInput[0]).toBeChecked();
      });

      const downloadConstructedWeight = screen.getAllByRole('link');
      expect(downloadConstructedWeight[0]).toHaveTextContent('Go to download page');
      expect(downloadConstructedWeight[0]).toHaveAttribute(
        'href',
        'https://www.ustranscom.mil/dp3/weightestimator.cfm',
      );

      expect(screen.getByText('weight estimator.xls')).toBeInTheDocument();
      const deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons[0]).toBeInTheDocument();
      expect(screen.getByText('800KB')).toBeInTheDocument();
      expect(screen.getByText('Uploaded 22 Jun 2022 11:25 PM')).toBeInTheDocument();

      expect(missingWeightInput[1]).toBeChecked();

      expect(downloadConstructedWeight[1]).toHaveTextContent('Go to download page');
      expect(downloadConstructedWeight[1]).toHaveAttribute(
        'href',
        'https://www.ustranscom.mil/dp3/weightestimator.cfm',
      );

      expect(screen.getByText('weight estimator.xlsx')).toBeInTheDocument();
      expect(deleteButtons[1]).toBeInTheDocument();
      expect(screen.getByText('400KB')).toBeInTheDocument();
      expect(screen.getByText('Uploaded 23 Jun 2022 11:25 PM')).toBeInTheDocument();

      const trailerCheckboxes = screen.getAllByLabelText('Yes');
      expect(trailerCheckboxes[0]).toBeChecked();
      expect(trailerCheckboxes[1]).toBeChecked();

      expect(screen.getByText('Upload proof of ownership')).toBeInstanceOf(HTMLLabelElement);

      expect(screen.getByText('trailer_title.pdf')).toBeInTheDocument();
      expect(deleteButtons[1]).toBeInTheDocument();
      expect(screen.getByText('200KB')).toBeInTheDocument();
      expect(screen.getByText('Uploaded 24 Jun 2022 11:25 PM')).toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });
  });
  describe('validates the form', () => {
    it('marks required fields of empty form', async () => {
      render(<WeightTicketForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      let invalidAlerts;
      await waitFor(() => {
        invalidAlerts = screen.getAllByRole('alert');
      });

      expect(invalidAlerts).toHaveLength(5);

      expect(invalidAlerts[0].nextSibling).toHaveAttribute('name', 'vehicleDescription');
      expect(within(invalidAlerts[1].nextSibling).getByLabelText('Empty weight')).toBeInTheDocument();
      // Had no luck trying to get the label of the file input with the aria-describedby
      expect(within(invalidAlerts[2].previousSibling).getByText('Upload empty weight ticket')).toBeInTheDocument();
      expect(within(invalidAlerts[3].nextSibling).getByLabelText('Full weight')).toBeInTheDocument();
      expect(within(invalidAlerts[4].previousSibling).getByText('Upload full weight ticket')).toBeInTheDocument();
    });

    it('triggers error if the full weight is less than the empty weight', async () => {
      render(<WeightTicketForm {...defaultProps} />);

      await userEvent.type(screen.getByLabelText('Empty weight'), '4999');
      await userEvent.type(screen.getByLabelText('Full weight'), '4999');

      await waitFor(() => {
        expect(screen.getByText('The full weight must be greater than the empty weight')).toBeInTheDocument();
      });
    });
  });
  describe('attaches button handler callbacks', () => {
    it('calls the onSubmit callback with required fields', async () => {
      render(<WeightTicketForm {...defaultProps} {...weightTicketUploadsOnlyProps} />);

      await userEvent.type(screen.getByLabelText('Vehicle description'), 'DMC Delorean');
      await userEvent.type(screen.getByLabelText('Empty weight'), '4999');
      await userEvent.type(screen.getByLabelText('Full weight'), '6999');

      /* testing-library's upload helper doesn't seem to be detected with our use of filepond

      // we can't query for the file inputs because they aren't accessible roles and the hidden aria-labelledby
      // isn't found by testing-library
      const uploadFileHints = screen.getAllByText(DocumentAndImageUploadInstructions);
      const uploadEmptyWeight = uploadFileHints[0].nextSibling.firstChild;

      const emptyWeightFile = new File(['empty weight'], 'empty weight.png', { type: 'image/png' });
      await userEvent.upload(uploadEmptyWeight, emptyWeightFile);

      expect(uploadEmptyWeight.files[0]).toBe(emptyWeightFile);
      expect(uploadEmptyWeight.files.item(0)).toBe(emptyWeightFile);
      expect(uploadEmptyWeight.files).toHaveLength(1);

      const uploadFullWeight = uploadFileHints[1].nextSibling.firstChild;

      const fullWeightFile = new File(['full weight'], 'full weight.png', { type: 'image/png' });
      await userEvent.upload(uploadFullWeight, fullWeightFile);
      */

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(
          {
            vehicleDescription: 'DMC Delorean',
            emptyWeight: '4999',
            missingEmptyWeightTicket: false,
            emptyWeightTickets: [
              {
                id: '299e2fb4-432d-4261-bbed-d8280c6090af',
                created_at: '2022-06-22T23:25:50.490Z',
                bytes: 819200,
                url: 'a/fake/path',
                filename: 'empty_weight.jpg',
                content_type: 'image/jpg',
              },
            ],
            fullWeight: '6999',
            missingFullWeightTicket: false,
            fullWeightTickets: [
              {
                id: 'f70af8a1-38e9-4ae2-a837-3c0c61069a0d',
                created_at: '2022-06-23T23:25:50.490Z',
                bytes: 409600,
                url: 'a/fake/path',
                filename: 'full_weight.pdf',
                content_type: 'application/pdf',
              },
            ],
            hasOwnTrailer: 'false',
            hasClaimedTrailer: 'false',
            trailerOwnershipDocs: [],
          },
          expect.anything(),
        );
      });
    });
    it('calls the onSubmit callback with constructed weight and trailer ownership values', async () => {
      render(<WeightTicketForm {...defaultProps} {...constructedWeightTrailerProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(
          {
            vehicleDescription: 'DMC Delorean',
            emptyWeight: '3999',
            missingEmptyWeightTicket: true,
            emptyWeightTickets: [
              {
                id: '299e2fb4-432d-4261-bbed-d8280c6090af',
                created_at: '2022-06-22T23:25:50.490Z',
                bytes: 819200,
                url: 'a/fake/path',
                filename: 'weight estimator.xls',
                content_type: 'application/vnd.ms-excel',
              },
            ],
            fullWeight: '7111',
            missingFullWeightTicket: true,
            fullWeightTickets: [
              {
                id: 'f70af8a1-38e9-4ae2-a837-3c0c61069a0d',
                created_at: '2022-06-23T23:25:50.490Z',
                bytes: 409600,
                url: 'a/fake/path',
                filename: 'weight estimator.xlsx',
                content_type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
              },
            ],
            hasOwnTrailer: 'true',
            hasClaimedTrailer: 'true',
            trailerOwnershipDocs: [
              {
                id: 'fd4e80f8-d025-44b2-8c33-15240fac51ab',
                created_at: '2022-06-24T23:25:50.490Z',
                bytes: 204800,
                url: 'a/fake/path',
                filename: 'trailer_title.pdf',
                content_type: 'application/pdf',
              },
            ],
          },
          expect.anything(),
        );
      });
    });
    it('calls the onBack prop when the Finish Later button is clicked', async () => {
      render(<WeightTicketForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });
});
