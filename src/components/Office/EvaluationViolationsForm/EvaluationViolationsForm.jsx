import React from 'react';
import * as PropTypes from 'prop-types';
import { Grid, GridContainer, Button, FormGroup, Radio, Fieldset, Textarea, Label } from '@trussworks/react-uswds';
import { useParams, useHistory } from 'react-router';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import classnames from 'classnames';
import { useMutation, queryCache } from 'react-query';

import styles from './EvaluationViolationsForm.module.scss';
import SelectedViolation from './SelectedViolation/SelectedViolation';

import { EVALUATION_REPORT } from 'constants/queryKeys';
import ViolationsAccordion from 'components/Office/ViolationsAccordion/ViolationsAccordion';
import { saveEvaluationReport, associateReportViolations } from 'services/ghcApi';
import { DatePickerInput } from 'components/form/fields';
import { MILMOVE_LOG_LEVEL, milmoveLog } from 'utils/milmoveLog';
import { EvaluationReportShape, ReportViolationShape, PWSViolationShape } from 'types';
import { formatDateForSwagger } from 'shared/dates';

const EvaluationViolationsForm = ({ violations, evaluationReport, reportViolations }) => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

  const [mutateEvaluationReport] = useMutation(saveEvaluationReport, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
    onSuccess: () => {
      queryCache.refetchQueries([EVALUATION_REPORT, reportId]).then();
    },
  });

  const [mutateReportViolations] = useMutation(associateReportViolations, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

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

  const saveDraft = async (values) => {
    // pull out fields we dont want to save/update
    const {
      createdAt,
      updatedAt,
      shipmentID,
      id,
      moveID,
      moveReferenceID,
      type,
      officeUser,
      reportID,
      eTag,
      ...existingReportFields
    } = evaluationReport;

    let seriousIncident;
    if (values.seriousIncident) {
      seriousIncident = values.seriousIncident === 'yes';
    }

    const body = {
      ...existingReportFields,
      seriousIncident,
      seriousIncidentDesc: values.seriousIncidentDesc,
      observedClaimsResponseDate: formatDateForSwagger(values.observedClaimsResponseDate),
      observedPickupDate: formatDateForSwagger(values.observedPickupDate),
      observedPickupSpreadStartDate: formatDateForSwagger(values.observedPickupSpreadStartDate),
      observedPickupSpreadEndDate: formatDateForSwagger(values.observedPickupSpreadEndDate),
    };

    await mutateEvaluationReport({ reportID: reportId, ifMatchETag: eTag, body });

    // Also need to update any violations that were selected
    await mutateReportViolations({ reportID: reportId, body: { violations: values.selectedViolations } });
  };

  const handleSaveDraft = async (values) => {
    await saveDraft(values);

    history.push(`/moves/${moveCode}/evaluation-reports`, { showSaveDraftSuccess: true });
  };

  const getInitialValues = () => {
    const selectedViolations = reportViolations ? reportViolations.map((violation) => violation.violationID) : [];

    let seriousIncident;
    if (evaluationReport && Object.hasOwn(evaluationReport, 'seriousIncident')) {
      seriousIncident = evaluationReport.seriousIncident ? 'yes' : 'no';
    }

    const initialValues = {
      selectedViolations,
      seriousIncident,
      seriousIncidentDesc: evaluationReport.seriousIncidentDesc,
      observedClaimsResponseDate: evaluationReport.observedClaimsResponseDate,
      observedPickupDate: evaluationReport.observedPickupDate,
      observedPickupSpreadStartDate: evaluationReport.observedPickupSpreadStartDate,
      observedPickupSpreadEndDate: evaluationReport.observedPickupSpreadEndDate,
    };

    return initialValues;
  };

  return (
    <Formik
      initialValues={getInitialValues()}
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
                      key={`${violationId}-selected`}
                    />
                  ))}
                </Grid>
              </Grid>

              <Grid row>
                <Grid col className={styles.claimDatePicker}>
                  <div>
                    {kpiDates.includes('observedClaimDate') && (
                      <DatePickerInput
                        className={styles.datePicker}
                        label="Observed claims response date"
                        name="observedClaimsResponseDate"
                        hint="Only enter a date here if the claim has a response."
                        showOptional
                      />
                    )}
                    {kpiDates.includes('observedPickupDate') && (
                      <DatePickerInput
                        label="Observed pickup date"
                        name="observedPickupDate"
                        hint="Enter the date you witnessed the pickup."
                      />
                    )}
                    {kpiDates.includes('observedPickupSpreadDates') && (
                      <DatePickerInput label="Observed pickup spread start date" name="observedPickupSpreadStartDate" />
                    )}
                    {kpiDates.includes('observedPickupSpreadDates') && (
                      <DatePickerInput label="Observed pickup spread end date" name="observedPickupSpreadEndDate" />
                    )}
                  </div>
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
                              <Label className={styles.incidentTextAreaLabel} htmlFor="seriousIncidentDesc">
                                Serious incident description
                              </Label>
                              <Field as={Textarea} name="seriousIncidentDesc" id="seriousIncidentDesc" />
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
                    <Button
                      data-testid="saveDraft"
                      type="button"
                      className="usa-button--secondary"
                      onClick={() => handleSaveDraft(values)}
                    >
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

EvaluationViolationsForm.propTypes = {
  violations: PropTypes.arrayOf(PWSViolationShape).isRequired,
  evaluationReport: EvaluationReportShape.isRequired,
  reportViolations: PropTypes.arrayOf(ReportViolationShape).isRequired,
};
