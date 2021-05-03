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
    expect(screen.getByLabelText('Name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Email')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Phone')).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all backup contact info inputs', () => {
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
      expect(screen.getByLabelText('Name')).toHaveValue(initialValues.name);
      expect(screen.getByLabelText('Email')).toHaveValue(initialValues.email);
      expect(screen.getByLabelText('Phone')).toHaveValue(initialValues.telephone);
    });
  });

  it('can namespace fields', () => {
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
    expect(screen.getByLabelText('Name')).toHaveValue(initialBackupInfo.name);
    expect(screen.getByLabelText('Email')).toHaveValue(initialBackupInfo.email);
    expect(screen.getByLabelText('Phone')).toHaveValue(initialBackupInfo.telephone);
  });
});
