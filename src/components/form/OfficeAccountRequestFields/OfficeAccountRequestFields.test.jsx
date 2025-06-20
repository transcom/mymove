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
    expect(screen.getByTestId('headquartersCheckBox')).toBeInTheDocument();
    expect(screen.getByTestId('taskOrderingOfficerCheckBox')).toBeInTheDocument();
    expect(screen.getByTestId('taskInvoicingOfficerCheckBox')).toBeInTheDocument();
    expect(screen.getByTestId('transportationContractingOfficerCheckBox')).toBeInTheDocument();
    expect(screen.getByTestId('servicesCounselorCheckBox')).toBeInTheDocument();
    expect(screen.getByTestId('qualityAssuranceEvaluatorCheckBox')).toBeInTheDocument();
    expect(screen.getByTestId('customerSupportRepresentativeCheckBox')).toBeInTheDocument();
    expect(screen.getByTestId('governmentSurveillanceRepresentativeCheckbox')).toBeInTheDocument();
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

    const headquartersCheckbox = screen.getByTestId('headquartersCheckBox');

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

    const tooCheckbox = screen.getByTestId('taskOrderingOfficerCheckBox');
    const tioCheckbox = screen.getByTestId('taskInvoicingOfficerCheckBox');

    await userEvent.click(tooCheckbox);
    await userEvent.click(tioCheckbox);

    expect(
      await screen.findByText(
        'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.',
      ),
    ).toBeInTheDocument();
  });

  it('renders asterisks for required fields', async () => {
    render(
      <Formik initialValues={initialValues} validationSchema={officeAccountRequestSchema}>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
    expect(screen.getByLabelText('First Name *')).toBeInTheDocument();
    expect(screen.getByLabelText('Last Name *')).toBeInTheDocument();
    expect(screen.getByLabelText('Email *')).toBeInTheDocument();
    expect(screen.getByLabelText('Confirm Email *')).toBeInTheDocument();
    expect(screen.getByLabelText('Telephone *')).toBeInTheDocument();
    expect(screen.getByLabelText('Transportation Office *')).toBeInTheDocument();
    expect(screen.getByTestId('requestedRolesHeadingSpan')).toBeInTheDocument();
    expect(screen.getByTestId('requestedRolesHeadingSpan')).toHaveTextContent('*');
  });
});
