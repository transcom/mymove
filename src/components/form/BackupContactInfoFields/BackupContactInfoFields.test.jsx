import React from 'react';
import { render } from '@testing-library/react';
import { Formik } from 'formik';

import { BackupContactInfoFields } from './index';

describe('BackupContactInfoFields component', () => {
  it('renders a legend and all backup contact info inputs', () => {
    const { getByText, getByLabelText } = render(
      <Formik>
        <BackupContactInfoFields legend="Backup contact" name="backupContact" />
      </Formik>,
    );
    expect(getByText('Backup contact')).toBeInstanceOf(HTMLLegendElement);
    expect(getByLabelText('Name')).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText('Email')).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText('Phone')).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all backup contact info inputs', () => {
      const initialValues = {
        email: 'test@example.com',
        name: 'test',
        telephone: '555-123-4567',
      };

      const { getByLabelText } = render(
        <Formik initialValues={initialValues}>
          <BackupContactInfoFields legend="Backup contact" name="backupContact" />
        </Formik>,
      );
      expect(getByLabelText('Name')).toHaveValue(initialValues.name);
      expect(getByLabelText('Email')).toHaveValue(initialValues.email);
      expect(getByLabelText('Phone')).toHaveValue(initialValues.telephone);
    });
  });
});
