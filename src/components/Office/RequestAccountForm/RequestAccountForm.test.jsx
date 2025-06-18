import React from 'react';
import { screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';
import selectEvent from 'react-select-event';

import RequestAccountForm from './RequestAccountForm';

import { renderWithRouter } from 'testUtils';
import { searchTransportationOfficesOpen } from 'services/ghcApi';
import { isBooleanFlagEnabledUnauthenticatedOffice } from 'utils/featureFlags';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  searchTransportationOfficesOpen: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabledUnauthenticatedOffice: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const mockRolesWithPrivs = [
  { roleType: 'headquarters', roleName: 'Headquarters' },
  { roleType: 'task_ordering_officer', roleName: 'Task Ordering Officer' },
  { roleType: 'task_invoicing_officer', roleName: 'Task Invoicing Officer' },
  { roleType: 'contracting_officer', roleName: 'Contracting Officer' },
  { roleType: 'services_counselor', roleName: 'Services Counselor' },
  { roleType: 'qae', roleName: 'Quality Assurance Evaluator' },
  { roleType: 'customer_service_representative', roleName: 'Customer Service Representative' },
  { roleType: 'gsr', roleName: 'Government Surveillance Representative' },
];
const mockPrivileges = [{ privilegeType: 'supervisor', privilegeName: 'Supervisor' }];

describe('RequestAccountForm component', () => {
  const testProps = {
    initialValues: {
      officeAccountRequestFirstName: '',
      officeAccountRequestMiddleInitial: '',
      officeAccountRequestLastName: '',
      officeAccountRequestEmail: '',
      officeAccountRequestTelephone: '',
      officeAccountRequestEdipi: '',
      officeAccountRequestOtherUniqueId: '',
      officeAccountTransportationOffice: undefined,
    },
    onSubmit: jest.fn(),
    onCancel: jest.fn(),
    rolesWithPrivs: mockRolesWithPrivs,
    privileges: mockPrivileges,
  };

  it('renders the form inputs', async () => {
    isBooleanFlagEnabledUnauthenticatedOffice.mockImplementation(() => Promise.resolve(true));
    renderWithRouter(<RequestAccountForm {...testProps} />);

    const firstName = screen.getByTestId('officeAccountRequestFirstName');
    expect(firstName).toBeInstanceOf(HTMLInputElement);
    expect(firstName).toHaveValue(testProps.initialValues.officeAccountRequestFirstName);

    const middleInitial = screen.getByTestId('officeAccountRequestMiddleInitial');
    expect(middleInitial).toBeInstanceOf(HTMLInputElement);
    expect(middleInitial).toHaveValue(testProps.initialValues.officeAccountRequestMiddleInitial);

    const lastName = screen.getByTestId('officeAccountRequestLastName');
    expect(lastName).toBeInstanceOf(HTMLInputElement);
    expect(lastName).toHaveValue(testProps.initialValues.officeAccountRequestLastName);

    const email = screen.getByTestId('officeAccountRequestEmail');
    expect(email).toBeInstanceOf(HTMLInputElement);
    expect(email).toHaveValue(testProps.initialValues.officeAccountRequestEmail);

    const emailConfirmation = screen.getByTestId('emailConfirmation');
    expect(emailConfirmation).toBeInstanceOf(HTMLInputElement);
    expect(emailConfirmation).toHaveValue(testProps.initialValues.officeAccountRequestEmail);

    const telephone = screen.getByTestId('officeAccountRequestTelephone');
    expect(telephone).toBeInstanceOf(HTMLInputElement);
    expect(telephone).toHaveValue(testProps.initialValues.officeAccountRequestTelephone);

    const edipi = screen.getByTestId('officeAccountRequestEdipi');
    expect(edipi).toBeInstanceOf(HTMLInputElement);
    expect(edipi).toHaveValue(testProps.initialValues.officeAccountRequestEdipi);

    const edipiConfirmation = screen.getByTestId('edipiConfirmation');
    expect(edipiConfirmation).toBeInstanceOf(HTMLInputElement);
    expect(edipiConfirmation).toHaveValue(testProps.initialValues.officeAccountRequestEdipi);

    const uniqueId = screen.getByTestId('officeAccountRequestOtherUniqueId');
    expect(uniqueId).toBeInstanceOf(HTMLInputElement);
    expect(uniqueId).toHaveValue(testProps.initialValues.officeAccountRequestOtherUniqueId);

    const uniqueIdConfirmation = screen.getByTestId('otherUniqueIdConfirmation');
    expect(uniqueIdConfirmation).toBeInstanceOf(HTMLInputElement);
    expect(uniqueIdConfirmation).toHaveValue(testProps.initialValues.officeAccountRequestOtherUniqueId);

    const transportationOffice = screen.getByLabelText(/^Transportation Office/i);
    expect(transportationOffice).toBeInstanceOf(HTMLInputElement);
    expect(transportationOffice).toHaveTextContent('');

    const hqCheckbox = screen.getByTestId('headquartersCheckbox');
    expect(hqCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(hqCheckbox).not.toBeChecked(false);

    const tooCheckbox = screen.getByTestId('task_ordering_officerCheckbox');
    expect(tooCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(tooCheckbox).not.toBeChecked(false);

    const tioCheckbox = screen.getByTestId('task_invoicing_officerCheckbox');
    expect(tioCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(tioCheckbox).not.toBeChecked(false);

    const coCheckbox = screen.getByTestId('contracting_officerCheckbox');
    expect(coCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(coCheckbox).not.toBeChecked(false);

    const csrCheckbox = screen.getByTestId('customer_service_representativeCheckbox');
    expect(csrCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(csrCheckbox).not.toBeChecked(false);

    const scCheckbox = screen.getByTestId('services_counselorCheckbox');
    expect(scCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(scCheckbox).not.toBeChecked(false);

    const qsaCheckbox = screen.getByTestId('qaeCheckbox');
    expect(qsaCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(qsaCheckbox).not.toBeChecked(false);

    const gsrCheckbox = screen.getByTestId('gsrCheckbox');
    expect(gsrCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(gsrCheckbox).not.toBeChecked(false);

    await waitFor(() => {
      const supervisorPrivilegeCheckbox = screen.getByTestId('supervisorPrivilegeCheckbox');
      expect(supervisorPrivilegeCheckbox).toBeInstanceOf(HTMLInputElement);
      expect(supervisorPrivilegeCheckbox).not.toBeChecked(false);
    });
  });

  it('cancels requesting office account when cancel button is clicked', async () => {
    renderWithRouter(<RequestAccountForm {...testProps} />);

    const cancelButton = await screen.getByTestId('requestOfficeAccountCancelButton');
    await userEvent.click(cancelButton);
    expect(testProps.onCancel).toHaveBeenCalled();
  });

  it('submits requesting office account form when submit button is clicked', async () => {
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
    searchTransportationOfficesOpen.mockImplementation(mockSearchTransportationOfficesOpen);

    renderWithRouter(<RequestAccountForm {...testProps} />);

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

    const submitButton = await screen.getByTestId('requestOfficeAccountSubmitButton');
    await userEvent.click(submitButton);

    expect(testProps.onSubmit).toHaveBeenCalled();
  });

  it('submits requesting office account with supervisor privilege form when submit button is clicked', async () => {
    isBooleanFlagEnabledUnauthenticatedOffice.mockImplementation(() => Promise.resolve(true));
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
    searchTransportationOfficesOpen.mockImplementation(mockSearchTransportationOfficesOpen);

    renderWithRouter(<RequestAccountForm {...testProps} />);

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

    await waitFor(() => {
      const supervisorPrivilegeCheckbox = screen.getByTestId('supervisorPrivilegeCheckbox');
      userEvent.click(supervisorPrivilegeCheckbox);
    });

    const submitButton = await screen.getByTestId('requestOfficeAccountSubmitButton');
    await userEvent.click(submitButton);

    expect(testProps.onSubmit).toHaveBeenCalled();
  });

  it('Throws error requesting office account with invalid email domain', async () => {
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
    searchTransportationOfficesOpen.mockImplementation(mockSearchTransportationOfficesOpen);

    renderWithRouter(<RequestAccountForm {...testProps} />);

    await userEvent.type(screen.getByTestId('officeAccountRequestFirstName'), 'Bob');
    await userEvent.type(screen.getByTestId('officeAccountRequestLastName'), 'Banks');
    await userEvent.type(screen.getByTestId('officeAccountRequestEmail'), 'banks@gmail.com');

    const tooCheckbox = screen.getByTestId('task_ordering_officerCheckbox');
    await userEvent.click(tooCheckbox);

    expect(screen.getAllByText('Domain must be .mil, .gov or .edu').length).toBe(1);
  });

  describe('Role selection validation', () => {
    const checkboxTestIds = [
      'headquartersCheckbox',
      'task_ordering_officerCheckbox',
      'task_invoicing_officerCheckbox',
      'contracting_officerCheckbox',
      'services_counselorCheckbox',
      'qaeCheckbox',
      'customer_service_representativeCheckbox',
      'gsrCheckbox',
    ];

    it.each(checkboxTestIds)('shows and clears error for %s', async (testId) => {
      renderWithRouter(<RequestAccountForm {...testProps} />);

      const checkbox = screen.getByTestId(testId);

      await userEvent.click(checkbox); // check
      await userEvent.click(checkbox); // uncheck to trigger validation

      const error = await screen.findByText('You must select at least one role.');
      expect(error).toBeInTheDocument();

      await userEvent.click(checkbox); // check again
      expect(screen.queryByText('You must select at least one role.')).not.toBeInTheDocument();
    });
  });

  it('shows policy error when both TOO and TIO checkboxes are both selected, and goes away after unselecting one of them', async () => {
    renderWithRouter(<RequestAccountForm {...testProps} />);

    const tooCheckbox = screen.getByTestId('task_ordering_officerCheckbox');
    const tioCheckbox = screen.getByTestId('task_invoicing_officerCheckbox');

    // Click both the TOO and TIO role checkboxes
    await userEvent.click(tooCheckbox);
    await userEvent.click(tioCheckbox);

    // Check that the validation error appears
    const policyVerrs = await screen.findAllByText(
      'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.',
    );
    expect(policyVerrs.length).toBeGreaterThan(0);

    // Check that it goes away after unselecting either TIO or TOO checkbox
    await userEvent.click(tioCheckbox);
    expect(
      screen.queryByText(
        'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.',
      ),
    ).not.toBeInTheDocument();
  });

  afterEach(jest.resetAllMocks);
});
