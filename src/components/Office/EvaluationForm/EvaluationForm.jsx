import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import 'styles/office.scss';
import { GridContainer, Grid, Button, Radio, FormGroup, Fieldset, Label, Textarea } from '@trussworks/react-uswds';
import { useParams, useHistory, useLocation } from 'react-router';
import { useMutation, queryCache } from 'react-query';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './EvaluationForm.module.scss';

import { EVALUATION_REPORT } from 'constants/queryKeys';
import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import ConnectedEvaluationReportConfirmationModal from 'components/ConfirmationModals/EvaluationReportConfirmationModal';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { deleteEvaluationReport, saveEvaluationReport, submitEvaluationReport } from 'services/ghcApi';
import { DatePickerInput, DropdownInput } from 'components/form/fields';
import { MILMOVE_LOG_LEVEL, milmoveLog } from 'utils/milmoveLog';
import { formatDateForSwagger } from 'shared/dates';
import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';
import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';

const EvaluationForm = ({
  evaluationReport,
  reportViolations,
  mtoShipments,
  customerInfo,
  grade,
  destinationDutyLocationPostalCode,
}) => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();
  const location = useLocation();

  const [deleteEvaluationReportMutation] = useMutation(deleteEvaluationReport);
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

  const [mutateEvaluationReport] = useMutation(saveEvaluationReport, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
    onSuccess: () => {
      queryCache.refetchQueries([EVALUATION_REPORT, reportId]).then();
    },
  });

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isSubmitModalOpen, setIsSubmitModalOpen] = useState(false);

  // whether or not the delete report modal is displaying
  const toggleDeleteReportModal = () => {
    setIsDeleteModalOpen(!isDeleteModalOpen);
  };

  // cancel report updates but don't delete, just re-route back to reports page
  const cancelForUpdatedReport = () => {
    history.push(`/moves/${moveCode}/evaluation-reports`);
  };

  const deleteReport = async () => {
    // Close the modal
    setIsDeleteModalOpen(!isDeleteModalOpen);

    // Mark as deleted in database
    await deleteEvaluationReportMutation(reportId);

    // Reroute back to eval report page, include flag to know to show alert
    history.push(`/moves/${moveCode}/evaluation-reports`, { showCanceledSuccess: true });
  };

  // passed to the confrimation modal
  const submitReport = async () => {
    // close the modal
    setIsSubmitModalOpen(!isSubmitModalOpen);

    // mark as submitted in the DB
    await submitEvaluationReportMutation({ reportID: reportId, ifMatchETag: evaluationReport.eTag });
  };

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

    // format the inspection type if its there
    const { evaluationType } = values;
    let inspectionType;
    if (evaluationType) {
      if (evaluationType === 'dataReview') {
        inspectionType = 'DATA_REVIEW';
      } else {
        inspectionType = values.evaluationType.toUpperCase();
      }
    }
    // format the location if it's there
    let locationDescription;
    let { evaluationLocation } = values;
    if (evaluationLocation) {
      evaluationLocation = values.evaluationLocation.toUpperCase();
      if (evaluationLocation === 'OTHER') {
        locationDescription = values.otherEvaluationLocation;
      }
    }

    if (values.evaluationLocation === 'other') {
      locationDescription = values.otherEvaluationLocation;
    }

    let timeDepart;
    let evalStart;
    let evalEnd;
    if (inspectionType === 'PHYSICAL' && values.evaluationLocation !== 'other') {
      timeDepart = `${values.timeDepartHour}:${values.timeDepartMinute}`;
      evalStart = `${values.evalStartHour}:${values.evalStartMinute}`;
      evalEnd = `${values.evalEndHour}:${values.evalEndMinute}`;
    }

    let violations;
    if (values.violationsObserved) {
      violations = values.violationsObserved === 'yes';
    }

    let observedDate;
    if (values.observedDate !== 'Invalid date') {
      observedDate = formatDateForSwagger(values.observedDate);
    }

    const body = {
      ...existingReportFields,
      location: evaluationLocation,
      locationDescription,
      inspectionType,
      remarks: values.remarks,
      // this is a yes or no boolean and not a list of the violations
      violationsObserved: violations,
      inspectionDate: formatDateForSwagger(values.inspectionDate),
      timeDepart,
      evalStart,
      evalEnd,
      observedDate,
    };

    await mutateEvaluationReport({ reportID: reportId, ifMatchETag: eTag, body });
  };

  const handleSaveDraft = async (values) => {
    await saveDraft(values);

    history.push(`/moves/${moveCode}/evaluation-reports`, { showSaveDraftSuccess: true });
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

  const handleSelectViolations = async (values) => {
    await saveDraft(values);

    // Reroute to currentURL/violations
    history.push(`${location.pathname}/violations`);
  };

  const isShipment = evaluationReport.type === EVALUATION_REPORT_TYPE.SHIPMENT;

  const modalTitle = (
    <div className={styles.title}>
      <h3>{`Preview and submit ${evaluationReport.type.toLowerCase()} report`}</h3>
      <p>Is all the information shown correct?</p>
    </div>
  );

  const initialValues = {
    remarks: evaluationReport.remarks,
    inspectionDate: evaluationReport.inspectionDate,
    observedDate: evaluationReport.observedDate,
  };

  if (evaluationReport.location) {
    initialValues.evaluationLocation = evaluationReport.location.toLowerCase();
  }

  if (evaluationReport.locationDescription) {
    initialValues.otherEvaluationLocation = evaluationReport.locationDescription;
  }
  if (evaluationReport.inspectionType) {
    if (evaluationReport.inspectionType === 'DATA_REVIEW') {
      initialValues.evaluationType = 'dataReview';
    } else {
      initialValues.evaluationType = evaluationReport.inspectionType.toLowerCase();
    }
  }

  if (evaluationReport.timeDepart) {
    const timeDepartValues = evaluationReport.timeDepart.split(':');
    [initialValues.timeDepartHour, initialValues.timeDepartMinute] = timeDepartValues;
  }

  if (evaluationReport.evalStart) {
    const timeDepartValues = evaluationReport.evalStart.split(':');
    [initialValues.evalStartHour, initialValues.evalStartMinute] = timeDepartValues;
  }

  if (evaluationReport.evalEnd) {
    const timeDepartValues = evaluationReport.evalEnd.split(':');
    [initialValues.evalEndHour, initialValues.evalEndMinute] = timeDepartValues;
  }

  if (evaluationReport.violationsObserved !== undefined) {
    initialValues.violationsObserved = evaluationReport.violationsObserved ? 'yes' : 'no';
  }

  const EvalDurationSchema = Yup.string().when(['evaluationType', 'evaluationLocation'], {
    is: (evaluationType, evaluationLocation) => evaluationType === 'physical' && evaluationLocation !== 'other',
    then: Yup.string().required(),
  });

  const validationSchema = Yup.object().shape(
    {
      inspectionDate: Yup.date().required(),
      evaluationType: Yup.string().required(),
      timeDepartHour: EvalDurationSchema,
      timeDepartMinute: EvalDurationSchema,
      evalStartHour: EvalDurationSchema,
      evalStartMinute: EvalDurationSchema,
      evalEndMinute: EvalDurationSchema,
      evalEndHour: EvalDurationSchema,
      evaluationLocation: Yup.string().required(),
      violationsObserved: Yup.string().required(),
      remarks: Yup.string().required(),
      otherEvaluationLocation: Yup.string().when('evaluationLocation', {
        is: 'other',
        then: Yup.string().required(),
      }),
    },
    [
      ['timeDepartMinute', 'evaluationType'],
      ['timeDepartMinute', 'evaluationLocation'],
      ['timeDepartHour', 'evaluationType'],
      ['timeDepartHour', 'evaluationLocation'],
      ['evalStartMinute', 'evaluationType'],
      ['evalStartMinute', 'evaluationLocation'],
      ['evalEndMinute', 'evaluationType'],
      ['evalEndMinute', 'evaluationLocation'],
      ['evalEndHour', 'evaluationType'],
      ['evalEndHour', 'evaluationLocation'],
    ],
  );

  const minutes = [];

  for (let i = 0; i < 10; i += 1) {
    minutes[i] = { key: `0${String(i)}`, value: `0${String(i)}` };
  }

  for (let i = 10; i < 60; i += 1) {
    minutes[i] = { key: String(i), value: String(i) };
  }

  const hours = [];

  for (let i = 0; i < 10; i += 1) {
    hours[i] = { key: `0${String(i)}`, value: `0${String(i)}` };
  }

  for (let i = 10; i < 25; i += 1) {
    hours[i] = { key: String(i), value: String(i) };
  }

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

  return (
    <>
      <ConnectedDeleteEvaluationReportConfirmationModal
        isOpen={isDeleteModalOpen}
        closeModal={toggleDeleteReportModal}
        submitModal={deleteReport}
      />
      <ConnectedEvaluationReportConfirmationModal
        modalTopRightClose={setIsSubmitModalOpen}
        isOpen={isSubmitModalOpen}
        modalTitle={modalTitle}
        evaluationReport={evaluationReport}
        reportViolations={reportViolations}
        moveCode={moveCode}
        customerInfo={customerInfo}
        grade={grade}
        mtoShipments={mtoShipments}
        modalActions={submitModalActions}
        bordered
        destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
      />

      <Formik
        initialValues={initialValues}
        enableReinitialize
        onSubmit={handlePreviewReport}
        validationSchema={validationSchema}
        validateOnMount
      >
        {({ values, setFieldValue, isValid }) => {
          const showObservedDeliveryDate =
            values.evaluationType === 'physical' && values.evaluationLocation === 'destination' && isShipment;
          const showObservedPickupDate =
            values.evaluationType === 'physical' && values.evaluationLocation === 'origin' && isShipment;

          const showTimeDepartStartEnd = values.evaluationType === 'physical' && values.evaluationLocation !== 'other';

          return (
            <Form className={classnames(formStyles.form, styles.form)} data-testid="evaluationReportForm">
              <GridContainer className={styles.cardContainer}>
                <Grid row className={styles.evalInfoSection}>
                  <Grid col>
                    <h2>Evaluation form</h2>
                    <h3>Evaluation information</h3>
                    <DatePickerInput label="Date of inspection" name="inspectionDate" disableErrorLabel />
                    <FormGroup>
                      <Fieldset className={styles.radioGroup}>
                        <legend className="usa-label">Evaluation type</legend>
                        <Field
                          as={Radio}
                          label="Data review"
                          id="dataReview"
                          name="evaluationType"
                          value="dataReview"
                          title="Data review"
                          type="radio"
                          checked={values.evaluationType === 'dataReview'}
                        />
                        <Field
                          as={Radio}
                          label="Virtual"
                          id="virtual"
                          name="evaluationType"
                          value="virtual"
                          title="Virtual"
                          type="radio"
                          checked={values.evaluationType === 'virtual'}
                        />
                        <Field
                          as={Radio}
                          label="Physical"
                          id="physical"
                          name="evaluationType"
                          value="physical"
                          title="Physical"
                          type="radio"
                          checked={values.evaluationType === 'physical'}
                        />
                      </Fieldset>
                    </FormGroup>
                    {showTimeDepartStartEnd && (
                      <>
                        <legend className="usa-label">Time departed for evaluation</legend>
                        <div className={styles.durationPickers}>
                          <div>
                            <DropdownInput
                              id="timeDepartHour"
                              name="timeDepartHour"
                              label="Hours"
                              className={styles.hourPicker}
                              onChange={(e) => {
                                setFieldValue('timeDepartHour', e.target.value);
                              }}
                              disableErrorLabel
                              options={hours}
                            />
                          </div>
                          <div>
                            <DropdownInput
                              id="timeDepartMinute"
                              name="timeDepartMinute"
                              label="Minutes"
                              className={styles.minutePicker}
                              onChange={(e) => {
                                setFieldValue('timeDepartMinute', e.target.value);
                              }}
                              disableErrorLabel
                              options={minutes}
                            />
                          </div>
                        </div>
                        <legend className="usa-label">Time evaluation started</legend>
                        <div className={styles.durationPickers}>
                          <div>
                            <DropdownInput
                              id="evalStartHour"
                              name="evalStartHour"
                              label="Hours"
                              className={styles.hourPicker}
                              onChange={(e) => {
                                setFieldValue('evalStartHour', e.target.value);
                              }}
                              disableErrorLabel
                              options={hours}
                            />
                          </div>
                          <div>
                            <DropdownInput
                              id="evalStartMinute"
                              name="evalStartMinute"
                              label="Minutes"
                              className={styles.minutePicker}
                              onChange={(e) => {
                                setFieldValue('evalStartMinute', e.target.value);
                              }}
                              disableErrorLabel
                              options={minutes}
                            />
                          </div>
                        </div>
                        <legend className="usa-label">Time evaluation ended</legend>
                        <div className={styles.durationPickers}>
                          <div>
                            <DropdownInput
                              id="evalEndHour"
                              name="evalEndHour"
                              label="Hours"
                              className={styles.hourPicker}
                              onChange={(e) => {
                                setFieldValue('evalEndHour', e.target.value);
                              }}
                              disableErrorLabel
                              options={hours}
                            />
                          </div>
                          <div>
                            <DropdownInput
                              id="evalEndMinute"
                              name="evalEndMinute"
                              label="Minutes"
                              className={styles.minutePicker}
                              onChange={(e) => {
                                setFieldValue('evalEndMinute', e.target.value);
                              }}
                              disableErrorLabel
                              options={minutes}
                            />
                          </div>
                        </div>
                      </>
                    )}
                    <FormGroup>
                      <Fieldset className={styles.radioGroup}>
                        <legend className="usa-label">Evaluation location</legend>
                        <Field
                          as={Radio}
                          label="Origin"
                          id="origin"
                          name="evaluationLocation"
                          value="origin"
                          title="Origin"
                          type="radio"
                          checked={values.evaluationLocation === 'origin'}
                        />
                        {isShipment && (
                          <Field
                            as={Radio}
                            label="Destination"
                            id="destination"
                            name="evaluationLocation"
                            value="destination"
                            title="Destination"
                            type="radio"
                            checked={values.evaluationLocation === 'destination'}
                          />
                        )}
                        <Field
                          as={Radio}
                          label="Other"
                          id="other"
                          name="evaluationLocation"
                          value="other"
                          title="Other"
                          type="radio"
                          checked={values.evaluationLocation === 'other'}
                        />
                        {values.evaluationLocation === 'other' && (
                          <Field
                            as={Textarea}
                            name="otherEvaluationLocation"
                            id="otherEvaluationLocation"
                            className={styles.textArea}
                          />
                        )}
                      </Fieldset>
                    </FormGroup>
                    {showObservedDeliveryDate && (
                      <div className={styles.showOptional}>
                        <DatePickerInput
                          label="Observed delivery date"
                          name="observedDate"
                          hint="Only enter a date here if the delivery you witnessed did not happen on the scheduled delivery date"
                          showOptional
                        />
                      </div>
                    )}
                    {showObservedPickupDate && (
                      <div className={styles.showOptional}>
                        <DatePickerInput
                          label="Observed pickup date"
                          name="observedDate"
                          hint="Only enter a date here if the pickup you witnessed did not happen on the scheduled pickup date"
                          showOptional
                        />
                      </div>
                    )}
                  </Grid>
                </Grid>
                <Grid row className={styles.evalInfoSection}>
                  <Grid col>
                    <h3>Violations</h3>
                    <FormGroup className={styles.violationsGroup}>
                      <Fieldset>
                        <legend className="usa-label">Violations observed</legend>
                        <Field
                          as={Radio}
                          label="No"
                          id="noViolations"
                          name="violationsObserved"
                          value="no"
                          title="No"
                          type="radio"
                          checked={values.violationsObserved === 'no'}
                          data-testid="noViolationsRadioOption"
                          className={styles.radioGroup}
                        />
                        <Field
                          as={Radio}
                          label="Yes"
                          id="yesViolations"
                          name="violationsObserved"
                          value="yes"
                          title="Yes"
                          type="radio"
                          checked={values.violationsObserved === 'yes'}
                          data-testid="yesViolationsRadioOption"
                          className={styles.radioGroup}
                        />
                        {values.violationsObserved === 'yes' && (
                          <p className={styles.violationsInfo}>
                            <small>You will select the specific PWS paragraphs violated on the next screen.</small>
                          </p>
                        )}
                      </Fieldset>
                    </FormGroup>
                  </Grid>
                </Grid>
                <Grid row>
                  <Grid col>
                    <h3>QAE remarks</h3>
                    <Label htmlFor="evaluationRemarks">Evaluation remarks</Label>
                    <Field
                      as={Textarea}
                      name="remarks"
                      id="evaluationRemarks"
                      title="Evaluation remarks"
                      className={styles.textArea}
                    />
                  </Grid>
                </Grid>
              </GridContainer>
              <GridContainer className={styles.buttonContainer}>
                <Grid row>
                  <Grid col>
                    <div className={styles.buttonRow}>
                      {evaluationReport.updatedAt === evaluationReport.createdAt && (
                        <Button
                          className="usa-button--unstyled"
                          onClick={toggleDeleteReportModal}
                          type="button"
                          data-testid="cancelReport"
                        >
                          Cancel
                        </Button>
                      )}
                      {!(evaluationReport.updatedAt === evaluationReport.createdAt) && (
                        <Button
                          className="usa-button--unstyled"
                          data-testid="cancelReport"
                          onClick={cancelForUpdatedReport}
                          type="button"
                        >
                          Cancel
                        </Button>
                      )}
                      <Button type="button" className="usa-button--secondary" onClick={() => handleSaveDraft(values)}>
                        Save draft
                      </Button>
                      {values.violationsObserved === 'yes' ? (
                        <Button
                          disabled={!isValid}
                          onClick={() => handleSelectViolations(values)}
                          type="button"
                          data-testid="selectViolations"
                        >
                          Next: select violations
                        </Button>
                      ) : (
                        <Button
                          disabled={!isValid}
                          type="button"
                          data-testid="reviewAndSubmit"
                          onClick={() => handlePreviewReport(values)}
                        >
                          Review and submit
                        </Button>
                      )}
                    </div>
                  </Grid>
                </Grid>
              </GridContainer>
            </Form>
          );
        }}
      </Formik>
    </>
  );
};

EvaluationForm.propTypes = {
  evaluationReport: EvaluationReportShape.isRequired,
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  destinationDutyLocationPostalCode: PropTypes.string.isRequired,
};

EvaluationForm.defaultProps = {
  mtoShipments: null,
};

export default EvaluationForm;
