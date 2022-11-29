import React, { useState } from 'react';
import { string } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Form, FormGroup, Label, Radio } from '@trussworks/react-uswds';

import styles from './ReviewWeightTicket.module.scss';

import { PPMShipmentShape, WeightTicketShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import ApproveReject from 'components/form/ApproveReject/ApproveReject';
import formStyles from 'styles/form.module.scss';
import { formatWeight, formatDate, formatCentsTruncateWhole } from 'utils/formatters';

export default function ReviewWeightTicket({ ppmShipment, weightTicket, tripNumber, ppmNumber }) {
  const {
    vehicleDescription,
    missingEmptyWeightTicket,
    missingFullWeightTicket,
    emptyWeight,
    fullWeight,
    ownsTrailer,
    trailerMeetsCriteria,
    status,
    reason,
  } = weightTicket || {};
  const [canEditRejection, setCanEditRejection] = useState(true);
  const constructedOrWeightTicket =
    !missingEmptyWeightTicket && !missingFullWeightTicket ? 'weightTicket' : 'constructedWeight';
  const {
    actualPickupPostalCode,
    actualDestinationPostalCode,
    actualMoveDate,
    hasReceivedAdvance,
    advanceAmountReceived,
  } = ppmShipment || {};

  const initialValues = {
    weightType: weightTicket?.id ? constructedOrWeightTicket : '',
    emptyWeight: emptyWeight || '',
    fullWeight: fullWeight || '',
    ownsTrailer: ownsTrailer ? 'true' : 'false',
    trailerMeetsCriteria: trailerMeetsCriteria ? 'true' : 'false',
    status: status || '',
    rejectionReason: reason || '',
  };
  return (
    <div className={styles.container}>
      <Formik initialValues={initialValues}>
        {({ handleReset, handleChange, submitForm, setValues, values }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            submitForm().then(() => {
              setCanEditRejection(true);
            });
          };

          const handleRejectChange = (event) => {
            handleChange(event);
            submitForm().then(() => {
              setCanEditRejection(false);
            });
          };

          const handleRejectCancel = (event) => {
            if (initialValues.rejectionReason) {
              setCanEditRejection(false);
            }

            handleReset(event);
          };

          const handleFormReset = () => {
            setValues({
              status: 'REQUESTED',
              rejectionReason: '',
            });
            submitForm().then(() => {
              setCanEditRejection(true);
            });
          };

          return (
            <div className="container--accent--ppm">
              <Form className={classnames(formStyles.form, styles.ReviewWeightTicket)}>
                <header className={styles.header}>
                  <div>
                    <h2>PPM {ppmNumber}</h2>
                    <section>
                      <div>
                        <Label className={styles.headerLabel}>Departure date</Label>
                        <span className={styles.light}>{formatDate(actualMoveDate)}</span>
                      </div>
                      <div>
                        <Label className={styles.headerLabel}>Starting ZIP</Label>
                        <span className={styles.light}>{actualPickupPostalCode}</span>
                      </div>
                      <div>
                        <Label className={styles.headerLabel}>Ending ZIP</Label>
                        <span className={styles.light}>{actualDestinationPostalCode}</span>
                      </div>
                      <div>
                        <Label className={styles.headerLabel}>Advance recieved</Label>
                        <span className={styles.light}>
                          {hasReceivedAdvance ? `Yes, $${formatCentsTruncateWhole(advanceAmountReceived)}` : 'No'}
                        </span>
                      </div>
                    </section>
                  </div>
                  <hr />
                </header>
                <h3 className={styles.tripNumber}>Trip {tripNumber}</h3>
                <legend className={classnames('usa-label', styles.label)}>Vehicle description</legend>
                <div className={styles.displayValue}>{vehicleDescription}</div>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Weight type</legend>
                    <Field
                      as={Radio}
                      id="weight-tickets"
                      label="Weight tickets"
                      name="weightType"
                      value="weightTicket"
                      checked={values.weightType === 'weightTicket'}
                    />
                    <Field
                      as={Radio}
                      id="constructed-weight"
                      label="Constructed weight"
                      name="weightType"
                      value="constructedWeight"
                      checked={values.weightType === 'constructedWeight'}
                    />
                  </Fieldset>
                </FormGroup>
                <MaskedTextField
                  defaultValue="0"
                  name="emptyWeight"
                  label={values.weightType === 'weightTicket' ? 'Empty weight' : 'Empty constructed weight'}
                  id="emptyWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                />
                <MaskedTextField
                  defaultValue="0"
                  name="fullWeight"
                  label={values.weightType === 'weightTicket' ? 'Full weight' : 'Full constructed weight'}
                  id="fullWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                />
                <Label className={styles.label}>Net weight</Label>
                <div className={styles.displayValue}>{formatWeight(values.fullWeight - values.emptyWeight)}</div>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Did they use a trailer they owned</legend>
                    <Field
                      as={Radio}
                      id="ownsTrailerYes"
                      label="Yes"
                      name="ownsTrailer"
                      value="true"
                      checked={values.ownsTrailer === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="ownsTrailerNo"
                      label="No"
                      name="ownsTrailer"
                      value="false"
                      checked={values.ownsTrailer === 'false'}
                    />
                  </Fieldset>
                </FormGroup>
                {values.ownsTrailer === 'true' && (
                  <FormGroup>
                    <Fieldset>
                      <legend className="usa-label">{`Is the trailer's weight claimable?`}</legend>
                      <Field
                        as={Radio}
                        id="trailerCriteriaYes"
                        label="Yes"
                        name="trailerMeetsCriteria"
                        value="true"
                        checked={values.trailerMeetsCriteria === 'true'}
                      />
                      <Field
                        as={Radio}
                        id="trailerCriteriaNo"
                        label="No"
                        name="trailerMeetsCriteria"
                        value="false"
                        checked={values.trailerMeetsCriteria === 'false'}
                      />
                    </Fieldset>
                  </FormGroup>
                )}
                <h3 className={styles.reviewHeader}>Review trip {tripNumber}</h3>
                <p>Add a review for this weight ticket</p>
                <ApproveReject
                  id="ApproveReject"
                  currentStatus={values.status}
                  rejectionReason={values.rejectionReason}
                  requestComplete={false}
                  approvedStatus="APPROVED"
                  deniedStatus="DENIED"
                  canEditRejection={canEditRejection}
                  setCanEditRejection={setCanEditRejection}
                  handleApprovalChange={handleApprovalChange}
                  handleRejectChange={handleRejectChange}
                  handleRejectCancel={handleRejectCancel}
                  handleChange={handleChange}
                  handleFormReset={handleFormReset}
                />
              </Form>
            </div>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewWeightTicket.propTypes = {
  weightTicket: WeightTicketShape,
  ppmShipment: PPMShipmentShape,
  tripNumber: string.isRequired,
  ppmNumber: string.isRequired,
};

ReviewWeightTicket.defaultProps = {
  weightTicket: undefined,
  ppmShipment: undefined,
};
