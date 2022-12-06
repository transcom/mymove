import React, { useState } from 'react';
import { number } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Form, FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewWeightTicket.module.scss';

import { PPMShipmentShape, WeightTicketShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import { formatWeight } from 'utils/formatters';
import ppmDocumentStatus from 'constants/ppms';

const validationSchema = Yup.object().shape({
  emptyWeight: Yup.number().required('Enter the empty weight'),
  fullWeight: Yup.number()
    .required('Required')
    .when('emptyWeight', (emptyWeight, schema) => {
      return emptyWeight != null
        ? schema.min(emptyWeight + 1, 'The full weight must be greater than the empty weight')
        : schema;
    }),
  status: Yup.string().required('Reviewing this weight ticket is required'),
  rejectionReason: Yup.string().when('status', {
    is: ppmDocumentStatus.REJECTED,
    then: (schema) => schema.required('Add a reason why this weight ticket is rejected'),
  }),
});

export default function ReviewWeightTicket({ ppmShipment, weightTicket, tripNumber, ppmNumber }) {
  const [canEditRejection, setCanEditRejection] = useState(true);

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

  const initialValues = {
    weightType: missingEmptyWeightTicket || missingFullWeightTicket ? 'constructedWeight' : 'weightTicket',
    emptyWeight: emptyWeight ? `${emptyWeight}` : '',
    fullWeight: fullWeight ? `${fullWeight}` : '',
    ownsTrailer: ownsTrailer ? 'true' : 'false',
    trailerMeetsCriteria: trailerMeetsCriteria ? 'true' : 'false',
    status: status || '',
    rejectionReason: reason || '',
  };
  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik initialValues={initialValues} validationSchema={validationSchema}>
        {({ handleChange, values }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            setCanEditRejection(true);
          };

          return (
            <Form className={classnames(formStyles.form, styles.ReviewWeightTicket)}>
              <PPMHeaderSummary ppmShipment={ppmShipment} ppmNumber={ppmNumber} />
              <hr />
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
                  <legend className="usa-label">Did they use a trailer they owned?</legend>
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
              <Fieldset>
                <div
                  className={classnames(approveRejectStyles.statusOption, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.APPROVED,
                  })}
                >
                  <Radio
                    id={`approve-${weightTicket?.id}`}
                    checked={values.status === ppmDocumentStatus.APPROVED}
                    value={ppmDocumentStatus.APPROVED}
                    name="status"
                    label="Approve"
                    onChange={handleApprovalChange}
                    data-testid="approveRadio"
                  />
                </div>
                <div
                  className={classnames(approveRejectStyles.statusOption, styles.reject, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.REJECTED,
                  })}
                >
                  <Radio
                    id={`reject-${weightTicket?.id}`}
                    checked={values.status === ppmDocumentStatus.REJECTED}
                    value={ppmDocumentStatus.REJECTED}
                    name="status"
                    label="Reject"
                    onChange={handleChange}
                    data-testid="rejectRadio"
                  />

                  {values.status === ppmDocumentStatus.REJECTED && (
                    <FormGroup className={styles.reason}>
                      <Label htmlFor={`rejectReason-${weightTicket?.id}`}>Reason</Label>
                      {!canEditRejection && (
                        <p data-testid="rejectionReasonReadOnly">{weightTicket?.reason || values.rejectionReason}</p>
                      )}

                      {canEditRejection && (
                        <>
                          <Textarea
                            id={`rejectReason-${weightTicket?.id}`}
                            name="rejectionReason"
                            onChange={handleChange}
                            value={values.rejectionReason}
                            placeholder="Type something"
                          />
                          <div className={styles.hint}>500 characters</div>
                        </>
                      )}
                    </FormGroup>
                  )}
                </div>
              </Fieldset>
            </Form>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewWeightTicket.propTypes = {
  weightTicket: WeightTicketShape,
  ppmShipment: PPMShipmentShape,
  tripNumber: number.isRequired,
  ppmNumber: number.isRequired,
};

ReviewWeightTicket.defaultProps = {
  weightTicket: undefined,
  ppmShipment: undefined,
};
