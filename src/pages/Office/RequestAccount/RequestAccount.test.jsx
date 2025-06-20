import React from 'react';
import { screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';
import selectEvent from 'react-select-event';

import { RequestAccount } from './RequestAccount';

import { generalRoutes } from 'constants/routes';
import { createOfficeAccountRequest, searchTransportationOfficesOpen } from 'services/ghcApi';
import { renderWithProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  useRolesPrivilegesQueriesOfficeApp: () => ({
    result: {
      privileges: [{ privilegeType: 'supervisor', privilegeName: 'Supervisor' }],
      rolesWithPrivs: [
        { roleType: 'headquarters', roleName: 'Headquarters' },
        { roleType: 'task_ordering_officer', roleName: 'Task Ordering Officer' },
        { roleType: 'task_invoicing_officer', roleName: 'Task Invoicing Officer' },
        { roleType: 'contracting_officer', roleName: 'Contracting Officer' },
        { roleType: 'services_counselor', roleName: 'Services Counselor' },
        { roleType: 'qae', roleName: 'Quality Assurance Evaluator' },
        { roleType: 'customer_service_representative', roleName: 'Customer Service Representative' },
        { roleType: 'gsr', roleName: 'Government Surveillance Representative' },
      ],
    },
  }),
}));

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
    renderWithProviders(<RequestAccount />);

    const formHeader = screen.getByRole('heading', { name: 'Request Office Account', level: 2 });
    expect(formHeader).toBeInTheDocument();
  });

  it('should navigate to sign in page after submit', async () => {
    const props = {
      setFlashMessage: jest.fn(),
    };
    renderWithProviders(<RequestAccount {...props} />);

    const mockResponse = {
      ok: true,
      status: 200,
    };

    searchTransportationOfficesOpen.mockImplementation(mockSearchTransportationOfficesOpen);
    createOfficeAccountRequest.mockImplementation(() => Promise.resolve(mockResponse));

    await userEvent.type(screen.getByTestId('officeAccountRequestFirstName'), 'Bob');
    await userEvent.type(screen.getByTestId('officeAccountRequestLastName'), 'Banks');
    await userEvent.type(screen.getByTestId('officeAccountRequestEmail'), 'banks@us.navy.mil');
    await userEvent.type(screen.getByTestId('emailConfirmation'), 'banks@us.navy.mil');
    await userEvent.type(screen.getByTestId('officeAccountRequestTelephone'), '333-333-3333');
    await userEvent.type(screen.getByTestId('officeAccountRequestEdipi'), '1111111111');
    await userEvent.type(screen.getByTestId('edipiConfirmation'), '1111111111');
    await userEvent.type(screen.getByTestId('officeAccountRequestOtherUniqueId'), 'uniqueID123');
    await userEvent.type(screen.getByTestId('otherUniqueIdConfirmation'), 'uniqueID123');

    const transportationOfficeInput = screen.getByLabelText(/^Transportation Office/i);
    await fireEvent.change(transportationOfficeInput, { target: { value: 'Tester' } });
    await act(() => selectEvent.select(transportationOfficeInput, /Tester/));

    const tooCheckbox = screen.getByTestId('task_ordering_officerCheckbox');
    await userEvent.click(tooCheckbox);

    const saveBtn = screen.getByTestId('requestOfficeAccountSubmitButton');
    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(props.setFlashMessage).toHaveBeenCalledWith(
        'OFFICE_ACCOUNT_REQUEST_SUCCESS',
        'success',
        'You have successfully requested access to MilMove. This request must be processed by an administrator prior to login. Once this process is completed, an approval or rejection email will be sent notifying you of the status of your account request.',
        '',
        true,
      );
    });
    expect(mockNavigate).toHaveBeenCalledWith(generalRoutes.SIGN_IN_PATH);
  });

  it('should display error message on failed submit', async () => {
    renderWithProviders(<RequestAccount />);

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

    await userEvent.type(screen.getByTestId('officeAccountRequestFirstName'), 'Bob');
    await userEvent.type(screen.getByTestId('officeAccountRequestLastName'), 'Banks');
    await userEvent.type(screen.getByTestId('officeAccountRequestEmail'), 'banks@us.navy.mil');
    await userEvent.type(screen.getByTestId('emailConfirmation'), 'banks@us.navy.mil');
    await userEvent.type(screen.getByTestId('officeAccountRequestTelephone'), '333-333-3333');
    await userEvent.type(screen.getByTestId('officeAccountRequestEdipi'), '1111111111');
    await userEvent.type(screen.getByTestId('edipiConfirmation'), '1111111111');
    await userEvent.type(screen.getByTestId('officeAccountRequestOtherUniqueId'), 'uniqueID123');
    await userEvent.type(screen.getByTestId('otherUniqueIdConfirmation'), 'uniqueID123');

    const transportationOfficeInput = screen.getByLabelText(/^Transportation Office/i);
    await fireEvent.change(transportationOfficeInput, { target: { value: 'Tester' } });
    await act(() => selectEvent.select(transportationOfficeInput, /Tester/));

    const coCheckbox = screen.getByTestId('contracting_officerCheckbox');
    await userEvent.click(coCheckbox);

    const scCheckbox = screen.getByTestId('services_counselorCheckbox');
    await userEvent.click(scCheckbox);

    const qaeCheckbox = screen.getByTestId('qaeCheckbox');
    await userEvent.click(qaeCheckbox);

    const gsrCheckbox = screen.getByTestId('gsrCheckbox');
    await userEvent.click(gsrCheckbox);

    const saveBtn = screen.getByTestId('requestOfficeAccountSubmitButton');
    await userEvent.click(saveBtn);

    expect(await screen.findByText('An error occurred')).toBeVisible();
  });

  it('goes back to the sign in page when the cancel button is clicked', async () => {
    renderWithProviders(<RequestAccount />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(generalRoutes.SIGN_IN_PATH);
  });

  afterEach(jest.resetAllMocks);
});
