import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import OfficeAccountRequestFields from './OfficeAccountRequestFields';

import { officeAccountRequestSchema } from 'utils/validation';

const initialValues = {
  officeAccountRequestEdipi: '',
  edipiConfirmation: '',
  officeAccountRequestOtherUniqueId: '',
  otherUniqueIdConfirmation: '',
  officeAccountRequestEmail: '',
  emailConfirmation: '',
};

describe('OfficeAccountRequestFields component', () => {
  it('renders the form inputs', async () => {
    render(
      <Formik initialValues={initialValues} validationSchema={officeAccountRequestSchema}>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    expect(screen.getByTestId('officeAccountRequestFirstName')).toBeInTheDocument();
    expect(screen.getByTestId('officeAccountRequestMiddleInitial')).toBeInTheDocument();
    expect(screen.getByTestId('officeAccountRequestLastName')).toBeInTheDocument();
    expect(screen.getByTestId('officeAccountRequestEmail')).toBeInTheDocument();
    expect(screen.getByTestId('emailConfirmation')).toBeInTheDocument();
    expect(screen.getByTestId('officeAccountRequestTelephone')).toBeInTheDocument();
    expect(screen.getByTestId('officeAccountRequestEdipi')).toBeInTheDocument();
    expect(screen.getByTestId('edipiConfirmation')).toBeInTheDocument();
    expect(screen.getByTestId('officeAccountRequestOtherUniqueId')).toBeInTheDocument();
    expect(screen.getByTestId('otherUniqueIdConfirmation')).toBeInTheDocument();
    expect(screen.getByTestId('headquartersCheckbox')).toBeInTheDocument();
    expect(screen.getByTestId('task_ordering_officerCheckbox')).toBeInTheDocument();
    expect(screen.getByTestId('task_invoicing_officerCheckbox')).toBeInTheDocument();
    expect(screen.getByTestId('contracting_officerCheckbox')).toBeInTheDocument();
    expect(screen.getByTestId('services_counselorCheckbox')).toBeInTheDocument();
    expect(screen.getByTestId('qaeCheckbox')).toBeInTheDocument();
    expect(screen.getByTestId('customer_service_representativeCheckbox')).toBeInTheDocument();
    expect(screen.getByTestId('gsrCheckbox')).toBeInTheDocument();
  });

  it('validates that EDIPI and EDIPI confirmation match', async () => {
    render(
      <Formik initialValues={initialValues} validationSchema={officeAccountRequestSchema}>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    const edipiInput = screen.getByTestId('officeAccountRequestEdipi');
    const edipiConfirmInput = screen.getByTestId('edipiConfirmation');

    await userEvent.type(edipiInput, '1234567890');
    await userEvent.type(edipiConfirmInput, '0987654321');
    await userEvent.tab();

    expect(await screen.findByText('DODID#s must match')).toBeInTheDocument();

    await userEvent.clear(edipiConfirmInput);
    await userEvent.type(edipiConfirmInput, '1234567890');
    await userEvent.tab();

    expect(screen.queryByText('DODID#s must match')).not.toBeInTheDocument();
  });

  it('validates that Other Unique ID and its confirmation match', async () => {
    render(
      <Formik initialValues={initialValues} validationSchema={officeAccountRequestSchema}>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    const uniqueIdInput = screen.getByTestId('officeAccountRequestOtherUniqueId');
    const uniqueIdConfirmInput = screen.getByTestId('otherUniqueIdConfirmation');

    await userEvent.type(uniqueIdInput, 'ABCD1234');
    await userEvent.type(uniqueIdConfirmInput, 'XYZ9876');
    await userEvent.tab();

    expect(await screen.findByText('Unique IDs must match')).toBeInTheDocument();

    await userEvent.clear(uniqueIdConfirmInput);
    await userEvent.type(uniqueIdConfirmInput, 'ABCD1234');
    await userEvent.tab();

    expect(screen.queryByText('Unique IDs must match')).not.toBeInTheDocument();
  });

  it('validates that email and email confirmation match', async () => {
    render(
      <Formik initialValues={initialValues} validationSchema={officeAccountRequestSchema}>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    const emailInput = screen.getByTestId('officeAccountRequestEmail');
    const emailConfirmInput = screen.getByTestId('emailConfirmation');

    await userEvent.type(emailInput, 'test@example.com');
    await userEvent.type(emailConfirmInput, 'wrong@example.com');
    await userEvent.tab();

    expect(await screen.findByText('Emails must match')).toBeInTheDocument();

    await userEvent.clear(emailConfirmInput);
    await userEvent.type(emailConfirmInput, 'test@example.com');
    await userEvent.tab();

    expect(screen.queryByText('Emails must match')).not.toBeInTheDocument();
  });

  it('shows a validation error if no roles are selected after interaction', async () => {
    render(
      <Formik initialValues={initialValues} validationSchema={officeAccountRequestSchema}>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    const headquartersCheckbox = screen.getByTestId('headquartersCheckbox');

    await userEvent.click(headquartersCheckbox); // check
    await userEvent.click(headquartersCheckbox); // uncheck
    await userEvent.tab();

    expect(await screen.findByText('You must select at least one role.')).toBeInTheDocument();
  });

  it('shows a validation error if both Task Ordering and Task Invoicing Officer are selected', async () => {
    render(
      <Formik initialValues={initialValues} validationSchema={officeAccountRequestSchema}>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    const tooCheckbox = screen.getByTestId('task_ordering_officerCheckbox');
    const tioCheckbox = screen.getByTestId('task_invoicing_officerCheckbox');

    await userEvent.click(tooCheckbox);
    await userEvent.click(tioCheckbox);

    expect(
      await screen.findByText(
        'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.',
      ),
    ).toBeInTheDocument();
  });
});
