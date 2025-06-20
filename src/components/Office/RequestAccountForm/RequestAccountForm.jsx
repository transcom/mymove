import React from 'react';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import classnames from 'classnames';

import requestAccountFormStyles from './RequestAccountForm.module.scss';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { OfficeAccountRequestFields } from 'components/form/OfficeAccountRequestFields/OfficeAccountRequestFields';
import '@trussworks/react-uswds/lib/index.css';
import { Form } from 'components/form/Form';
import { withContext } from 'shared/AppContext';
import { officeAccountRequestSchema } from 'utils/validation';

const RequestAccountForm = ({ initialValues, onSubmit, onCancel, rolesWithPrivs, privileges }) => {
  const sectionStyles = classnames(formStyles.formSection, requestAccountFormStyles.formSection);

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validateOnMount
      validationSchema={officeAccountRequestSchema}
    >
      {({ isValid, handleSubmit }) => {
        return (
          <Form className={classnames(formStyles.form, requestAccountFormStyles.form)}>
            <SectionWrapper className={sectionStyles}>
              <h2 className={requestAccountFormStyles.center}>Request Office Account</h2>
              <p>
                To request an office acount, please fill out the form below with all required fields. Once submitted, a
                MilMove admin will review your account. If approved, you will be notified and allowed to log in.
              </p>
              <p>
                <b>NOTE:</b> When filling out your DODID# or unique identifier, you <b>MUST</b> provide one or the
                other. If using CAC, fill out the EDIPI field. If using PIV or ECA certs, please fill out the unique
                identifier.
              </p>

              <OfficeAccountRequestFields rolesWithPrivs={rolesWithPrivs} privileges={privileges} />
            </SectionWrapper>

            <div className={requestAccountFormStyles.buttonRow}>
              <Button
                type="button"
                onClick={() => onCancel()}
                data-testid="requestOfficeAccountCancelButton"
                outline
                secondary
              >
                Cancel
              </Button>
              <Button
                type="button"
                disabled={!isValid}
                onClick={() => handleSubmit()}
                data-testid="requestOfficeAccountSubmitButton"
              >
                Submit
              </Button>
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

export default withContext(RequestAccountForm);
