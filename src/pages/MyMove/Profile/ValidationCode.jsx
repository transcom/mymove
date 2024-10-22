import React, { useState } from 'react';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { Formik } from 'formik';
import * as Yup from 'yup';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { getResponseError, validateCode } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import TextField from 'components/form/fields/TextField/TextField';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

export const ValidationCode = () => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);
  const [validationError, setValidationError] = useState(null);

  const initialValues = {
    code: '',
  };

  const onSubmit = async (values) => {
    const body = {
      parameterValue: values.code,
      parameterName: 'validation_code',
    };

    await validateCode(body)
      .then((response) => {
        const { parameterValue } = response.body;
        if (parameterValue === body.parameterValue) {
          navigate(customerRoutes.DOD_INFO_PATH);
        } else {
          setValidationError('Please try again');
        }
      })
      .catch((e) => {
        const { response } = e;
        const errorMessage = getResponseError(response, 'There was an internal server error submitting your request.');
        setServerError(errorMessage);
      });
  };

  const validationSchema = Yup.object().shape({
    code: Yup.string().required('Required').max(20, 'Enter up to 20 digits'),
  });

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}
      {validationError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="Incorrect validation code">
              {validationError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
            {({ isValid, handleSubmit }) => {
              return (
                <Form className={formStyles.form}>
                  <h1>Please enter a validation code to begin creating a move</h1>
                  <TextField
                    label="Validation code"
                    name="code"
                    id="code"
                    required
                    maxLength="20"
                    inputMode="numeric"
                    pattern="[0-9]{10}"
                  />

                  <div className={formStyles.formActions}>
                    <WizardNavigation isFirstPage disableNext={!isValid} onNextClick={handleSubmit} />
                  </div>
                </Form>
              );
            }}
          </Formik>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(ValidationCode);
