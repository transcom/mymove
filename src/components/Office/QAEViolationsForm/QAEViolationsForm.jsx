import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import { Grid, GridContainer, Button, FormGroup, Radio, Fieldset, Textarea } from '@trussworks/react-uswds';
import { useParams, useNavigate } from 'react-router-dom';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { useMutation, useQueryClient } from '@tanstack/react-query';

import styles from './QAEViolationsForm.module.scss';
import SelectedViolation from './SelectedViolation/SelectedViolation';

import { EVALUATION_REPORT, REPORT_VIOLATIONS } from 'constants/queryKeys';
import ViolationsAccordion from 'components/Office/ViolationsAccordion/ViolationsAccordion';
import ConnectedEvaluationReportConfirmationModal from 'components/ConfirmationModals/EvaluationReportConfirmationModal';
import { saveEvaluationReport, associateReportViolations, submitEvaluationReport } from 'services/ghcApi';
import { DatePickerInput } from 'components/form/fields';
import { milmoveLogger } from 'utils/milmoveLog';
import { EvaluationReportShape, ReportViolationShape, PWSViolationShape, CustomerShape, ShipmentShape } from 'types';
import { formatDateForSwagger } from 'shared/dates';

const QAEViolationsForm = ({
  violations,
  evaluationReport,
  reportViolations,
  customerInfo,
  grade,
  mtoShipments,
  destinationDutyLocationPostalCode,
}) => {
  const { moveCode, reportId } = useParams();
  const navigate = useNavigate();

  const queryClient = useQueryClient();
  const { mutate: mutateEvaluationReport } = useMutation(saveEvaluationReport, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
    onSuccess: async () => {
      await queryClient.refetchQueries([EVALUATION_REPORT, reportId]);
    },
  });

  const { mutate: submitEvaluationReportMutation } = useMutation(submitEvaluationReport, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const { mutate: mutateReportViolations } = useMutation(associateReportViolations, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
    onSuccess: async () => {
      await queryClient.refetchQueries([REPORT_VIOLATIONS, reportId]);
    },
  });

  const [isSubmitModalOpen, setIsSubmitModalOpen] = useState(false);

  // passed to the confrimation modal
  const submitReport = async () => {
    // close the modal
    setIsSubmitModalOpen(!isSubmitModalOpen);

    // mark as submitted in the DB
    submitEvaluationReportMutation(
      { reportID: reportId, ifMatchETag: evaluationReport.eTag },
      {
        onSuccess: () => {
          // Reroute back to eval report page, include flag to show success alert
          navigate(`/moves/${moveCode}/evaluation-reports`, { state: { showSubmitSuccess: true } });
        },
      },
    );
  };

  const modalTitle = (
    <div className={styles.title}>
      <h3>{`Preview and submit ${
        evaluationReport && evaluationReport.type ? evaluationReport.type.toLowerCase() : ''
      } report`}</h3>
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
        data-testid="backToEvalFromSubmit"
      >
        Cancel
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
    navigate(`/moves/${moveCode}/evaluation-reports/${reportId}`);
  };

  const cancelForViolations = () => {
    navigate(`/moves/${moveCode}/evaluation-reports`);
  };

  // Get distinct categories
  const categories = [...new Set(violations.map((item) => item.category))];

  /*

  Form is valid when:
    At least one validation from the list must be selected

    Serious incident radio button must be selected (Yes or No)
      if Yes -> serious incident description is required

    Additional KPI fields which aren't marked optional must be filled when visible

  */
  const validationSchema = Yup.object().shape({
    selectedViolations: Yup.array().of(Yup.string()).min(1),
    seriousIncident: Yup.string().required(),
    seriousIncidentDesc: Yup.string().when('seriousIncident', {
      is: 'yes',
      then: (schema) => schema.required(),
    }),
    observedClaimsResponseDate: Yup.date().when('kpiViolations', {
      is: (kpiViolations) => kpiViolations.includes('observedClaimsResponseDate'),
      then: (schema) => schema.optional(),
    }),
    observedPickupDate: Yup.date().when('kpiViolations', {
      is: (kpiViolations) => kpiViolations.includes('observedPickupDate'),
      then: (schema) => schema.required(),
    }),
    observedPickupSpreadStartDate: Yup.date().when('kpiViolations', {
      is: (kpiViolations) => kpiViolations.includes('observedPickupSpreadDates'),
      then: (schema) => schema.required(),
    }),
    observedPickupSpreadEndDate: Yup.date().when('kpiViolations', {
      is: (kpiViolations) => kpiViolations.includes('observedPickupSpreadDates'),
      then: (schema) => schema.required(),
    }),
    observedDeliveryDate: Yup.date().when('kpiViolations', {
      is: (kpiViolations) => kpiViolations.includes('observedDeliveryDate'),
      then: (schema) => schema.required(),
    }),
  });

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
      seriousIncidentDesc: seriousIncident === false ? null : values.seriousIncidentDesc,
      observedClaimsResponseDate: formatDateForSwagger(values.observedClaimsResponseDate),
      observedPickupDate: formatDateForSwagger(values.observedPickupDate),
      observedPickupSpreadStartDate: formatDateForSwagger(values.observedPickupSpreadStartDate),
      observedPickupSpreadEndDate: formatDateForSwagger(values.observedPickupSpreadEndDate),
      observedDeliveryDate: formatDateForSwagger(values.observedDeliveryDate),
    };

    mutateEvaluationReport({ reportID: reportId, ifMatchETag: eTag, body });

    // Also need to update any violations that were selected
    mutateReportViolations({ reportID: reportId, body: { violations: values.selectedViolations } });
  };

  const handleSaveDraft = async (values) => {
    saveDraft(values);

    navigate(`/moves/${moveCode}/evaluation-reports`, { state: { showSaveDraftSuccess: true } });
  };

  const kpiViolationList = violations.filter((item) => item.isKpi);

  const getInitialValues = () => {
    const selectedViolations = reportViolations ? reportViolations.map((violation) => violation.violationID) : [];

    let seriousIncident;
    if (evaluationReport && 'seriousIncident' in evaluationReport) {
      seriousIncident = evaluationReport.seriousIncident ? 'yes' : 'no';
    }

    const kpiViolations = [];

    if (reportViolations) {
      reportViolations.forEach((entry) => {
        if (entry.violation?.isKpi) {
          const ade = entry.violation.additionalDataElem;
          kpiViolations.push(ade);
        }
      });
    }

    const initialValues = {
      selectedViolations,
      seriousIncident,
      seriousIncidentDesc: evaluationReport?.seriousIncidentDesc,
      observedPickupSpreadStartDate: evaluationReport?.observedPickupSpreadStartDate,
      observedPickupSpreadEndDate: evaluationReport?.observedPickupSpreadEndDate,
      observedClaimsResponseDate: evaluationReport?.observedClaimsResponseDate,
      observedPickupDate: evaluationReport?.observedPickupDate,
      observedDeliveryDate: evaluationReport?.observedDeliveryDate,
      kpiViolations,
    };

    initialValues.kpiViolations = kpiViolations;

    return initialValues;
  };

  // Review and Submit button
  // Saves report changes
  // displays report preview ahead of final submission
  const handlePreviewReport = (values) => {
    // save updates
    saveDraft(values);

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
        grade={grade}
        mtoShipments={mtoShipments}
        modalActions={submitModalActions}
        reportViolations={reportViolations}
        bordered
        destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
      />
      <Formik
        initialValues={getInitialValues()}
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

            const kpiViolation = kpiViolationList.find((entry) => entry.id === id);
            const prevSelectedKpiViolations = values.kpiViolations || [];

            // removing
            if (prevSelectedViolations.includes(id)) {
              setFieldValue(
                fieldKey,
                prevSelectedViolations.filter((violationId) => violationId !== id),
              );

              if (kpiViolation) {
                // perform a lookup rather than straight filter as we may have more than one of a single KPI violation in the array that are still valid
                const location = prevSelectedKpiViolations.findIndex(
                  (entry) => entry === kpiViolation.additionalDataElem,
                );

                prevSelectedKpiViolations.splice(location, 1);

                setFieldValue('kpiViolations', prevSelectedKpiViolations);
              }
            } else {
              setFieldValue(fieldKey, [...prevSelectedViolations, id]);

              if (kpiViolation) {
                setFieldValue('kpiViolations', [...prevSelectedKpiViolations, kpiViolation.additionalDataElem]);
              }
            }
          };

          return (
            <>
              <GridContainer className={styles.cardContainer} data-testid="evaluationViolationsForm">
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
                      {values.kpiViolations.includes('observedClaimsReponseDate') && (
                        <DatePickerInput
                          className={styles.datePicker}
                          label="Observed claims response date"
                          name="observedClaimsResponseDate"
                          hint="Only enter a date here if the claim has a response."
                          showOptional
                        />
                      )}
                      {values.kpiViolations.includes('observedPickupDate') && (
                        <DatePickerInput
                          label="Observed pickup date"
                          name="observedPickupDate"
                          hint="Enter the date you witnessed the pickup."
                        />
                      )}
                      {values.kpiViolations.includes('observedPickupSpreadDates') && (
                        <DatePickerInput
                          label="Observed pickup spread start date"
                          name="observedPickupSpreadStartDate"
                        />
                      )}
                      {values.kpiViolations.includes('observedPickupSpreadDates') && (
                        <DatePickerInput label="Observed pickup spread end date" name="observedPickupSpreadEndDate" />
                      )}
                      {values.kpiViolations.includes('observedClaimsResponseDate') && (
                        <DatePickerInput
                          className={styles.datePicker}
                          label="Observed claims response date"
                          name="observedClaimsResponseDate"
                          hint="Only enter a date here if the claim has a response."
                          showOptional
                        />
                      )}
                      {values.kpiViolations.includes('observedDeliveryDate') && (
                        <DatePickerInput
                          className={styles.datePicker}
                          label="Observed delivery date"
                          name="observedDeliveryDate"
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
                                <Field as={Textarea} name="seriousIncidentDesc" />
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
                    <div className={styles.buttonRowSplit}>
                      <Button unstyled type="button" onClick={handleBackToEvalForm} data-testid="backToEvalForm">
                        {'< Back to evaluation form'}
                      </Button>
                      <div className={styles.buttonRow}>
                        <Button
                          type="button"
                          onClick={cancelForViolations}
                          data-testid="cancelReport"
                          unstyled
                          className={styles.cancelBtn}
                        >
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
                        <Button
                          disabled={!isValid}
                          type="button"
                          onClick={() => handlePreviewReport(values)}
                          data-testid="reviewAndSubmit"
                        >
                          Review and submit
                        </Button>
                      </div>
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

QAEViolationsForm.propTypes = {
  violations: PropTypes.arrayOf(PWSViolationShape).isRequired,
  evaluationReport: EvaluationReportShape.isRequired,
  reportViolations: PropTypes.arrayOf(ReportViolationShape),
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
  destinationDutyLocationPostalCode: PropTypes.string.isRequired,
};

QAEViolationsForm.defaultProps = {
  mtoShipments: null,
  reportViolations: null,
};

export default QAEViolationsForm;
