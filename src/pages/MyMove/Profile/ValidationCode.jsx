import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { Formik } from 'formik';
import * as Yup from 'yup';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { customerRoutes } from 'constants/routes';
import { ServiceMemberShape } from 'types/customerShapes';
import ValidationCodeForm from 'components/Customer/ValidationCodeForm/ValidationCodeForm';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import TextField from 'components/form/fields/TextField/TextField';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

export const ValidationCode = ({ updateServiceMember, serviceMember }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    code: '',
  };

  const handleBack = () => {
    navigate(customerRoutes.CONUS_OCONUS_PATH);
  };

  const handleNext = () => {
    navigate(customerRoutes.NAME_PATH);
  };

  const onSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      code: values.code,
    };

    console.log('payload', payload);

    navigate(customerRoutes.CONUS_OCONUS_PATH);

    // return patchServiceMember(payload)
    //   .then(updateServiceMember)
    //   .then(handleNext)
    //   .catch((e) => {
    //     // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
    //     const { response } = e;
    //     const errorMessage = getResponseError(response, 'failed to update service member due to server error');
    //     setServerError(errorMessage);
    //   });
  };

  const validationSchema = Yup.object().shape({
    code: Yup.string()
      .matches(/[0-9]{20}/, 'Enter a 20-digit number')
      .required('Required'),
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

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
            {({ isValid, handleSubmit }) => {
              return (
                <Form className={formStyles.form}>
                  <h1>Please enter validation code to begin creating a move</h1>
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

ValidationCode.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(ValidationCode);
