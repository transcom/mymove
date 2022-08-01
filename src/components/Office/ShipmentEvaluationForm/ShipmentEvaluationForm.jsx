import React, { useState } from 'react';
import 'styles/office.scss';
import {
  GridContainer,
  Grid,
  Button,
  Radio,
  FormGroup,
  Fieldset,
  Dropdown,
  Label,
  Textarea,
} from '@trussworks/react-uswds';
import { useParams, useHistory } from 'react-router';
import { useMutation } from 'react-query';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import classnames from 'classnames';

import styles from './ShipmentEvaluationForm.module.scss';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import { deleteEvaluationReport } from 'services/ghcApi';
import { DatePickerInput } from 'components/form/fields';
import Hint from 'components/Hint';

const ShipmentEvaluationForm = () => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

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

  const submitReport = async () => {};

  const initialValues = {};
  const validationSchema = Yup.object().shape({});

  const hours = [...Array(13).keys()];
  const minutes = [0, 15, 30, 45];

  return (
    <>
      <ConnectedDeleteEvaluationReportConfirmationModal
        isOpen={isDeleteModelOpen}
        closeModal={toggleCancelModel}
        submitModal={cancelReport}
      />
      <Formik initialValues={initialValues} onSubmit={submitReport} validationSchema={validationSchema} validateOnMount>
        {({ values }) => {
          return (
            <Form className={classnames(formStyles.form, styles.form)}>
              <GridContainer className={styles.cardContainer}>
                <Grid row className={styles.evalInfoSection}>
                  <Grid col>
                    <h2>Evaluation form</h2>
                    <h3>Evaluation information</h3>
                    <DatePickerInput label="Date of inspection" name="dateOfInspection" />
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
                              <Hint htmlFor="hour" className={styles.hourLabel}>
                                Hour
                              </Hint>
                              <Dropdown id="hour" name="hour" label="Hour" className={styles.hourPicker}>
                                {hours.map((option) => (
                                  <option key={option} value={option}>
                                    {option}
                                  </option>
                                ))}
                              </Dropdown>
                            </div>
                            <div>
                              <Hint htmlFor="minute" className={styles.minuteLabel}>
                                Minute
                              </Hint>
                              <Dropdown id="minute" name="minute" label="Minute" className={styles.minutePicker}>
                                {minutes.map((option) => (
                                  <option key={option} value={option}>
                                    {option}
                                  </option>
                                ))}
                              </Dropdown>
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
                        name="observedDeliveryDate"
                        hint="Only enter a date here if the delivery you witnessed did not happen on the scheduled delivery date"
                      />
                    )}
                    {values.evaluationType === 'physical' && values.evaluationLocation === 'origin' && (
                      <DatePickerInput
                        label="Observed pickup date"
                        name="observedPickupDate"
                        hint="Only enter a date here if the pickup you witnessed did not happen on the scheduled pickup date"
                      />
                    )}
                    <FormGroup>
                      <Fieldset>
                        <legend className="usa-label">Evaluation length</legend>
                        <div className={styles.durationPickers}>
                          <div>
                            <Hint htmlFor="hour" className={styles.hourLabel}>
                              Hour
                            </Hint>
                            <Dropdown id="hour" name="evalLengthHour" label="Hour" className={styles.hourPicker}>
                              {hours.map((option) => (
                                <option key={option} value={option}>
                                  {option}
                                </option>
                              ))}
                            </Dropdown>
                          </div>
                          <div>
                            <Hint htmlFor="minute" className={styles.minuteLabel}>
                              Minute
                            </Hint>
                            <Dropdown
                              id="minute"
                              name="evalLengthMinute"
                              label="Minute"
                              className={styles.minutePicker}
                            >
                              {minutes.map((option) => (
                                <option key={option} value={option}>
                                  {option}
                                </option>
                              ))}
                            </Dropdown>
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
                      name="evaluationRemarks"
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
                      <Button className="usa-button--secondary">Save draft</Button>
                      <Button type="submit" disabled>
                        Review and submit
                      </Button>
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

export default ShipmentEvaluationForm;
