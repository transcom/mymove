import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import { BackupContactInfoFields } from './index';

describe('BackupContactInfoFields component', () => {
  it('renders a legend and all backup contact info inputs', () => {
    render(
      <Formik>
        <BackupContactInfoFields legend="Backup contact" />
      </Formik>,
    );
    expect(screen.getByText('Backup contact')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('First Name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Last Name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Email')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Phone')).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all backup contact info inputs', async () => {
      const initialValues = {
        email: 'test@example.com',
        firstName: 'test',
        lastName: 'case',
        telephone: '555-123-4567',
      };

      render(
        <Formik initialValues={initialValues}>
          <BackupContactInfoFields legend="Backup contact" />
        </Formik>,
      );
      expect(await screen.findByLabelText('First Name')).toHaveValue(initialValues.firstName);
      expect(await screen.findByLabelText('Last Name')).toHaveValue(initialValues.lastName);
      expect(await screen.findByLabelText('Email')).toHaveValue(initialValues.email);
      expect(await screen.findByLabelText('Phone')).toHaveValue(initialValues.telephone);
    });
  });

  it('can namespace fields', async () => {
    const namespace = 'backup_contact';

    const initialBackupInfo = {
      email: 'test@example.com',
      firstName: 'test',
      lastName: 'test',
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

    expect(await screen.findByLabelText('First Name')).toHaveValue(initialBackupInfo.firstName);
    expect(await screen.findByLabelText('Last Name')).toHaveValue(initialBackupInfo.lastName);
    expect(await screen.findByLabelText('Email')).toHaveValue(initialBackupInfo.email);
    expect(await screen.findByLabelText('Phone')).toHaveValue(initialBackupInfo.telephone);
  });
});
