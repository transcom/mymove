import React from 'react';
import { Formik } from 'formik';

import { BackupContactInfoFields } from './index';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Components/Fieldsets/BackupContactInfoFields',
};

export const Basic = () => (
  <Formik
    initialValues={{
      name: '',
      telephone: '',
      email: '',
    }}
  >
    {() => (
      <Form className={formStyles.form}>
        <BackupContactInfoFields legend="Backup contact" />
      </Form>
    )}
  </Formik>
);

export const WithAdditionalText = () => (
  <Formik
    initialValues={{
      name: '',
      telephone: '',
      email: '',
    }}
  >
    {() => (
      <Form className={formStyles.form}>
        <BackupContactInfoFields
          legend="Backup contact"
          render={(fields) => (
            <>
              <p>
                If we can&apos;t reach you, who can we contact? Any person you assign as a backup contact must be 18
                years of age or older.
              </p>
              {fields}
            </>
          )}
        />
      </Form>
    )}
  </Formik>
);
