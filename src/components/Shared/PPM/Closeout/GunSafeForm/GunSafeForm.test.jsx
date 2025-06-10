import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import { MockProviders } from 'testUtils';
import GunSafeForm from 'components/Shared/PPM/Closeout/GunSafeForm/GunSafeForm';
import { APP_NAME } from 'constants/apps';

const defaultProps = {
  gunSafe: {
    id: '58350bae-8e87-4e83-bd75-74027fb4853f',
    shipmentId: '8be77cb9-e8af-4ff0-b0a2-ade17cf6653c',
    weight: 145,
  },
  entitlements: {
    gunSafeWeight: 450,
  },
  onCreateUpload: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
  onBack: jest.fn(),
  onSubmit: jest.fn(),
};

const mockGunSafeEntitlements = {
  gunSafe: 500,
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectGunSafeEntitlements: jest.fn().mockImplementation(() => mockGunSafeEntitlements),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const gunSafeProps = {
  gunSafe: {
    document: {},
    weight: 1,
    description: 'Description',
    hasWeightTickets: '',
  },
};

const gunSafeNoWeightProps = {
  gunSafe: {
    ...gunSafeProps.gunSafe,
    weight: 0,
  },
};

const gunSafeWithDocumentProps = {
  gunSafe: {
    ...gunSafeProps.gunSafe,
    document: {
      uploads: [
        {
          id: '299e2fb4-432d-4261-bbed-d8280c6090af',
          createdAt: '2022-06-22T23:25:50.490Z',
          bytes: 819200,
          url: 'a/fake/path',
          filename: 'empty_weight.pdf',
          contentType: 'image/pdf',
        },
      ],
    },
  },
};

describe('GunSafeForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults - Customer page', () => {
      render(<GunSafeForm {...defaultProps} appName={APP_NAME.MYMOVE} />, { wrapper: MockProviders });

      expect(screen.getByRole('heading', { level: 2, name: 'Set 1' })).toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('renders blank form on load with defaults - Office page', () => {
      render(<GunSafeForm {...defaultProps} appName={APP_NAME.OFFICE} />, { wrapper: MockProviders });

      expect(screen.getByRole('heading', { level: 2, name: 'Set 1' })).toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });
  });

  describe('validates', () => {
    it('when all required fields are filled', async () => {
      render(<GunSafeForm {...defaultProps} {...gunSafeWithDocumentProps} />, { wrapper: MockProviders });
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
      });
    });
  });

  describe('attaches button handler callbacks', () => {
    it('calls the onSubmit callback with required fields', async () => {
      const expectedPayload = {
        document: gunSafeWithDocumentProps.gunSafe.document.uploads,
        weight: '1',
        description: 'Description',
        missingWeightTicket: false,
      };
      render(<GunSafeForm {...defaultProps} {...gunSafeWithDocumentProps} />, {
        wrapper: MockProviders,
      });
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(expectedPayload, expect.anything());
      });
    });

    it('calls the onBack prop when the Cancel button is clicked', async () => {
      render(<GunSafeForm {...defaultProps} />, { wrapper: MockProviders });

      await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });

    it('calls the onBack prop when the Cancel button is clicked - Office page', async () => {
      render(<GunSafeForm {...defaultProps} appName={APP_NAME.OFFICE} />, { wrapper: MockProviders });

      await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });

  describe('handles entitlements', () => {
    it("displays self's gun safe maximum. - Customer page", () => {
      const { container } = render(<GunSafeForm {...defaultProps} {...gunSafeProps} appName={APP_NAME.MYMOVE} />, {
        wrapper: MockProviders,
      });
      expect(container).toHaveTextContent(`Your maximum allowance is ${defaultProps.entitlements.gunSafeWeight} lbs.`);
    });

    it("displays self's gun safe maximum. - Office page", () => {
      const { container } = render(<GunSafeForm {...defaultProps} {...gunSafeProps} appName={APP_NAME.OFFICE} />, {
        wrapper: MockProviders,
      });
      expect(container).toHaveTextContent(`Your maximum allowance is ${defaultProps.entitlements.gunSafeWeight} lbs.`);
    });

    it('invalidates if weight exceeds the maximum. - Customer page', async () => {
      render(<GunSafeForm {...defaultProps} {...gunSafeProps} appName={APP_NAME.MYMOVE} />, {
        wrapper: MockProviders,
      });
      await act(async () => {
        await userEvent.type(screen.getByRole('textbox', { name: /^Shipment's gun safe weight/ }), '2000');
      });
      await waitFor(() => {
        expect(screen.getByText(/Weight must be lower than 450 lbs./)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });
    });

    it('invalidates if weight exceeds the maximum. - Office page', async () => {
      render(<GunSafeForm {...defaultProps} {...gunSafeProps} appName={APP_NAME.OFFICE} />, {
        wrapper: MockProviders,
      });
      await act(async () => {
        await userEvent.type(screen.getByRole('textbox', { name: /^Shipment's gun safe weight/ }), '8,501');
      });
      await waitFor(() => {
        expect(screen.getByText(/Weight must be lower than 450 lbs./)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });
    });
  });

  describe('invalidates fields', () => {
    it('invalidates if weight is zero', async () => {
      render(<GunSafeForm {...defaultProps} {...gunSafeNoWeightProps} />, { wrapper: MockProviders });
      await userEvent.type(screen.getByRole('textbox', { name: /^Shipment's gun safe weight/ }), '0');
      await waitFor(() => {
        expect(screen.getByText(/Enter a weight greater than 0 lbs./)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });
    });
  });
});
