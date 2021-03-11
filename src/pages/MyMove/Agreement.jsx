import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import moment from 'moment';
import { connect } from 'react-redux';
import { push as pushAction } from 'connected-react-router';

import SubmitMoveForm from 'components/Customer/SubmitMoveForm/SubmitMoveForm';
import ScrollToTop from 'components/ScrollToTop';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import { getPPMsForMove, submitMoveForApproval } from 'services/internalApi';
import { selectCurrentPPM, selectCurrentMove } from 'store/entities/selectors';
import { updatePPMs as updatePPMsAction, updateMove as updateMoveAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

export const Agreement = ({ moveId, ppmId, updatePPMs, updateMove, push, setFlashMessage }) => {
  const [serverError, setServerError] = useState(null);

  useEffect(() => {
    getPPMsForMove(moveId).then((response) => updatePPMs(response));
  });

  const handleBack = () => push(`/moves/${moveId}/review`);

  const handleSubmit = (values) => {
    const submitDate = moment().format();

    const data = {
      certification_text: completeCertificationText,
      date: submitDate,
      signature: values.signature,
      personally_procured_move_id: ppmId,
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
  updatePPMs: PropTypes.func.isRequired,
  updateMove: PropTypes.func.isRequired,
  ppmId: PropTypes.string,
};

Agreement.defaultProps = {
  ppmId: undefined,
};

const mapStateToProps = (state) => ({
  moveId: selectCurrentMove(state)?.id,
  ppmId: selectCurrentPPM(state)?.id,
});

const mapDispatchToProps = {
  updatePPMs: updatePPMsAction,
  updateMove: updateMoveAction,
  push: pushAction,
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(Agreement);
