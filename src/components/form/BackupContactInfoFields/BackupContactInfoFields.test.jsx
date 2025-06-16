import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import { BackupContactInfoFields } from './index';

describe('BackupContactInfoFields component', () => {
  it('renders a legend and all backup contact info inputs and asterisks for required fields', () => {
    render(
      <Formik>
        <BackupContactInfoFields legend="Backup contact" />
      </Formik>,
    );
    expect(screen.getByTestId('reqAsteriskMsg')).toBeInTheDocument();

    expect(screen.getByText('Backup contact')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('Name *')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText(/Name */)).toBeInTheDocument();
    expect(screen.getByLabelText(/Name */)).toBeRequired();
    expect(screen.getByLabelText('Email *')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText(/Email */)).toBeInTheDocument();
    expect(screen.getByLabelText(/Email */)).toBeRequired();
    expect(screen.getByLabelText('Phone *')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText(/Phone */)).toBeInTheDocument();
    expect(screen.getByLabelText(/Phone */)).toBeRequired();
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all backup contact info inputs', async () => {
      const initialValues = {
        email: 'test@example.com',
        name: 'test',
        telephone: '555-123-4567',
      };

      render(
        <Formik initialValues={initialValues}>
          <BackupContactInfoFields legend="Backup contact" />
        </Formik>,
      );
      expect(await screen.findByLabelText('Name *')).toHaveValue(initialValues.name);
      expect(await screen.findByLabelText('Email *')).toHaveValue(initialValues.email);
      expect(await screen.findByLabelText('Phone *')).toHaveValue(initialValues.telephone);
    });
  });

  it('can namespace fields', async () => {
    const namespace = 'backup_contact';

    const initialBackupInfo = {
      email: 'test@example.com',
      name: 'test',
      telephone: '555-123-4567',
    };

    const initialValues = {
      [namespace]: initialBackupInfo,
    };

    render(
      <Formik initialValues={initialValues}>
        <BackupContactInfoFields legend="Backup contact" name={namespace} />
      </Formik>,
    );

    expect(await screen.findByLabelText('Name *')).toHaveValue(initialBackupInfo.name);
    expect(await screen.findByLabelText('Email *')).toHaveValue(initialBackupInfo.email);
    expect(await screen.findByLabelText('Phone *')).toHaveValue(initialBackupInfo.telephone);
  });
});
