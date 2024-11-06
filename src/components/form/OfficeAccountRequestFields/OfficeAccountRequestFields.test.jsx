import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import OfficeAccountRequestFields from './OfficeAccountRequestFields';

describe('OfficeAccountRequestFields component', () => {
  it('renders the form inputs', async () => {
    render(
      <Formik>
        <OfficeAccountRequestFields />
      </Formik>,
    );

    const firstName = await screen.findByLabelText('First Name');
    expect(firstName).toBeInstanceOf(HTMLInputElement);

    const middleInitial = await screen.findByLabelText('First Name');
    expect(middleInitial).toBeInstanceOf(HTMLInputElement);

    const lastName = await screen.findByLabelText('Last Name');
    expect(lastName).toBeInstanceOf(HTMLInputElement);

    const email = await screen.findByLabelText('Email');
    expect(email).toBeInstanceOf(HTMLInputElement);

    const telephone = await screen.findByLabelText('Telephone');
    expect(telephone).toBeInstanceOf(HTMLInputElement);

    const edipi = await screen.getByTestId('officeAccountRequestEdipi');
    expect(edipi).toBeInstanceOf(HTMLInputElement);

    const uniqueId = await screen.getByTestId('officeAccountRequestOtherUniqueId');
    expect(uniqueId).toBeInstanceOf(HTMLInputElement);

    const transportationOffice = await screen.getByLabelText('Transportation Office');
    expect(transportationOffice).toBeInstanceOf(HTMLInputElement);

    const hqCheckbox = await screen.getByTestId('headquartersCheckBox');
    expect(hqCheckbox).toBeInstanceOf(HTMLInputElement);

    const tooCheckbox = await screen.getByTestId('taskOrderingOfficerCheckBox');
    expect(tooCheckbox).toBeInstanceOf(HTMLInputElement);

    const tioCheckbox = await screen.getByTestId('taskInvoicingOfficerCheckBox');
    expect(tioCheckbox).toBeInstanceOf(HTMLInputElement);

    const tcoCheckbox = await screen.getByTestId('transportationContractingOfficerCheckBox');
    expect(tcoCheckbox).toBeInstanceOf(HTMLInputElement);

    const scCheckbox = await screen.getByTestId('servicesCounselorCheckBox');
    expect(scCheckbox).toBeInstanceOf(HTMLInputElement);

    const qaeCheckbox = await screen.getByTestId('qualityAssuranceEvaluatorCheckBox');
    expect(qaeCheckbox).toBeInstanceOf(HTMLInputElement);

    const csrCheckbox = await screen.getByTestId('customerSupportRepresentativeCheckBox');
    expect(csrCheckbox).toBeInstanceOf(HTMLInputElement);
  });
});
