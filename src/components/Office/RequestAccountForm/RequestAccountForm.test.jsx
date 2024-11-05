import React from 'react';
import { screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';
import selectEvent from 'react-select-event';

import RequestAccountForm from './RequestAccountForm';

import { renderWithRouter } from 'testUtils';
import { searchTransportationOfficesOpen } from 'services/ghcApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  searchTransportationOfficesOpen: jest.fn(),
}));

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
  };

  it('renders the form inputs', async () => {
    renderWithRouter(<RequestAccountForm {...testProps} />);

    const firstName = await screen.findByLabelText('First Name');
    expect(firstName).toBeInstanceOf(HTMLInputElement);
    expect(firstName).toHaveValue(testProps.initialValues.officeAccountRequestFirstName);

    const middleInitial = await screen.findByLabelText('First Name');
    expect(middleInitial).toBeInstanceOf(HTMLInputElement);
    expect(middleInitial).toHaveValue(testProps.initialValues.officeAccountRequestMiddleInitial);

    const lastName = await screen.findByLabelText('Last Name');
    expect(lastName).toBeInstanceOf(HTMLInputElement);
    expect(lastName).toHaveValue(testProps.initialValues.officeAccountRequestLastName);

    const email = await screen.findByLabelText('Email');
    expect(email).toBeInstanceOf(HTMLInputElement);
    expect(email).toHaveValue(testProps.initialValues.officeAccountRequestEmail);

    const telephone = await screen.findByLabelText('Telephone');
    expect(telephone).toBeInstanceOf(HTMLInputElement);
    expect(telephone).toHaveValue(testProps.initialValues.officeAccountRequestTelephone);

    const edipi = screen.getByTestId('officeAccountRequestEdipi');
    expect(edipi).toBeInstanceOf(HTMLInputElement);
    expect(edipi).toHaveValue(testProps.initialValues.officeAccountRequestEdipi);

    const uniqueId = screen.getByTestId('officeAccountRequestOtherUniqueId');
    expect(uniqueId).toBeInstanceOf(HTMLInputElement);
    expect(uniqueId).toHaveValue(testProps.initialValues.officeAccountRequestOtherUniqueId);

    const transportationOffice = screen.getByLabelText('Transportation Office');
    expect(transportationOffice).toBeInstanceOf(HTMLInputElement);
    expect(transportationOffice).toHaveTextContent('');

    const hqCheckbox = screen.getByTestId('headquartersCheckBox');
    expect(hqCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(hqCheckbox).not.toBeChecked(false);

    const tooCheckbox = screen.getByTestId('taskOrderingOfficerCheckBox');
    expect(tooCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(tooCheckbox).not.toBeChecked(false);

    const tioCheckbox = screen.getByTestId('taskInvoicingOfficerCheckBox');
    expect(tioCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(tioCheckbox).not.toBeChecked(false);

    const tcoCheckbox = screen.getByTestId('transportationContractingOfficerCheckBox');
    expect(tcoCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(tcoCheckbox).not.toBeChecked(false);

    const scCheckbox = screen.getByTestId('servicesCounselorCheckBox');
    expect(scCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(scCheckbox).not.toBeChecked(false);

    const qsaCheckbox = screen.getByTestId('qualityAssuranceEvaluatorCheckBox');
    expect(qsaCheckbox).toBeInstanceOf(HTMLInputElement);
    expect(qsaCheckbox).not.toBeChecked(false);
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

    await userEvent.type(screen.getByLabelText('First Name'), 'Bob');
    await userEvent.type(screen.getByLabelText('Last Name'), 'Banks');
    await userEvent.type(screen.getByLabelText('Email'), 'banks@us.navy.mil');
    await userEvent.type(screen.getByLabelText('Telephone'), '333-333-3333');
    await userEvent.type(screen.getByTestId('officeAccountRequestEdipi'), '1111111111');
    await userEvent.type(screen.getByTestId('officeAccountRequestOtherUniqueId'), '1111111111');

    const transportationOfficeInput = screen.getByLabelText('Transportation Office');
    await fireEvent.change(transportationOfficeInput, { target: { value: 'Tester' } });
    await act(() => selectEvent.select(transportationOfficeInput, /Tester/));

    const tooCheckbox = screen.getByTestId('taskOrderingOfficerCheckBox');
    await userEvent.click(tooCheckbox);

    const submitButton = await screen.getByTestId('requestOfficeAccountSubmitButton');
    await userEvent.click(submitButton);

    expect(testProps.onSubmit).toHaveBeenCalled();
  });

  it('Throws error requesting office account with invalid email domanin', async () => {
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

    await userEvent.type(screen.getByLabelText('First Name'), 'Bob');
    await userEvent.type(screen.getByLabelText('Last Name'), 'Banks');
    await userEvent.type(screen.getByLabelText('Email'), 'banks@gmail.com');

    const tooCheckbox = screen.getByTestId('taskOrderingOfficerCheckBox');
    await userEvent.click(tooCheckbox);

    expect(screen.getAllByText('Domain must be .mil, .gov or .edu').length).toBe(1);
  });

  it('shows policy error when both TOO and TIO checkboxes are both selected, and goes away after unselecting one of them', async () => {
    renderWithRouter(<RequestAccountForm {...testProps} />);

    const tooCheckbox = screen.getByTestId('taskOrderingOfficerCheckBox');
    const tioCheckbox = screen.getByTestId('taskInvoicingOfficerCheckBox');

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
