import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import { Grid, GridContainer, Button, FormGroup, Radio, Fieldset, Textarea } from '@trussworks/react-uswds';
import { useParams, useHistory } from 'react-router';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import classnames from 'classnames';
import { useMutation, queryCache } from 'react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './EvaluationViolationsForm.module.scss';
import SelectedViolation from './SelectedViolation/SelectedViolation';

import { EVALUATION_REPORT } from 'constants/queryKeys';
import ViolationsAccordion from 'components/Office/ViolationsAccordion/ViolationsAccordion';
import ConnectedEvaluationReportConfirmationModal from 'components/ConfirmationModals/EvaluationReportConfirmationModal';
import { saveEvaluationReport, associateReportViolations, submitEvaluationReport } from 'services/ghcApi';
import { DatePickerInput } from 'components/form/fields';
import { MILMOVE_LOG_LEVEL, milmoveLog } from 'utils/milmoveLog';
import { EvaluationReportShape, ReportViolationShape, PWSViolationShape } from 'types';

const EvaluationViolationsForm = ({ violations, evaluationReport, reportViolations, customerInfo, mtoShipments }) => {
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

  const [submitEvaluationReportMutation] = useMutation(submitEvaluationReport, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
    onSuccess: () => {
      // Reroute back to eval report page, include flag to show success alert
      history.push(`/moves/${moveCode}/evaluation-reports`, { showSubmitSuccess: true });
    },
  });

  const [mutateReportViolations] = useMutation(associateReportViolations, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const [isSubmitModalOpen, setIsSubmitModalOpen] = useState(false);

  // passed to the confrimation modal
  const submitReport = async () => {
    // close the modal
    setIsSubmitModalOpen(!isSubmitModalOpen);

    // mark as submitted in the DB
    await submitEvaluationReportMutation({ reportID: reportId, ifMatchETag: evaluationReport.eTag });
  };

  const modalTitle = (
    <div className={styles.title}>
      <h3>{`Preview and submit ${evaluationReport.type.toLowerCase()} report`}</h3>
      <p>Is all the information shown correct?</p>
    </div>
  );

  const submitModalActions = (
    <div className={styles.modalActions}>
      <Button
        type="button"
        onClick={() => setIsSubmitModalOpen()}
        aria-label="Back to Evaluation form"
        unstyled
        className={styles.closeModalBtn}
      >
        <FontAwesomeIcon icon="chevron-left" className={styles.backIcon} /> Back to Evaluation form
      </Button>
      <Button
        type="submit"
        onClick={submitReport}
        data-testid="modalSubmitButton"
        aria-label="Submit"
        className={styles.submitModalBtn}
      >
        Submit
      </Button>
    </div>
  );

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
    const { createdAt, updatedAt, shipmentID, id, moveID, moveReferenceID, ...existingReportFields } = evaluationReport;
    const body = {
      ...existingReportFields,
      // TODO: Add serious incident and date fields that are on the form
    };
    const { eTag } = evaluationReport;
    await mutateEvaluationReport({ reportID: reportId, ifMatchETag: eTag, body });

    // Also need to update any violations that were selected
    await mutateReportViolations({ reportID: reportId, body: { violations: values.selectedViolations } });
  };

  const handleSaveDraft = async (values) => {
    await saveDraft(values);

    history.push(`/moves/${moveCode}/evaluation-reports`, { showSaveDraftSuccess: true });
  };

  const initialValues = {
    selectedViolations: reportViolations ? reportViolations.map((violation) => violation.violationID) : [],
  };

  // Review and Submit button
  // Saves report changes
  // displays report preview ahead of final submission
  const handlePreviewReport = async (values) => {
    // save updates
    await saveDraft(values);

    // open the modal to submit
    setIsSubmitModalOpen(!isSubmitModalOpen);
  };

  return (
    <>
      <ConnectedEvaluationReportConfirmationModal
        modalTopRightClose={setIsSubmitModalOpen}
        isOpen={isSubmitModalOpen}
        modalTitle={modalTitle}
        evaluationReport={evaluationReport}
        moveCode={moveCode}
        customerInfo={customerInfo}
        grade={customerInfo.grade}
        mtoShipments={mtoShipments}
        modalActions={submitModalActions}
        reportViolations={reportViolations}
        bordered
      />
      <Formik
        initialValues={initialValues}
        enableReinitialize
        onSubmit={handlePreviewReport}
        validationSchema={validationSchema}
        validateOnMount
      >
        {({ values, setFieldValue, isValid }) => {
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
                          name="observedClaimDate"
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
                        <DatePickerInput
                          label="Observed pickup spread start date"
                          name="observedpickupStartDateScheduling"
                        />
                      )}
                      {kpiDates.includes('observedPickupSpreadDates') && (
                        <DatePickerInput
                          label="Observed pickup spread end date"
                          name="observedpickupEndDateScheduling"
                        />
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
                      <Button
                        data-testid="saveDraft"
                        type="button"
                        className="usa-button--secondary"
                        onClick={() => handleSaveDraft(values)}
                      >
                        Save draft
                      </Button>
                      <Button disabled={!isValid} type="button" onClick={() => handlePreviewReport(values)}>
                        Review and submit
                      </Button>
                    </div>
                  </Grid>
                </Grid>
              </GridContainer>
            </>
          );
        }}
      </Formik>
    </>
  );
};

export default EvaluationViolationsForm;

EvaluationViolationsForm.propTypes = {
  violations: PropTypes.arrayOf(PWSViolationShape).isRequired,
  evaluationReport: EvaluationReportShape.isRequired,
  reportViolations: PropTypes.arrayOf(ReportViolationShape).isRequired,
};
