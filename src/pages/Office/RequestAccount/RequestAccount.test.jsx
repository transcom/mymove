import React from 'react';
import { MemoryRouter } from 'react-router';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';
import selectEvent from 'react-select-event';

import { RequestAccount } from './RequestAccount';

import { generalRoutes } from 'constants/routes';
import { createOfficeAccountRequest, searchTransportationOfficesOpen } from 'services/ghcApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  createOfficeAccountRequest: jest.fn(),
  searchTransportationOfficesOpen: jest.fn(),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

const mockOfficeId = '3210a533-19b8-4805-a564-7eb452afce10';
const mockTransportationOffice = {
  address: {
    city: 'Test City',
    country: 'United States',
    id: 'a13806fc-0e7d-4dc3-91ca-b802d9da50f1',
    postalCode: '85309',
    state: 'AZ',
    streetAddress1: '7383 N Litchfield Rd',
    streetAddress2: 'Rm 1122',
  },
  created_at: '2018-05-28T14:27:39.198Z',
  gbloc: 'KKFA',
  id: mockOfficeId,
  name: 'Tester',
  phone_lines: [],
  updated_at: '2018-05-28T14:27:39.198Z',
};

const mockSearchTransportationOfficesOpen = () => Promise.resolve([mockTransportationOffice]);

describe('RequestAccount page', () => {
  it('renders the RequestAccount form', async () => {
    render(
      <MemoryRouter>
        <RequestAccount />
      </MemoryRouter>,
    );

    const formHeader = screen.getByRole('heading', { name: 'Request Office Account', level: 2 });
    expect(formHeader).toBeInTheDocument();
  });

  it('should navigate to sign in page after submit', async () => {
    const props = {
      setFlashMessage: jest.fn(),
    };
    render(
      <MemoryRouter>
        <RequestAccount {...props} />
      </MemoryRouter>,
    );

    const mockResponse = {
      ok: true,
      status: 200,
    };

    searchTransportationOfficesOpen.mockImplementation(mockSearchTransportationOfficesOpen);
    createOfficeAccountRequest.mockImplementation(() => Promise.resolve(mockResponse));

    await userEvent.type(screen.getByLabelText('First Name'), 'Bob');
    await userEvent.type(screen.getByLabelText('Last Name'), 'Banks');
    await userEvent.type(screen.getByLabelText('Email'), 'banks@us.af.mil');
    await userEvent.type(screen.getByLabelText('Telephone'), '333-333-3333');
    await userEvent.type(screen.getByTestId('officeAccountRequestEdipi'), '1111111111');
    await userEvent.type(screen.getByTestId('officeAccountRequestOtherUniqueId'), '1111111111');

    const transportationOfficeInput = screen.getByLabelText('Transportation Office');
    await fireEvent.change(transportationOfficeInput, { target: { value: 'Tester' } });
    await act(() => selectEvent.select(transportationOfficeInput, /Tester/));

    const tooCheckbox = screen.getByTestId('taskOrderingOfficerCheckBox');
    await userEvent.click(tooCheckbox);

    const saveBtn = screen.getByTestId('requestOfficeAccountSubmitButton');
    await userEvent.click(saveBtn);

    expect(mockNavigate).toHaveBeenCalledWith(generalRoutes.SIGN_IN_PATH);
  });

  it('should display error message on failed submit', async () => {
    render(
      <MemoryRouter>
        <RequestAccount />
      </MemoryRouter>,
    );

    const mockResponse = {
      status: 500,
      response: {
        body: {
          detail: 'test',
          invalid_fields: {
            email: 'Test',
          },
        },
      },
    };

    searchTransportationOfficesOpen.mockImplementation(mockSearchTransportationOfficesOpen);
    createOfficeAccountRequest.mockImplementation(() => Promise.reject(mockResponse));

    await userEvent.type(screen.getByLabelText('First Name'), 'Bob');
    await userEvent.type(screen.getByLabelText('Last Name'), 'Banks');
    await userEvent.type(screen.getByLabelText('Email'), 'banks@test.edu');
    await userEvent.type(screen.getByLabelText('Telephone'), '333-333-3333');
    await userEvent.type(screen.getByTestId('officeAccountRequestEdipi'), '1111111111');
    await userEvent.type(screen.getByTestId('officeAccountRequestOtherUniqueId'), '1111111111');

    const transportationOfficeInput = screen.getByLabelText('Transportation Office');
    await fireEvent.change(transportationOfficeInput, { target: { value: 'Tester' } });
    await act(() => selectEvent.select(transportationOfficeInput, /Tester/));

    const tcoCheckbox = screen.getByTestId('transportationContractingOfficerCheckBox');
    await userEvent.click(tcoCheckbox);

    const scCheckbox = screen.getByTestId('servicesCounselorCheckBox');
    await userEvent.click(scCheckbox);

    const qsaCheckbox = screen.getByTestId('qualityAssuranceEvaluatorCheckBox');
    await userEvent.click(qsaCheckbox);

    const gsrCheckbox = screen.getByTestId('governmentSurveillanceRepresentativeCheckbox');
    await userEvent.click(gsrCheckbox);

    const saveBtn = screen.getByTestId('requestOfficeAccountSubmitButton');
    await userEvent.click(saveBtn);

    expect(await screen.findByText('An error occurred')).toBeVisible();
  });

  it('goes back to the sign in page when the cancel button is clicked', async () => {
    render(
      <MemoryRouter>
        <RequestAccount />
      </MemoryRouter>,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(generalRoutes.SIGN_IN_PATH);
  });

  afterEach(jest.resetAllMocks);
});
