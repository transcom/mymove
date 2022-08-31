import React from 'react';
import 'styles/office.scss';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams, useHistory } from 'react-router';
import { Formik } from 'formik';
import * as Yup from 'yup';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationViolationsStyles from './EvaluationViolations.module.scss';

import { useEvaluationReportQueries, usePWSViolationsQueries } from 'hooks/queries';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';
import ViolationsAccordion from 'components/Office/ViolationsAccordion/ViolationsAccordion';

const EvaluationViolations = () => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

  const { evaluationReport } = useEvaluationReportQueries(reportId);
  const { violations } = usePWSViolationsQueries();

  const handleBackToEvalForm = () => {
    // TODO: Save as draft before rerouting
    history.push(`/moves/${moveCode}/evaluation-reports/${reportId}`);
  };

  const cancelForViolations = () => {
    history.push(`/moves/${moveCode}/evaluation-reports`);
  };

  // Get distinct categories
  const categories = [...new Set(violations.map((item) => item.category))];
  const validationSchema = Yup.object().shape({});

  return (
    <div className={classnames(styles.tabContent, evaluationViolationsStyles.tabContent)}>
      <GridContainer>
        <QaeReportHeader report={evaluationReport} />

        <Formik
          initialValues={{}}
          enableReinitialize
          onSubmit={() => {}}
          validationSchema={validationSchema}
          validateOnMount
        >
          {({ values, setFieldValue }) => {
            const handleAccordionChange = (id) => {
              const fieldKey = 'selectedViolations';
              const prevSelectedViolations = values[fieldKey] || [];
              if (prevSelectedViolations.includes(id)) {
                setFieldValue(
                  fieldKey,
                  prevSelectedViolations.filter((violationId) => violationId !== id),
                );
              } else {
                setFieldValue(fieldKey, [...prevSelectedViolations, id]);
              }
            };

            return (
              <>
                <GridContainer className={evaluationViolationsStyles.cardContainer}>
                  <Grid row>
                    <Grid col>
                      <h2>Select violations</h2>
                      <p className={evaluationViolationsStyles.detailText}>
                        Select the paragraph from the Performance Work Statement (PWS) that the GHC Prime has violated.
                      </p>
                    </Grid>
                  </Grid>
                  {categories.map((category) => (
                    <ViolationsAccordion
                      onChange={handleAccordionChange}
                      violations={violations.filter((violation) => violation.category === category)}
                      key={`${category}-category`}
                      selected={values.selectedViolations}
                    />
                  ))}
                  <Grid row>
                    <Grid col>
                      <h3>Violations Selected ({values.selectedViolations && values.selectedViolations.length})</h3>
                      {values.selectedViolations?.map((violationId) => {
                        const violation = violations.find((v) => v.id === violationId);
                        return (
                          <div key={`${violationId}-violation`}>
                            <h5
                              className={styles.checkboxLabel}
                            >{`${violation.paragraphNumber} ${violation.title}`}</h5>
                            <small>{violation.requirementSummary}</small>
                            <Button
                              type="button"
                              unstyled
                              onClick={() => {
                                setFieldValue(
                                  'selectedViolations',
                                  values.selectedViolations.filter((id) => id !== violationId),
                                );
                              }}
                            >
                              Remove
                            </Button>
                          </div>
                        );
                      })}
                    </Grid>
                  </Grid>
                </GridContainer>
                <GridContainer className={evaluationViolationsStyles.buttonContainer}>
                  <Grid row>
                    <Grid col>
                      <div className={evaluationViolationsStyles.buttonRow}>
                        <Button
                          className={classnames(evaluationViolationsStyles.backToEvalButton, 'usa-button--unstyled')}
                          type="button"
                          onClick={handleBackToEvalForm}
                        >
                          {'< Back to Evaluation form'}
                        </Button>
                        <div className={evaluationViolationsStyles.grow} />
                        <Button className="usa-button--unstyled" type="button" onClick={cancelForViolations}>
                          Cancel
                        </Button>
                        <Button data-testid="saveDraft" type="submit" className="usa-button--secondary">
                          Save draft
                        </Button>
                        <Button disabled>Review and submit</Button>
                      </div>
                    </Grid>
                  </Grid>
                </GridContainer>
              </>
            );
          }}
        </Formik>
      </GridContainer>
    </div>
  );
};

export default EvaluationViolations;
