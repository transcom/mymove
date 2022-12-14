import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import { MockProviders } from 'testUtils';
import ProGearForm from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm';

const defaultProps = {
  onBack: jest.fn(),
  onSubmit: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
  onCreateUpload: jest.fn(),
};

const mockProGearEntitlements = {
  proGear: 1234,
  proGearSpouse: 987,
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectProGearEntitlements: jest.fn().mockImplementation(() => mockProGearEntitlements),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const proGearProps = {
  proGear: {
    belongsToSelf: true,
    document: {},
    weight: 1,
    description: 'Description',
    hasWeightTickets: '',
  },
};

const proGearNoWeightProps = {
  proGear: {
    ...proGearProps.proGear,
    weight: 0,
  },
};

const spouseProGearProps = {
  proGear: {
    belongsToSelf: false,
  },
};

const proGearWithDocumentProps = {
  proGear: {
    ...proGearProps.proGear,
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

describe('ProGearForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', () => {
      render(<ProGearForm {...defaultProps} />, { wrapper: MockProviders });

      expect(screen.getByRole('heading', { level: 2, name: 'Set 1' })).toBeInTheDocument();
      expect(screen.getByText('Who does this pro-gear belong to?')).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getByLabelText('Me')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('My spouse')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('does not select a radio when belongsToSelf is null', () => {
      render(<ProGearForm {...defaultProps} />, { wrapper: MockProviders });
      expect(screen.getByLabelText('Me')).not.toBeChecked();
      expect(screen.getByLabelText('My spouse')).not.toBeChecked();
    });

    it('selects "Me" radio when belongsToSelf is true', () => {
      const { container } = render(<ProGearForm {...defaultProps} {...proGearProps} />, { wrapper: MockProviders });
      expect(screen.getByLabelText('Me')).toBeChecked();
      expect(screen.getByLabelText('My spouse')).not.toBeChecked();
      expect(container).toHaveTextContent("You have to separate yours and your spouse's pro-gear.");
    });

    it('selects "My spouse" radio when belongsToSelf is false', () => {
      render(<ProGearForm {...defaultProps} {...spouseProGearProps} />, { wrapper: MockProviders });
      expect(screen.getByLabelText('My spouse')).toBeChecked();
      expect(screen.getByLabelText('Me')).not.toBeChecked();
    });
  });
  describe('validates', () => {
    it('when all required fields are filled', async () => {
      render(<ProGearForm {...defaultProps} {...proGearWithDocumentProps} />, { wrapper: MockProviders });
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
      });
    });
  });
  describe('attaches button handler callbacks', () => {
    it('calls the onSubmit callback with required fields', async () => {
      const expectedPayload = {
        belongsToSelf: 'true',
        document: proGearWithDocumentProps.proGear.document.uploads,
        weight: '1',
        description: 'Description',
        missingWeightTicket: false,
      };
      render(<ProGearForm {...defaultProps} {...proGearWithDocumentProps} />, { wrapper: MockProviders });
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(expectedPayload, expect.anything());
      });
    });
    it('calls the onBack prop when the Return To Homepage button is clicked', async () => {
      render(<ProGearForm {...defaultProps} />, { wrapper: MockProviders });

      await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });
  describe('handles entitlements', () => {
    it("displays self's pro-gear maximum.", () => {
      const { container } = render(<ProGearForm {...defaultProps} {...proGearProps} />, { wrapper: MockProviders });
      expect(container).toHaveTextContent('Your maximum allowance is 1,234 lbs.');
    });
    it("displays spouse's pro-gear maximum.", () => {
      const { container } = render(<ProGearForm {...defaultProps} {...proGearProps} {...spouseProGearProps} />, {
        wrapper: MockProviders,
      });
      expect(container).toHaveTextContent('Your maximum allowance is 987 lbs.');
    });
    it('invalidates if weight exceeds the maximum.', async () => {
      render(<ProGearForm {...defaultProps} {...proGearProps} />, { wrapper: MockProviders });
      await act(async () => {
        await userEvent.type(screen.getByRole('textbox', { name: /^Shipment's pro-gear weight/ }), '2000');
      });
      await waitFor(() => {
        expect(screen.getByText(/Pro gear weight must be less than or equal to 1,234 lbs./)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });
    });
    it('invalidates if a valid weight is entered but a lower maximum is subsequently selected', async () => {
      render(<ProGearForm {...defaultProps} {...proGearProps} />, { wrapper: MockProviders });
      await act(async () => {
        await userEvent.clear(screen.getByRole('textbox', { name: /^Shipment's pro-gear weight/ }));
        await userEvent.type(screen.getByRole('textbox', { name: /^Shipment's pro-gear weight/ }), '1000');
      });
      await waitFor(() => {
        expect(screen.queryByText(/Pro gear weight must be less than or equal to 1,234 lbs./)).not.toBeInTheDocument();
      });
      await act(async () => {
        await userEvent.click(screen.getByLabelText('My spouse'));
      });
      await waitFor(() => {
        expect(screen.getByText(/Pro gear weight must be less than or equal to 987 lbs./)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });
    });
  });
  describe('invalidates fields', () => {
    it('invalidates if weight is zero', async () => {
      render(<ProGearForm {...defaultProps} {...proGearNoWeightProps} />, { wrapper: MockProviders });
      await userEvent.type(screen.getByRole('textbox', { name: /^Shipment's pro-gear weight/ }), '0');
      await waitFor(() => {
        expect(screen.getByText(/Enter a weight greater than 0 lbs./)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });
    });
  });
});
