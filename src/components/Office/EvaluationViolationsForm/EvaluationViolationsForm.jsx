import React from 'react';
import { Grid, GridContainer, Button } from '@trussworks/react-uswds';
import { useParams, useHistory } from 'react-router';
import * as Yup from 'yup';
import { Formik } from 'formik';
import classnames from 'classnames';

import styles from './EvaluationViolationsForm.module.scss';

import ViolationsAccordion from 'components/Office/ViolationsAccordion/ViolationsAccordion';
import { DatePickerInput } from 'components/form/fields';

const EvaluationViolationsForm = ({ violations }) => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

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
        let kpiItems = [];
        if (values.selectedViolations) {
          kpiItems = violations.filter((item) => values.selectedViolations.includes(item.id) && item.isKpi);
        }
        const kpiDates = [...new Set(kpiItems.map((item) => item.additionalDataElem))];

        return (
          <>
            <GridContainer className={styles.cardContainer}>
              <Grid row>
                <Grid col>
                  <h2>Select violations</h2>
                  <p className={styles.detailText}>
                    Select the paragraph from the Performance Work Statement (PWS) that the GHC Prime has violated.
                  </p>
                </Grid>
              </Grid>

              {/* Violations Accordions */}
              {categories.map((category) => (
                <ViolationsAccordion
                  onChange={handleAccordionChange}
                  violations={violations.filter((violation) => violation.category === category)}
                  key={`${category}-category`}
                  selected={values.selectedViolations}
                />
              ))}

              {/* Selected Violations */}
              <Grid row>
                <Grid col>
                  <h3>
                    Violations Selected ({(values.selectedViolations && values.selectedViolations.length) || '0'})
                  </h3>
                  {values.selectedViolations?.map((violationId) => {
                    const violation = violations.find((v) => v.id === violationId);
                    return (
                      <div
                        key={`${violationId}-violation`}
                        style={{ display: 'flex', border: '1px solid #A9AEB1', borderRadius: '2px' }}
                      >
                        <div>
                          <h5 className={styles.checkboxLabel}>{`${violation.paragraphNumber} ${violation.title}`}</h5>
                          <small>{violation.requirementSummary}</small>
                        </div>
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

              <Grid row>
                <Grid col>
                  <div>
                    <h5>Optional</h5>
                    {kpiDates.includes('Observed Claim Response Date (Or None)') && (
                      <DatePickerInput
                        label="Observed claim response date (optional)"
                        name="observedClaimDate"
                        hint="Only enter a date here if the claim has a response."
                      />
                    )}
                    {kpiDates.includes('QAE Observed Pickup Date') && (
                      <DatePickerInput
                        label="Observed pickup date"
                        name="observedPickupDate"
                        hint="Enter the date you witnessed the pickup."
                      />
                    )}
                    {kpiDates.includes(
                      'QAE Observed Pickup Spread Start (Date); QAE Observed Pickup Spread End (Date)',
                    ) && (
                      <DatePickerInput
                        label="Observed pickup spread start date"
                        name="observedpickupStartDateScheduling"
                      />
                    )}
                    {kpiDates.includes(
                      'QAE Observed Pickup Spread Start (Date); QAE Observed Pickup Spread End (Date)',
                    ) && (
                      <DatePickerInput label="Observed pickup spread End date" name="observedpickupEndDateScheduling" />
                    )}
                  </div>
                </Grid>
              </Grid>
            </GridContainer>

            {/* Buttons */}
            <GridContainer className={styles.buttonContainer}>
              <Grid row>
                <Grid col>
                  <div className={styles.buttonRow}>
                    <Button
                      className={classnames(styles.backToEvalButton, 'usa-button--unstyled')}
                      type="button"
                      onClick={handleBackToEvalForm}
                    >
                      {'< Back to Evaluation form'}
                    </Button>
                    <div className={styles.grow} />
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
  );
};

export default EvaluationViolationsForm;
