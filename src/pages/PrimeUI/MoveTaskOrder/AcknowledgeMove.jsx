import React, { useState } from 'react';
import { Formik } from 'formik';
import { func } from 'prop-types';
import { connect } from 'react-redux';
import { useMutation } from '@tanstack/react-query';
import { Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams, useNavigate, generatePath } from 'react-router-dom';

import { acknowledgeMovesAndShipments } from 'services/primeApi';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { DatePickerInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { formatDateForSwagger, formatDateWithUTC } from 'shared/dates';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import scrollToTop from 'shared/scrollToTop';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import Hint from 'components/Hint/index';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const AcknowledgeMove = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const navigate = useNavigate();
  const { moveCodeOrID } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  const { mutate: acknowledgeMoveRequestMutation } = useMutation(acknowledgeMovesAndShipments, {
    onSuccess: () => {
      setFlashMessage(`ACKNOWLEDGE_MOVE_SUCCESS${moveCodeOrID}`, 'success', 'Successfully acknowledged move', '', true);

      navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        let additionalDetails = '';
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            additionalDetails += `:\n${key} - ${body.invalidFields[key]}`;
          });
        }

        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}${additionalDetails}`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred.',
        });
      }
      scrollToTop();
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const initialValues = {
    primeAcknowledgedAt: formatDateWithUTC(moveTaskOrder.primeAcknowledgedAt, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
  };
  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const onSubmit = (values) => {
    const { primeAcknowledgedAt } = values;

    const body = [
      {
        id: moveTaskOrder.id,
        primeAcknowledgedAt: formatDateForSwagger(primeAcknowledgedAt),
      },
    ];

    acknowledgeMoveRequestMutation({ body });
  };

  return (
    <div className={classnames('grid-container-desktop-lg', 'usa-prose', primeStyles.primeContainer)}>
      <div className="grid-row">
        <div className="grid-col-12">
          {errorMessage?.detail && (
            <div className={primeStyles.errorContainer}>
              <Alert headingLevel="h4" slim type="error">
                <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
              </Alert>
            </div>
          )}
          <h2>Acknowledge Move</h2>
          <SectionWrapper className={formStyles.formSection}>
            <Formik initialValues={initialValues} onSubmit={onSubmit}>
              {({ isValid, isSubmitting, handleSubmit, dirty }) => {
                return (
                  <Form className={formStyles.form}>
                    <dl className={descriptionListStyles.descriptionList} data-testid="moveDetails">
                      <h2>Move</h2>
                      <div className={descriptionListStyles.row}>
                        <dt>Move Code:</dt>
                        <dd>{moveTaskOrder.moveCode}</dd>
                      </div>
                      <div className={descriptionListStyles.row}>
                        <dt>Move Id:</dt>
                        <dd>{moveTaskOrder.id}</dd>
                      </div>
                      {requiredAsteriskMessage}
                      <DatePickerInput
                        data-testid="primeAcknowledgedAt"
                        name="primeAcknowledgedAt"
                        label="Prime Acknowledged At"
                        showRequiredAsterisk
                        required
                        disabled={moveTaskOrder.primeAcknowledgedAt}
                      />
                      <Hint id="primeAcknowledgedAtHint" data-testid="primeAcknowledgedAtHint">
                        Prime Acknowledged At date can only be saved one time.
                      </Hint>
                    </dl>
                    <div className={formStyles.formActions}>
                      <WizardNavigation
                        editMode
                        disableNext={!isValid || !dirty || isSubmitting}
                        onCancelClick={handleClose}
                        onNextClick={handleSubmit}
                      />
                    </div>
                  </Form>
                );
              }}
            </Formik>
          </SectionWrapper>
        </div>
      </div>
    </div>
  );
};

AcknowledgeMove.propTypes = {
  setFlashMessage: func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(AcknowledgeMove);
