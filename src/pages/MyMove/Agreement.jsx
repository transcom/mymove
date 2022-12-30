import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import moment from 'moment';
import { connect } from 'react-redux';
import { generatePath, useNavigate } from 'react-router-dom';

import { customerRoutes } from 'constants/routes';
import SubmitMoveForm from 'components/Customer/SubmitMoveForm/SubmitMoveForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import { submitMoveForApproval } from 'services/internalApi';
import { selectCurrentMove } from 'store/entities/selectors';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { formatSwaggerDate } from 'utils/formatters';

export const Agreement = ({ moveId, updateMove, setFlashMessage }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    signature: '',
    date: formatSwaggerDate(new Date()),
  };

  const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

  const handleBack = () => navigate(reviewPath);

  const handleSubmit = (values) => {
    const submitDate = moment().format();

    const data = {
      certification_text: completeCertificationText,
      date: submitDate,
      signature: values.signature,
      certification_type: SIGNED_CERT_OPTIONS.SHIPMENT,
    };

    submitMoveForApproval(moveId, data)
      .then((response) => {
        updateMove(response);
        setFlashMessage('MOVE_SUBMIT_SUCCESS', 'success', 'You’ve submitted your move request.');
        navigate('/');
      })
      .catch((error) => {
        // TODO - log error internally?
        setServerError(error);
      });
  };

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <SubmitMoveForm
            initialValues={initialValues}
            onBack={handleBack}
            onSubmit={handleSubmit}
            certificationText={completeCertificationText}
            error={serverError}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

Agreement.propTypes = {
  moveId: PropTypes.string.isRequired,
  setFlashMessage: PropTypes.func.isRequired,
  updateMove: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => ({
  moveId: selectCurrentMove(state)?.id,
});

const mapDispatchToProps = {
  updateMove: updateMoveAction,
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(Agreement);
