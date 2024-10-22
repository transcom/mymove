import React, { useState, useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { func, number, string, object } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Form, FormGroup, Label, Radio, Button, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewProGear.module.scss';

import { ErrorMessage } from 'components/form';
import { patchProGearWeightTicket } from 'services/ghcApi';
import { ProGearTicketShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import ppmDocumentStatus from 'constants/ppms';

const validationSchema = Yup.object().shape({
  belongsToSelf: Yup.bool().required('Required'),
  proGearWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .when('missingWeightTicket', {
      is: 'true',
      then: (schema) => schema.required('Enter the constructed pro-gear weight'),
      otherwise: (schema) => schema.required('Enter the weight with pro-gear'),
    }),
  description: Yup.string().required('Required'),
  missingWeightTicket: Yup.string(),
  status: Yup.string().required('Reviewing this pro-gear is required'),
  rejectionReason: Yup.string().when('status', {
    is: ppmDocumentStatus.REJECTED,
    then: (schema) => schema.required('Add a reason why this pro-gear is rejected'),
  }),
});

export default function ReviewProGear({
  ppmShipmentInfo,
  proGear,
  tripNumber,
  ppmNumber,
  onError,
  onSuccess,
  formRef,
  readOnly,
}) {
  const [canEditRejection, setCanEditRejection] = useState(true);

  const { mutate: patchProGearMutation } = useMutation(patchProGearWeightTicket, {
    onSuccess,
    onError,
  });

  const { belongsToSelf, description, hasWeightTickets, weight, status, reason } = proGear || {};

  const proGearValue = belongsToSelf ? 'true' : 'false';

  const missingWeightTicketValue = hasWeightTickets ? 'false' : 'true';

  const handleSubmit = (values) => {
    if (readOnly) {
      onSuccess();
      return;
    }

    let hasWeightTicketValue;
    if (values.missingWeightTicket === 'true') {
      hasWeightTicketValue = false;
    }
    if (values.missingWeightTicket === 'false') {
      hasWeightTicketValue = true;
    }
    const payload = {
      ppmShipmentId: proGear.ppmShipmentId,
      belongsToSelf: values.belongsToSelf === 'true',
      hasWeightTickets: hasWeightTicketValue,
      weight: parseInt(values.proGearWeight, 10),
      reason: values.status === ppmDocumentStatus.APPROVED ? null : values.rejectionReason,
      status: values.status,
    };
    patchProGearMutation({
      ppmShipmentId: proGear.ppmShipmentId,
      proGearWeightTicketId: proGear.id,
      payload,
      eTag: proGear.eTag,
    });
  };

  const initialValues = {
    belongsToSelf: proGearValue,
    status: status || '',
    rejectionReason: reason || '',
    missingWeightTicket: missingWeightTicketValue,
    description: description ? `${description}` : '',
    proGearWeight: weight ? `${weight}` : '',
  };

  useEffect(() => {
    if (formRef?.current) {
      formRef.current.resetForm();
      formRef.current.validateForm();
    }
  }, [formRef, proGear]);

  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        innerRef={formRef}
        onSubmit={handleSubmit}
        enableReinitialize
        validateOnMount
      >
        {({ handleChange, errors, touched, values }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            setCanEditRejection(true);
          };

          return (
            <>
              <div className={classnames(formStyles.form, styles.reviewProGear, styles.headerContainer)}>
                <PPMHeaderSummary
                  ppmShipmentInfo={ppmShipmentInfo}
                  ppmNumber={ppmNumber}
                  showAllFields={false}
                  readOnly={readOnly}
                />
              </div>
              <Form className={classnames(formStyles.form, styles.reviewProGear)}>
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
                      disabled={readOnly}
                    />
                    <Field
                      as={Radio}
                      id="spouse"
                      label="Spouse"
                      name="belongsToSelf"
                      value="false"
                      checked={values.belongsToSelf === 'false'}
                      disabled={readOnly}
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
                      disabled={readOnly}
                    />
                    <Field
                      as={Radio}
                      id="constructed-weight"
                      label="Constructed weight"
                      name="missingWeightTicket"
                      value="true"
                      checked={values.missingWeightTicket === 'true'}
                      disabled={readOnly}
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
                  disabled={readOnly}
                />
                <h3 className={styles.reviewHeader}>Review pro-gear {tripNumber}</h3>
                <p>Add a review for this pro-gear</p>
                <ErrorMessage display={!!errors?.status && !!touched?.status}>{errors.status}</ErrorMessage>
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
                      label="Accept"
                      onChange={handleApprovalChange}
                      data-testid="approveRadio"
                      className={styles.acceptRadio}
                      disabled={readOnly}
                    />
                  </div>
                  <div
                    className={classnames(approveRejectStyles.statusOption, styles.reject, {
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
                      className={styles.rejectRadio}
                      disabled={readOnly}
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
                              disabled={readOnly}
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
                            <ErrorMessage display={!!errors?.rejectionReason && !!touched?.rejectionReason}>
                              {errors.rejectionReason}
                            </ErrorMessage>
                            <Textarea
                              id={`rejectReason-${proGear?.id}`}
                              name="rejectionReason"
                              onChange={handleChange}
                              error={touched.rejectionReason ? errors.rejectionReason : null}
                              value={values.rejectionReason}
                              placeholder="Type something"
                              disabled={readOnly}
                            />
                            <div className={styles.hint}>{500 - values.rejectionReason.length} characters</div>
                          </>
                        )}
                      </FormGroup>
                    )}
                  </div>
                </Fieldset>
              </Form>
            </>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewProGear.propTypes = {
  proGear: ProGearTicketShape,
  tripNumber: number.isRequired,
  ppmNumber: string.isRequired,
  onSuccess: func,
  formRef: object,
};

ReviewProGear.defaultProps = {
  proGear: null,
  onSuccess: null,
  formRef: null,
};
