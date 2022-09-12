import React from 'react';
import { Grid, GridContainer, Button, FormGroup, Radio, Fieldset, Textarea } from '@trussworks/react-uswds';
import { useParams, useHistory } from 'react-router';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import classnames from 'classnames';

import styles from './EvaluationViolationsForm.module.scss';
import SelectedViolation from './SelectedViolation/SelectedViolation';

import ViolationsAccordion from 'components/Office/ViolationsAccordion/ViolationsAccordion';

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
        // Handles adding/removing violations form formic `values`
        const toggleSelectedViolation = (id) => {
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
                  onChange={toggleSelectedViolation}
                  violations={violations.filter((violation) => violation.category === category)}
                  key={`${category}-category`}
                  selected={values.selectedViolations}
                />
              ))}

              {/* Selected Violations */}
              <Grid row>
                <Grid col>
                  <>
                    <hr className={styles.divider} />
                    <h3>
                      Violations Selected ({(values.selectedViolations && values.selectedViolations.length) || '0'})
                    </h3>
                  </>
                  {values.selectedViolations?.map((violationId) => (
                    <SelectedViolation
                      violation={violations.find((v) => v.id === violationId)}
                      unselectViolation={toggleSelectedViolation}
                    />
                  ))}
                </Grid>
              </Grid>

              {/* Serious incident */}
              <Grid row>
                <Grid col>
                  <div className={styles.incident}>
                    <hr className={styles.divider} />
                    <h3 className={styles.siHeading}>Serious incident</h3>
                    <FormGroup>
                      <Fieldset>
                        <div className={styles.serious}>
                          <legend data-testid="seriousIncidentLegend" className="usa-label">
                            Serious incident
                          </legend>
                        </div>
                        <div className={styles.seriousIncident}>
                          <Field
                            as={Radio}
                            label="No"
                            id="no"
                            name="seriousIncident"
                            value="no"
                            title="No"
                            type="radio"
                            checked={values.seriousIncident === 'no'}
                          />
                          <Field
                            as={Radio}
                            label="Yes"
                            id="yes"
                            name="seriousIncident"
                            value="yes"
                            title="Yes"
                            type="radio"
                            checked={values.seriousIncident === 'yes'}
                          />
                          {values.seriousIncident === 'yes' && (
                            <>
                              <p className={styles.incidentTextAreaLabel}>Serious incident description</p>
                              <Field as={Textarea} name="yesSeriousIncident" />
                            </>
                          )}
                        </div>
                      </Fieldset>
                    </FormGroup>
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
