import React, { useState } from 'react';
import { number } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Form, FormGroup, Label, Radio, Button, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewProGear.module.scss';

import { PPMShipmentShape, ProGearTicketShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import ppmDocumentStatus from 'constants/ppms';

const validationSchema = Yup.object().shape({
  selfProGear: Yup.bool().required('Required'),
  proGearWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .when('missingWeightTicket', {
      is: 'true',
      then: (schema) => schema.required('Enter the pro-gear weight'),
      other: (schema) => schema.required('Enter the constructed pro-gear weight'),
    }),
  description: Yup.string().required('Required'),
  missingWeightTicket: Yup.string(),
  status: Yup.string().required('Reviewing this pro-gear is required'),
  rejectionReason: Yup.string().when('status', {
    is: ppmDocumentStatus.REJECTED,
    then: (schema) => schema.required('Add a reason why this pro-gear is rejected').max(500),
  }),
});

export default function ReviewProGear({ ppmShipment, proGear, tripNumber, ppmNumber }) {
  const [canEditRejection, setCanEditRejection] = useState(true);

  const { description, selfProGear, proGearWeight, missingWeightTicket, status, reason } = proGear || {};

  let proGearValue;
  if (selfProGear === true) {
    proGearValue = 'true';
  }
  if (selfProGear === false) {
    proGearValue = 'false';
  }

  const initialValues = {
    belongsToSelf: proGearValue,
    status: status || '',
    rejectionReason: reason || '',
    missingWeightTicket: missingWeightTicket ? `${missingWeightTicket}` : '',
    description: description ? `${description}` : '',
    proGearWeight: proGearWeight ? `${proGearWeight}` : '',
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
            <Form className={classnames(formStyles.form, styles.reviewProGear)}>
              <PPMHeaderSummary ppmShipment={ppmShipment} ppmNumber={ppmNumber} />
              <hr />
              <h3 className={styles.tripNumber}>Pro-gear {tripNumber}</h3>
              <FormGroup>
                <Fieldset>
                  <legend className="usa-label">Belongs to</legend>
                  <Field
                    as={Radio}
                    id="customer"
                    label="Customer"
                    name="belongsToSelf"
                    value="true"
                    checked={values.belongsToSelf === 'true'}
                  />
                  <Field
                    as={Radio}
                    id="spouse"
                    label="Spouse"
                    name="belongsToSelf"
                    value="false"
                    checked={values.belongsToSelf === 'false'}
                  />
                </Fieldset>
              </FormGroup>
              <legend className={classnames('usa-label', styles.label)}>Description</legend>
              <div className={styles.displayValue}>{values.description}</div>
              <FormGroup>
                <Fieldset>
                  <legend className="usa-label">Pro-gear weight</legend>
                  <Field
                    as={Radio}
                    id="weight-tickets"
                    label="Weight tickets"
                    name="missingWeightTicket"
                    value="false"
                    checked={values.missingWeightTicket === 'false'}
                  />
                  <Field
                    as={Radio}
                    id="constructed-weight"
                    label="Constructed weight"
                    name="missingWeightTicket"
                    value="true"
                    checked={values.missingWeightTicket === 'true'}
                  />
                </Fieldset>
              </FormGroup>
              <MaskedTextField
                defaultValue="0"
                name="proGearWeight"
                label={
                  values.missingWeightTicket === 'true' ? 'Constructed pro-gear weight' : "Shipment's pro-gear weight"
                }
                id="proGearWeight"
                mask={Number}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
                suffix="lbs"
              />
              <h3 className={styles.reviewHeader}>Review pro-gear {tripNumber}</h3>
              <p>Add a review for this pro-gear</p>
              <Fieldset>
                <div
                  className={classnames(approveRejectStyles.statusOption, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.APPROVED,
                  })}
                >
                  <Radio
                    id={`approve-${proGear?.id}`}
                    checked={values.status === ppmDocumentStatus.APPROVED}
                    value={ppmDocumentStatus.APPROVED}
                    name="status"
                    label="Approve"
                    onChange={handleApprovalChange}
                    data-testid="approveRadio"
                  />
                </div>
                <div
                  className={classnames(approveRejectStyles.statusOption, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.REJECTED,
                  })}
                >
                  <Radio
                    id={`reject-${proGear?.id}`}
                    checked={values.status === ppmDocumentStatus.REJECTED}
                    value={ppmDocumentStatus.REJECTED}
                    name="status"
                    label="Reject"
                    onChange={handleChange}
                    data-testid="rejectRadio"
                  />

                  {values.status === ppmDocumentStatus.REJECTED && (
                    <FormGroup className={styles.rejectionReason}>
                      <Label htmlFor={`rejectReason-${proGear?.id}`}>Reason</Label>
                      {!canEditRejection && (
                        <>
                          <p data-testid="rejectionReasonReadOnly">{proGear?.reason || values.rejectionReason}</p>
                          <Button
                            type="button"
                            unstyled
                            data-testid="editReasonButton"
                            className={styles.clearStatus}
                            onClick={() => setCanEditRejection(true)}
                            aria-label="Edit reason button"
                          >
                            <span className="icon">
                              <FontAwesomeIcon icon="pen" title="Edit reason" alt="" />
                            </span>
                            <span aria-hidden="true">Edit reason</span>
                          </Button>
                        </>
                      )}

                      {canEditRejection && (
                        <>
                          <Textarea
                            id={`rejectReason-${proGear?.id}`}
                            name="rejectionReason"
                            onChange={handleChange}
                            value={values.rejectionReason}
                            placeholder="Type something"
                            maxLength={500}
                          />
                          <p className={styles.characters}>{500 - values.rejectionReason.length} characters</p>
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

ReviewProGear.propTypes = {
  proGear: ProGearTicketShape,
  ppmShipment: PPMShipmentShape,
  tripNumber: number.isRequired,
  ppmNumber: number.isRequired,
};

ReviewProGear.defaultProps = {
  proGear: undefined,
  ppmShipment: undefined,
};
