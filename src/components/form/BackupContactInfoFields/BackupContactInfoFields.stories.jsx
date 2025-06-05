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
      firstName: '',
      lastName: '',
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

export const WithInitialValues = () => (
  <Formik
    initialValues={{
      firstName: 'Peyton',
      lastName: 'Wing',
      email: 'pw@example.com',
      telephone: '915-555-8761',
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
      firstName: '',
      lastName: '',
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
                If we cannot reach you, who can we contact? Any person you assign as a backup contact must be 18 years
                of age or older.
              </p>
              {fields}
            </>
          )}
        />
      </Form>
    )}
  </Formik>
);

export const WithNamespacedFields = () => {
  const namespace = 'backup_contact';

  return (
    <Formik
      initialValues={{
        [namespace]: {
          firstName: '',
          lastName: '',
          telephone: '',
          email: '',
        },
      }}
    >
      {() => (
        <Form className={formStyles.form}>
          <BackupContactInfoFields legend="Backup contact" name={namespace} />
        </Form>
      )}
    </Formik>
  );
};
