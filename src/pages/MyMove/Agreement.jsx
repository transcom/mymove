import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import moment from 'moment';
import { connect } from 'react-redux';
import { generatePath } from 'react-router';
import { push as pushAction } from 'connected-react-router';

import { customerRoutes } from 'constants/routes';
import SubmitMoveForm from 'components/Customer/SubmitMoveForm/SubmitMoveForm';
import ScrollToTop from 'components/ScrollToTop';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import { submitMoveForApproval } from 'services/internalApi';
import { selectCurrentMove } from 'store/entities/selectors';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { formatSwaggerDate } from 'utils/formatters';

export const Agreement = ({ moveId, updateMove, push, setFlashMessage }) => {
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    signature: '',
    date: formatSwaggerDate(new Date()),
  };

  const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

  const handleBack = () => push(reviewPath);

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
        push('/');
      })
      .catch((error) => {
        // TODO - log error internally?
        setServerError(error);
      });
  };

  return (
    <GridContainer>
      <ScrollToTop otherDep={serverError} />
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
  push: PropTypes.func.isRequired,
  updateMove: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => ({
  moveId: selectCurrentMove(state)?.id,
});

const mapDispatchToProps = {
  updateMove: updateMoveAction,
  push: pushAction,
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(Agreement);
