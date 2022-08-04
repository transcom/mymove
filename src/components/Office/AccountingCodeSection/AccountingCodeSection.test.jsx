import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import AccountingCodeSection from './AccountingCodeSection';

describe('components/Office/AccountingCodeSection', () => {
  it('should render without codes', () => {
    render(
      <Formik initialValues={{}}>
        <AccountingCodeSection
          label="Test Code Section"
          emptyMessage="No codes entered."
          fieldName="testSection"
          shipmentTypes={{}}
        />
      </Formik>,
    );

    expect(screen.getByText('No codes entered.')).toBeInTheDocument();
  });

  it('should render a single code and select it by default', async () => {
    render(
      <Formik initialValues={{}}>
        <AccountingCodeSection
          label="Test Code Section"
          emptyMessage="No codes entered."
          fieldName="testSection"
          shipmentTypes={{ HHG: '1234' }}
        />
      </Formik>,
    );

    await waitFor(() => expect(screen.getByLabelText('1234 (HHG)')).toBeChecked());
  });

  it('should render multiple codes and not default either of them', () => {
    render(
      <Formik initialValues={{}}>
        <AccountingCodeSection
          label="Test Code Section"
          emptyMessage="No codes entered."
          fieldName="testSection"
          shipmentTypes={{ HHG: '1234', NTS: '2345' }}
        />
      </Formik>,
    );

    expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
    expect(screen.getByLabelText('2345 (NTS)')).not.toBeChecked();
  });

  it('should respond to initialValues', () => {
    render(
      <Formik initialValues={{ testSection: 'NTS' }}>
        <AccountingCodeSection
          label="Test Code Section"
          emptyMessage="No codes entered."
          fieldName="testSection"
          shipmentTypes={{ HHG: '1234', NTS: '2345' }}
        />
      </Formik>,
    );

    expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
    expect(screen.getByLabelText('2345 (NTS)')).toBeChecked();
  });

  it('should clear values if the clear button is selected', async () => {
    render(
      <Formik initialValues={{ testSection: 'NTS' }}>
        <AccountingCodeSection
          label="Test Code Section"
          emptyMessage="No codes entered."
          fieldName="testSection"
          shipmentTypes={{ HHG: '1234', NTS: '2345' }}
        />
      </Formik>,
    );

    await userEvent.click(screen.getByText('Clear selection'));
    await waitFor(() => expect(screen.getByLabelText('2345 (NTS)')).not.toBeChecked());
  });
});
