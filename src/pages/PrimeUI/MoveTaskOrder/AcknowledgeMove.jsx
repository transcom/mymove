import React, { useState } from 'react';
import { Formik } from 'formik';
import { func } from 'prop-types';
import { connect } from 'react-redux';
import { Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams, useNavigate, generatePath } from 'react-router-dom';

import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { DatePickerInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { formatDateWithUTC } from 'shared/dates';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const AcknowledgeMove = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const navigate = useNavigate();
  const { moveCodeOrID } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const initialValues = {
    primeAcknowledgedAt: formatDateWithUTC(moveTaskOrder.primeAcknowledgedAt, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
  };
  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
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
          <SectionWrapper className={formStyles.formSection}>
            <Formik initialValues={initialValues} validateOnMount>
              {({ isValid, isSubmitting, handleSubmit }) => {
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
                      <DatePickerInput
                        data-testid="primeAcknowledgedAt"
                        name="primeAcknowledgedAt"
                        label="Prime Acknowledged At"
                      />
                    </dl>
                    <div className={formStyles.formActions}>
                      <WizardNavigation
                        editMode
                        disableNext={!isValid || isSubmitting}
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
