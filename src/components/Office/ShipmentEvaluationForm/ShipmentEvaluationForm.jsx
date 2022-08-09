import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import 'styles/office.scss';
import { GridContainer, Grid, Button, Radio, FormGroup, Fieldset, Label, Textarea } from '@trussworks/react-uswds';
import { useParams, useHistory, useLocation } from 'react-router';
import { useMutation } from 'react-query';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import classnames from 'classnames';

import styles from './ShipmentEvaluationForm.module.scss';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import { deleteEvaluationReport, saveEvaluationReport } from 'services/ghcApi';
import { DatePickerInput, DropdownInput } from 'components/form/fields';
import { MILMOVE_LOG_LEVEL, milmoveLog } from 'utils/milmoveLog';
import { formatDateForSwagger } from 'shared/dates';

const ShipmentEvaluationForm = ({ evaluationReport }) => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();
  const location = useLocation();

  const [isDeleteModelOpen, setIsDeleteModelOpen] = useState(false);

  const [deleteEvaluationReportMutation] = useMutation(deleteEvaluationReport);

  const toggleCancelModel = () => {
    setIsDeleteModelOpen(!isDeleteModelOpen);
  };

  const cancelReport = async () => {
    // Close the modal
    setIsDeleteModelOpen(!isDeleteModelOpen);

    // Mark as deleted in database
    await deleteEvaluationReportMutation(reportId);

    // Reroute back to eval report page, include flag to know to show alert
    history.push(`/moves/${moveCode}/evaluation-reports`, { showDeleteSuccess: true });
  };

  const [mutateEvaluationReport] = useMutation(saveEvaluationReport, {
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const convertToMinutes = (hours, minutes) => {
    return Number(hours || 0) * 60 + Number(minutes || 0);
  };

  const convertToHoursAndMinutes = (totalMinutes) => {
    // divide and round down to get hours
    const hours = Math.floor(totalMinutes / 60);
    // use modulus operator to get the remainder for minutes
    const minutes = totalMinutes % 60;
    return { hours, minutes };
  };

  const saveDraft = async (values) => {
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
    let { evaluationLocation } = values;
    if (evaluationLocation) {
      evaluationLocation = values.evaluationLocation.toUpperCase();
    }
    let evalMinutes;
    // calculate the minutes for evaluation length
    if (values.evalLengthHour || values.evalLengthMinute) {
      // convert hours to minutes and add to minutes
      evalMinutes = convertToMinutes(values.evalLengthHour, values.evalLengthMinute);
    }

    let travelMinutes;
    if (values.minute || values.hour) {
      travelMinutes = convertToMinutes(values.hour, values.minute);
    }

    let violations;
    if (values.violationsObserved) {
      violations = values.violationsObserved === 'yes';
    }

    const body = {
      location: evaluationLocation,
      locationDescription: values.otherEvaluationLocation,
      inspectionType,
      remarks: values.remarks,
      // hard coded until violations work
      violationsObserved: violations,
      inspectionDate: formatDateForSwagger(values.inspectionDate),
      evaluationLengthMinutes: evalMinutes,
      travelTimeMinutes: travelMinutes,
      observedDate: formatDateForSwagger(values.observedDate),
    };
    const { eTag } = evaluationReport;
    await mutateEvaluationReport({ reportID: reportId, ifMatchETag: eTag, body });
  };

  const handleSubmitSaveDraft = async (values) => {
    await saveDraft(values);

    history.push(`/moves/${moveCode}/evaluation-reports`, { showSaveDraftSuccess: true });
  };

  const handleSelectViolations = async (values) => {
    await saveDraft(values);

    // Reroute to currentURL/violations
    history.push(`${location.pathname}/violations`);
  };

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
  if (evaluationReport.evaluationLengthMinutes) {
    const { hours, minutes } = convertToHoursAndMinutes(evaluationReport.evaluationLengthMinutes);
    initialValues.evalLengthMinute = minutes;
    initialValues.evalLengthHour = hours;
  }

  if (evaluationReport.travelTimeMinutes) {
    const { hours, minutes } = convertToHoursAndMinutes(evaluationReport.travelTimeMinutes);
    initialValues.minute = minutes;
    initialValues.hour = hours;
  }

  if (evaluationReport.violationsObserved !== undefined) {
    initialValues.violationsObserved = evaluationReport.violationsObserved ? 'yes' : 'no';
  }

  const validationSchema = Yup.object().shape({});

  const minutes = [
    { key: '0', value: '0' },
    { key: '15', value: '15' },
    { key: '30', value: '30' },
    { key: '45', value: '45' },
  ];

  const hours = [];
  for (let i = 0; i < 13; i += 1) {
    hours[i] = { key: String(i), value: String(i) };
  }

  return (
    <>
      <ConnectedDeleteEvaluationReportConfirmationModal
        isOpen={isDeleteModelOpen}
        closeModal={toggleCancelModel}
        submitModal={cancelReport}
      />
      <Formik
        initialValues={initialValues}
        enableReinitialize
        onSubmit={handleSubmitSaveDraft}
        validationSchema={validationSchema}
        validateOnMount
      >
        {({ values, setFieldValue }) => {
          return (
            <Form className={classnames(formStyles.form, styles.form)}>
              <GridContainer className={styles.cardContainer}>
                <Grid row className={styles.evalInfoSection}>
                  <Grid col>
                    <h2>Evaluation form</h2>
                    <h3>Evaluation information</h3>
                    <DatePickerInput label="Date of inspection" name="inspectionDate" />
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
                    {values.evaluationType === 'physical' && (
                      <FormGroup>
                        <Fieldset>
                          <legend className="usa-label">Travel time to evaluation</legend>
                          <div className={styles.durationPickers}>
                            <div>
                              <DropdownInput
                                id="hour"
                                name="hour"
                                label="Hour"
                                className={styles.hourPicker}
                                onChange={(e) => {
                                  setFieldValue('hour', e.target.value);
                                }}
                                options={hours}
                              />
                            </div>
                            <div>
                              <DropdownInput
                                id="minute"
                                name="minute"
                                label="Minute"
                                className={styles.minutePicker}
                                onChange={(e) => {
                                  setFieldValue('minute', e.target.value);
                                }}
                                options={minutes}
                              />
                            </div>
                          </div>
                        </Fieldset>
                      </FormGroup>
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
                    {values.evaluationType === 'physical' && values.evaluationLocation === 'destination' && (
                      <DatePickerInput
                        label="Observed delivery date"
                        name="observedDate"
                        hint="Only enter a date here if the delivery you witnessed did not happen on the scheduled delivery date"
                      />
                    )}
                    {values.evaluationType === 'physical' && values.evaluationLocation === 'origin' && (
                      <DatePickerInput
                        label="Observed pickup date"
                        name="observedDate"
                        hint="Only enter a date here if the pickup you witnessed did not happen on the scheduled pickup date"
                      />
                    )}
                    <FormGroup>
                      <Fieldset>
                        <legend className="usa-label">Evaluation length</legend>
                        <div className={styles.durationPickers}>
                          <div>
                            <DropdownInput
                              id="hour"
                              name="evalLengthHour"
                              label="Hour"
                              className={styles.hourPicker}
                              onChange={(e) => {
                                setFieldValue('evalLengthHour', e.target.value);
                              }}
                              options={hours}
                            />
                          </div>
                          <div>
                            <DropdownInput
                              id="minute"
                              name="evalLengthMinute"
                              label="Minute"
                              className={styles.minutePicker}
                              onChange={(e) => {
                                setFieldValue('evalLengthMinute', e.target.value);
                              }}
                              options={minutes}
                            />
                          </div>
                        </div>
                      </Fieldset>
                    </FormGroup>
                  </Grid>
                </Grid>
                <Grid row className={styles.evalInfoSection}>
                  <Grid col>
                    <h3>Violations</h3>
                    <FormGroup className={styles.violationsGroup}>
                      <Fieldset className={styles.radioGroup}>
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
                        />
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
                      <Button className="usa-button--unstyled" onClick={toggleCancelModel} type="button">
                        Cancel
                      </Button>
                      <Button data-testid="saveDraft" type="submit" className="usa-button--secondary">
                        Save draft
                      </Button>
                      {values.violationsObserved === 'yes' ? (
                        <Button onClick={() => handleSelectViolations(values)} type="button">
                          Next: select violations
                        </Button>
                      ) : (
                        <Button disabled={!values.violationsObserved}>Review and submit</Button>
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

ShipmentEvaluationForm.propTypes = {
  evaluationReport: PropTypes.object,
};

ShipmentEvaluationForm.defaultProps = {
  evaluationReport: {},
};

export default ShipmentEvaluationForm;
