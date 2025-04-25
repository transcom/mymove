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
});
