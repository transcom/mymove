import React from 'react';
import { connect } from 'react-redux';
import { func, PropTypes } from 'prop-types';
import { Formik } from 'formik';
import { useNavigate } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import SelectableCard from 'components/Customer/SelectableCard';
import { setConusStatus } from 'store/onboarding/actions';
import { selectConusStatus } from 'store/onboarding/selectors';
import { CONUS_STATUS } from 'shared/constants';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { Form } from 'components/form/Form';
import { customerRoutes } from 'constants/routes';

const ConusOrNot = ({ setLocation, conusStatus }) => {
  const navigate = useNavigate();

  const oconusCardText = (
    <>
      <div>Starts or ends in Alaska, Hawaii, or International locations</div>
      <strong>MilMove does not support OCONUS moves yet.</strong> Contact your current transportation office to set up
      your move.
    </>
  );

  const onSubmit = (values) => {
    // const payload = {
    //   id: serviceMember.id,
    //   code: values.code,
    // };

    console.log('values', values);

    navigate(customerRoutes.DOD_INFO_PATH);

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

  return (
    <GridContainer>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <Formik validateOnMount onSubmit={onSubmit}>
            {({ isValid, handleSubmit }) => {
              return (
                <Form className={formStyles.form}>
                  <h1>Where are you moving?</h1>
                  <SectionWrapper className={formStyles.formSection}>
                    <SelectableCard
                      id={`input_${CONUS_STATUS.CONUS}`}
                      label="CONUS"
                      value={CONUS_STATUS.CONUS}
                      onChange={(e) => setLocation(e.target.value)}
                      name="conusStatus"
                      checked={conusStatus === CONUS_STATUS.CONUS}
                      cardText="Starts and ends in the continental US"
                    />
                    <SelectableCard
                      id={`input_${CONUS_STATUS.OCONUS}`}
                      label="OCONUS"
                      value={CONUS_STATUS.OCONUS}
                      onChange={(e) => setLocation(e.target.value)}
                      name="conusStatus"
                      checked={conusStatus === CONUS_STATUS.OCONUS}
                      disabled
                      cardText={oconusCardText}
                    />
                  </SectionWrapper>
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

ConusOrNot.propTypes = {
  setLocation: func.isRequired,
  conusStatus: PropTypes.string,
};

ConusOrNot.defaultProps = {
  conusStatus: '',
};

const mapStateToProps = (state) => ({
  conusStatus: selectConusStatus(state),
});

const mapDispatchToProps = {
  setLocation: setConusStatus,
};

export default connect(mapStateToProps, mapDispatchToProps)(ConusOrNot);
