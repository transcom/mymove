import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Alert, GridContainer, Grid } from '@trussworks/react-uswds';
import moment from 'moment';
import { connect } from 'react-redux';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import { customerRoutes } from 'constants/routes';
import SubmitMoveForm from 'components/Customer/SubmitMoveForm/SubmitMoveForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { SIGNED_CERT_OPTIONS, MOVE_LOCKED_WARNING, checkIfMoveIsLocked } from 'shared/constants';
import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import { submitMoveForApproval } from 'services/internalApi';
import { selectCurrentMove, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { formatServiceMemberNameToString, formatSwaggerDate } from 'utils/formatters';

export const Agreement = ({ updateMove, setFlashMessage, serviceMember, move }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);
  const [isMoveLocked, setIsMoveLocked] = useState(false);
  const { moveId } = useParams();
  const initialValues = {
    signature: '',
    date: formatSwaggerDate(new Date()),
  };

  const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

  const handleBack = () => navigate(reviewPath);

  const getServiceMemberName = (loggedInUser) => {
    if (loggedInUser) {
      return formatServiceMemberNameToString(serviceMember);
    }
    return '';
  };

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
        navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
      })
      .catch((error) => {
        // TODO - log error internally?
        setServerError(error);
      });
  };

  useEffect(() => {
    if (checkIfMoveIsLocked(move)) {
      setIsMoveLocked(true);
    }
  }, [move]);

  return (
    <>
      {isMoveLocked && (
        <Alert headingLevel="h4" type="warning">
          {MOVE_LOCKED_WARNING}
        </Alert>
      )}
      <GridContainer>
        <NotificationScrollToTop dependency={serverError} />
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <SubmitMoveForm
              initialValues={initialValues}
              onBack={handleBack}
              onSubmit={handleSubmit}
              certificationText={completeCertificationText}
              currentUser={getServiceMemberName(serviceMember)}
              error={serverError}
              isMoveLocked={isMoveLocked}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </>
  );
};

Agreement.propTypes = {
  setFlashMessage: PropTypes.func.isRequired,
  updateMove: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => {
  const move = selectCurrentMove(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    move,
    serviceMember,
  };
};

const mapDispatchToProps = {
  updateMove: updateMoveAction,
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(Agreement);
